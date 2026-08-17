package analytics

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

var overlapCategories = map[string]struct{}{
	"streaming": {},
	"music":     {},
	"cloud":     {},
}

var personalTariffRank = map[string]int{
	"mini":       1,
	"lite":       2,
	"student":    3,
	"basic":      4,
	"standard":   5,
	"plus":       6,
	"individual": 7,
	"premium":    8,
	"pro":        9,
	"max":        10,
}

var expensiveTariffs = map[string]struct{}{
	"premium": {},
	"pro":     {},
	"max":     {},
}

var familyTariffs = map[string]struct{}{
	"family": {},
	"duo":    {},
}

type crowdBucket struct {
	values []float64
}

type resolvedSub struct {
	Subscription
	service Service
	matched bool
}

func BuildRecommendations(input Input) []Recommendation {
	if input.Now.IsZero() {
		input.Now = time.Now()
	}

	if input.Convert == nil {
		input.Convert = func(amount float64, from string) float64 { return amount }
	}

	index := newServiceIndex(input.Services)
	country := countryCode(input.Country)
	recs := make([]Recommendation, 0, 16)

	resolved := make([]resolvedSub, 0, len(input.Subscriptions))
	var monthlyTotal float64

	for _, sub := range input.Subscriptions {
		item := resolvedSub{Subscription: sub}
		if service, ok := index.match(sub.Name); ok {
			item.service = service
			item.matched = true
		}

		resolved = append(resolved, item)
		if sub.IncludeInAnalytics {
			monthlyTotal += input.Convert(monthlyAmount(sub.SharePrice, sub.Period), sub.Currency)
		}
	}

	crowd := indexCrowd(index, input.Crowd)
	covered := make(map[uint64]struct{})

	recs = append(recs, overlapRecommendations(resolved, input.Categories, input.Convert, input.DisplayCurrency)...)
	recs = append(recs, familyRecommendations(resolved, input.Convert, input.DisplayCurrency)...)

	for _, rec := range crowdOverpayRecommendations(resolved, crowd, country, input.Convert, input.DisplayCurrency) {
		if rec.SubscriptionID != nil {
			covered[*rec.SubscriptionID] = struct{}{}
		}
		recs = append(recs, rec)
	}

	for _, rec := range downgradeRecommendations(resolved, crowd, country, input.Convert, input.DisplayCurrency) {
		if rec.SubscriptionID != nil {
			covered[*rec.SubscriptionID] = struct{}{}
		}
		recs = append(recs, rec)
	}

	recs = append(recs, yearlyCrowdRecommendations(resolved, crowd, country, input.Convert, input.DisplayCurrency)...)
	recs = append(recs, expensiveTariffRecommendations(resolved, monthlyTotal, covered, input.Convert, input.DisplayCurrency)...)
	recs = append(recs, hygieneRecommendations(input.Subscriptions, monthlyTotal, input.Convert, input.DisplayCurrency, input.Now)...)

	sort.SliceStable(recs, func(i, j int) bool {
		return recs[i].Priority > recs[j].Priority
	})

	if len(recs) > maxRecommendations {
		recs = recs[:maxRecommendations]
	}

	return recs
}

func overlapRecommendations(subs []resolvedSub, categories map[uint64][]string, convert func(float64, string) float64, displayCurrency string) []Recommendation {
	type grouped struct {
		sub     resolvedSub
		monthly float64
	}

	byCategory := map[string][]grouped{}

	for _, sub := range subs {
		if !sub.IncludeInAnalytics {
			continue
		}

		monthly := convert(monthlyAmount(sub.SharePrice, sub.Period), sub.Currency)
		slugs := categories[sub.ID]
		if len(slugs) == 0 {
			continue
		}

		for _, slug := range slugs {
			if _, ok := overlapCategories[slug]; !ok {
				continue
			}

			byCategory[slug] = append(byCategory[slug], grouped{sub: sub, monthly: monthly})
		}
	}

	recs := make([]Recommendation, 0, len(byCategory))

	for slug, items := range byCategory {
		if len(items) < 2 {
			continue
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].monthly < items[j].monthly
		})

		var savings float64
		names := make([]string, 0, len(items))
		for i, item := range items {
			names = append(names, displayName(item.sub.Subscription, item.sub.service))
			if i > 0 {
				savings += item.monthly
			}
		}

		if savings <= 0 {
			continue
		}

		recs = append(recs, Recommendation{
			ID:         "overlap-" + slug,
			Type:       "overlap",
			TitleKey:   "analytics.rec-overlap-title",
			DescKey:    "analytics.rec-overlap-desc",
			DescValues: map[string]any{"names": strings.Join(names, ", "), "category": slug, "amount": roundMoney(savings), "currency": displayCurrency},
			Priority:   100000 + savings,
		})
	}

	return recs
}

func familyRecommendations(subs []resolvedSub, convert func(float64, string) float64, displayCurrency string) []Recommendation {
	recs := make([]Recommendation, 0)

	for _, sub := range subs {
		if !sub.IncludeInAnalytics || !sub.IsOwner {
			continue
		}

		tariff := tariffOrNone(sub.Tariff)
		if _, ok := familyTariffs[tariff]; !ok {
			continue
		}

		if sub.SharePercent < 99.5 {
			continue
		}

		monthly := convert(monthlyAmount(sub.SharePrice, sub.Period), sub.Currency)
		name := displayName(sub.Subscription, sub.service)

		recs = append(recs, Recommendation{
			ID:             fmt.Sprintf("family-%d", sub.ID),
			Type:           "family-share",
			TitleKey:       "analytics.rec-family-title",
			DescKey:        "analytics.rec-family-desc",
			DescValues:     map[string]any{"name": name, "tariff": tariff, "amount": roundMoney(monthly), "currency": displayCurrency},
			SubscriptionID: ptrID(sub.ID),
			Priority:       70000 + monthly,
		})
	}

	return recs
}

func crowdOverpayRecommendations(subs []resolvedSub, crowd map[crowdKey]*crowdBucket, country string, convert func(float64, string) float64, displayCurrency string) []Recommendation {
	recs := make([]Recommendation, 0)

	for _, sub := range subs {
		if !sub.IncludeInAnalytics || !sub.matched {
			continue
		}

		tariff := tariffOrNone(sub.Tariff)
		if tariff == "none" {
			continue
		}

		period := "monthly"
		userPrice := monthlyAmount(sub.Price, sub.Period)
		if strings.EqualFold(sub.Period, "yearly") {
			period = "yearly"
			userPrice = sub.Price
		}

		median, sample, usedCountry := crowdMedian(crowd, sub.service.Slug, tariff, strings.ToUpper(sub.Currency), period, country)
		if sample < crowdMinSample && period == "yearly" {
			median, sample, usedCountry = crowdMedian(crowd, sub.service.Slug, tariff, strings.ToUpper(sub.Currency), "monthly", country)
			userPrice = monthlyAmount(sub.Price, sub.Period)
			period = "monthly"
		}

		if sample < crowdMinSample || median <= 0 || userPrice < median*crowdOverpayRatio {
			continue
		}

		percent := int(math.Round((userPrice/median - 1) * 100))
		if percent < 20 {
			continue
		}

		overpay := convert(userPrice-median, sub.Currency)
		if period == "yearly" {
			overpay = convert((userPrice-median)/12, sub.Currency)
		}
		values := map[string]any{
			"name":     sub.service.Name,
			"tariff":   tariff,
			"percent":  percent,
			"amount":   roundMoney(overpay),
			"currency": displayCurrency,
		}

		descKey := "analytics.rec-crowd-overpay-desc"
		if usedCountry != "" {
			descKey = "analytics.rec-crowd-overpay-country-desc"
			values["country"] = usedCountry
		}

		recs = append(recs, Recommendation{
			ID:             fmt.Sprintf("crowd-overpay-%d", sub.ID),
			Type:           "crowd-overpay",
			TitleKey:       "analytics.rec-crowd-overpay-title",
			DescKey:        descKey,
			DescValues:     values,
			SubscriptionID: ptrID(sub.ID),
			Priority:       80000 + overpay,
		})
	}

	return recs
}

func downgradeRecommendations(subs []resolvedSub, crowd map[crowdKey]*crowdBucket, country string, convert func(float64, string) float64, displayCurrency string) []Recommendation {
	recs := make([]Recommendation, 0)

	for _, sub := range subs {
		if !sub.IncludeInAnalytics || !sub.matched {
			continue
		}

		tariff := tariffOrNone(sub.Tariff)
		rank, ok := personalTariffRank[tariff]
		if !ok || rank <= personalTariffRank["basic"] {
			continue
		}

		userMonthly := monthlyAmount(sub.Price, sub.Period)
		bestTariff := ""
		bestMedian := 0.0

		for cheaper, cheaperRank := range personalTariffRank {
			if cheaperRank >= rank {
				continue
			}

			median, sample, _ := crowdMedian(crowd, sub.service.Slug, cheaper, strings.ToUpper(sub.Currency), "monthly", country)
			if sample < crowdMinSample || median <= 0 || median >= userMonthly*0.9 {
				continue
			}

			if bestTariff == "" || cheaperRank < personalTariffRank[bestTariff] {
				bestTariff = cheaper
				bestMedian = median
			}
		}

		if bestTariff == "" {
			continue
		}

		savings := convert(userMonthly-bestMedian, sub.Currency)
		typical := convert(bestMedian, sub.Currency)

		recs = append(recs, Recommendation{
			ID:       fmt.Sprintf("downgrade-%d", sub.ID),
			Type:     "downgrade",
			TitleKey: "analytics.rec-downgrade-title",
			DescKey:  "analytics.rec-downgrade-desc",
			DescValues: map[string]any{
				"name":           sub.service.Name,
				"tariff":         tariff,
				"cheaper_tariff": bestTariff,
				"typical":        roundMoney(typical),
				"amount":         roundMoney(savings),
				"currency":       displayCurrency,
			},
			SubscriptionID: ptrID(sub.ID),
			Priority:       90000 + savings,
		})
	}

	return recs
}

func yearlyCrowdRecommendations(subs []resolvedSub, crowd map[crowdKey]*crowdBucket, country string, convert func(float64, string) float64, displayCurrency string) []Recommendation {
	recs := make([]Recommendation, 0)

	for _, sub := range subs {
		if !sub.IncludeInAnalytics || !sub.matched || !strings.EqualFold(sub.Period, "monthly") {
			continue
		}

		tariff := tariffOrNone(sub.Tariff)
		yearlyMedian, sample, _ := crowdMedian(crowd, sub.service.Slug, tariff, strings.ToUpper(sub.Currency), "yearly", country)
		if sample < crowdMinSample || yearlyMedian <= 0 {
			continue
		}

		userYearly := monthlyAmount(sub.Price, sub.Period) * 12
		if yearlyMedian >= userYearly*0.85 {
			continue
		}

		savings := convert(userYearly-yearlyMedian, sub.Currency)
		percent := int(math.Round((1 - yearlyMedian/userYearly) * 100))
		name := displayName(sub.Subscription, sub.service)

		recs = append(recs, Recommendation{
			ID:       fmt.Sprintf("yearly-%d", sub.ID),
			Type:     "yearly-save",
			TitleKey: "analytics.rec-yearly-title",
			DescKey:  "analytics.rec-yearly-desc",
			DescValues: map[string]any{
				"name":     name,
				"percent":  percent,
				"amount":   roundMoney(savings),
				"currency": displayCurrency,
			},
			SubscriptionID: ptrID(sub.ID),
			Priority:       60000 + savings,
		})
	}

	return recs
}

func expensiveTariffRecommendations(subs []resolvedSub, monthlyTotal float64, covered map[uint64]struct{}, convert func(float64, string) float64, displayCurrency string) []Recommendation {
	if monthlyTotal <= 0 {
		return nil
	}

	recs := make([]Recommendation, 0)

	for _, sub := range subs {
		if !sub.IncludeInAnalytics {
			continue
		}

		if _, seen := covered[sub.ID]; seen {
			continue
		}

		tariff := tariffOrNone(sub.Tariff)
		if _, ok := expensiveTariffs[tariff]; !ok {
			continue
		}

		monthly := convert(monthlyAmount(sub.SharePrice, sub.Period), sub.Currency)
		if monthly/monthlyTotal < expensiveShareMin {
			continue
		}

		percent := int(math.Round((monthly / monthlyTotal) * 100))
		name := displayName(sub.Subscription, sub.service)

		recs = append(recs, Recommendation{
			ID:             fmt.Sprintf("expensive-%d", sub.ID),
			Type:           "expensive-tariff",
			TitleKey:       "analytics.rec-expensive-tariff-title",
			DescKey:        "analytics.rec-expensive-tariff-desc",
			DescValues:     map[string]any{"name": name, "tariff": tariff, "percent": percent, "currency": displayCurrency},
			SubscriptionID: ptrID(sub.ID),
			Priority:       50000 + monthly,
		})
	}

	return recs
}

func hygieneRecommendations(subs []Subscription, monthlyTotal float64, convert func(float64, string) float64, displayCurrency string, now time.Time) []Recommendation {
	recs := make([]Recommendation, 0)
	analyticsSubs := make([]Subscription, 0, len(subs))
	excluded := 0

	for _, sub := range subs {
		if !sub.IncludeInAnalytics {
			excluded++
			continue
		}

		analyticsSubs = append(analyticsSubs, sub)
	}

	if excluded > 0 {
		recs = append(recs, Recommendation{
			ID:         "excluded",
			Type:       "excluded",
			TitleKey:   "analytics.rec-excluded-title",
			DescKey:    "analytics.rec-excluded-desc",
			DescValues: map[string]any{"count": excluded},
			Priority:   5000,
		})
	}

	if monthlyTotal > 0 {
		sorted := append([]Subscription(nil), analyticsSubs...)
		sort.Slice(sorted, func(i, j int) bool {
			left := convert(monthlyAmount(sorted[i].SharePrice, sorted[i].Period), sorted[i].Currency)
			right := convert(monthlyAmount(sorted[j].SharePrice, sorted[j].Period), sorted[j].Currency)
			return left > right
		})

		if len(sorted) > 0 {
			top := convert(monthlyAmount(sorted[0].SharePrice, sorted[0].Period), sorted[0].Currency)
			if top/monthlyTotal >= concentrationShare {
				percent := int(math.Round((top / monthlyTotal) * 100))
				recs = append(recs, Recommendation{
					ID:             fmt.Sprintf("concentration-%d", sorted[0].ID),
					Type:           "concentration",
					TitleKey:       "analytics.rec-concentration-title",
					DescKey:        "analytics.rec-concentration-desc",
					DescValues:     map[string]any{"name": sorted[0].Name, "percent": percent},
					SubscriptionID: ptrID(sorted[0].ID),
					Priority:       40000 + top,
				})
			}
		}
	}

	dueSoon := make([]Subscription, 0)
	var dueTotal float64
	for _, sub := range analyticsSubs {
		days := daysUntil(sub.DatePay, sub.Period, now)
		if days >= 0 && days <= 7 {
			dueSoon = append(dueSoon, sub)
			dueTotal += convert(sub.SharePrice, sub.Currency)
		}
	}

	if len(dueSoon) >= 3 {
		recs = append(recs, Recommendation{
			ID:         "cluster",
			Type:       "cluster",
			TitleKey:   "analytics.rec-cluster-title",
			DescKey:    "analytics.rec-cluster-desc",
			DescValues: map[string]any{"count": len(dueSoon), "amount": roundMoney(dueTotal), "currency": displayCurrency},
			Priority:   20000 + dueTotal,
		})
	}

	if monthlyTotal > 0 {
		small := make([]Subscription, 0)
		var smallTotal float64
		for _, sub := range analyticsSubs {
			monthly := convert(monthlyAmount(sub.SharePrice, sub.Period), sub.Currency)
			if monthly > 0 && monthly/monthlyTotal < smallShareMax {
				small = append(small, sub)
				smallTotal += monthly
			}
		}

		if len(small) >= 4 {
			recs = append(recs, Recommendation{
				ID:         "small-subs",
				Type:       "small-subs",
				TitleKey:   "analytics.rec-small-title",
				DescKey:    "analytics.rec-small-desc",
				DescValues: map[string]any{"count": len(small), "amount": roundMoney(smallTotal), "currency": displayCurrency},
				Priority:   10000 + smallTotal,
			})
		}
	}

	next30 := outflowWithinDays(analyticsSubs, 30, convert, now)
	if monthlyTotal > 0 && next30 > monthlyTotal*1.35 {
		recs = append(recs, Recommendation{
			ID:         "upcoming-heavy",
			Type:       "upcoming-heavy",
			TitleKey:   "analytics.rec-upcoming-title",
			DescKey:    "analytics.rec-upcoming-desc",
			DescValues: map[string]any{"amount": roundMoney(next30), "currency": displayCurrency},
			Priority:   30000 + next30,
		})
	}

	return recs
}

type crowdKey struct {
	slug     string
	tariff   string
	currency string
	period   string
	country  string
}

func indexCrowd(index *serviceIndex, rows []CrowdPrice) map[crowdKey]*crowdBucket {
	buckets := map[crowdKey]*crowdBucket{}

	for _, row := range rows {
		service, ok := index.match(row.Name)
		if !ok {
			continue
		}

		period := strings.ToLower(strings.TrimSpace(row.Period))
		if period == "" {
			period = "monthly"
		}

		monthly := monthlyAmount(row.Price, period)
		if monthly <= 0 {
			continue
		}

		storePeriod := period
		value := monthly
		if period == "yearly" {
			storePeriod = "yearly"
			value = row.Price
		} else {
			storePeriod = "monthly"
		}

		put := func(country string) {
			key := crowdKey{
				slug:     service.Slug,
				tariff:   tariffOrNone(row.Tariff),
				currency: strings.ToUpper(row.Currency),
				period:   storePeriod,
				country:  country,
			}

			bucket := buckets[key]
			if bucket == nil {
				bucket = &crowdBucket{}
				buckets[key] = bucket
			}

			bucket.values = append(bucket.values, value)
		}

		put("")
		if code := countryCode(row.Country); code != "" {
			put(code)
		}
	}

	return buckets
}

func crowdMedian(crowd map[crowdKey]*crowdBucket, slug, tariff, currency, period, country string) (median float64, sample int, usedCountry string) {
	if country != "" {
		if bucket := crowd[crowdKey{slug: slug, tariff: tariff, currency: currency, period: period, country: country}]; bucket != nil && len(bucket.values) >= crowdMinSample {
			return medianOf(bucket.values), len(bucket.values), country
		}
	}

	if bucket := crowd[crowdKey{slug: slug, tariff: tariff, currency: currency, period: period, country: ""}]; bucket != nil && len(bucket.values) >= crowdMinSample {
		return medianOf(bucket.values), len(bucket.values), ""
	}

	return 0, 0, ""
}

func medianOf(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	cloned := append([]float64(nil), values...)
	sort.Float64s(cloned)
	mid := len(cloned) / 2
	if len(cloned)%2 == 1 {
		return cloned[mid]
	}

	return (cloned[mid-1] + cloned[mid]) / 2
}

func daysUntil(datePay time.Time, period string, now time.Time) int {
	next := nextPaymentDate(datePay, period, now)
	return int(next.Sub(now.Truncate(24*time.Hour)).Hours() / 24)
}

func nextPaymentDate(datePay time.Time, period string, now time.Time) time.Time {
	next := time.Date(datePay.Year(), datePay.Month(), datePay.Day(), 0, 0, 0, 0, time.UTC)
	today := now.Truncate(24 * time.Hour)

	for next.Before(today) {
		if strings.EqualFold(period, "yearly") {
			next = next.AddDate(1, 0, 0)
		} else {
			next = next.AddDate(0, 1, 0)
		}
	}

	return next
}

func outflowWithinDays(subs []Subscription, days int, convert func(float64, string) float64, now time.Time) float64 {
	today := now.Truncate(24 * time.Hour)
	horizon := today.AddDate(0, 0, days)
	var total float64

	for _, sub := range subs {
		next := nextPaymentDate(sub.DatePay, sub.Period, now)
		for !next.After(horizon) {
			if !next.Before(today) {
				total += convert(sub.SharePrice, sub.Currency)
			}

			if strings.EqualFold(sub.Period, "yearly") {
				next = next.AddDate(1, 0, 0)
			} else {
				next = next.AddDate(0, 1, 0)
			}
		}
	}

	return total
}

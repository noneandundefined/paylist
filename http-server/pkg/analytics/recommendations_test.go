package analytics

import (
	"testing"
	"time"
)

func identityConvert(amount float64, _ string) float64 {
	return amount
}

func TestFamilyShareRecommendation(t *testing.T) {
	recs := BuildRecommendations(Input{
		Subscriptions: []Subscription{{
			ID: 1, Name: "Netflix", Tariff: "family", Price: 1000, SharePrice: 1000, SharePercent: 100,
			Currency: "RUB", Period: "monthly", IncludeInAnalytics: true, IsOwner: true,
		}},
		DisplayCurrency: "RUB",
		Convert:         identityConvert,
		Now:             time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC),
	})

	if len(recs) == 0 || recs[0].Type != "family-share" {
		t.Fatalf("expected family-share, got %#v", recs)
	}
}

func TestCrowdOverpayAndDowngrade(t *testing.T) {
	services := []Service{{Slug: "netflix", Name: "Netflix", Category: "streaming", Aliases: []string{"нетфликс"}}}
	crowd := make([]CrowdPrice, 0, 12)
	for i := 0; i < 5; i++ {
		crowd = append(crowd, CrowdPrice{Name: "Netflix", Tariff: "premium", Price: 1000, Currency: "RUB", Period: "monthly", Country: "RU"})
		crowd = append(crowd, CrowdPrice{Name: "Netflix", Tariff: "basic", Price: 700, Currency: "RUB", Period: "monthly", Country: "RU"})
	}

	recs := BuildRecommendations(Input{
		Subscriptions: []Subscription{{
			ID: 7, Name: "Netflix Premium", Tariff: "premium", Price: 1300, SharePrice: 1300, SharePercent: 100,
			Currency: "RUB", Period: "monthly", IncludeInAnalytics: true, IsOwner: true,
		}},
		Services:        services,
		Crowd:           crowd,
		DisplayCurrency: "RUB",
		Country:         "RU",
		Convert:         identityConvert,
		Now:             time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC),
	})

	var overpay, downgrade bool
	for _, rec := range recs {
		if rec.Type == "crowd-overpay" {
			overpay = true
			if rec.DescValues["percent"] != 30 {
				t.Fatalf("overpay percent: %#v", rec.DescValues["percent"])
			}
			if rec.DescValues["country"] != "RU" {
				t.Fatalf("country: %#v", rec.DescValues["country"])
			}
		}
		if rec.Type == "downgrade" {
			downgrade = true
			if rec.DescValues["cheaper_tariff"] != "basic" {
				t.Fatalf("cheaper tariff: %#v", rec.DescValues["cheaper_tariff"])
			}
			if rec.DescValues["typical"] != 700.0 {
				t.Fatalf("typical: %#v", rec.DescValues["typical"])
			}
			if rec.DescValues["amount"] != 600.0 {
				t.Fatalf("savings: %#v", rec.DescValues["amount"])
			}
		}
	}

	if !overpay || !downgrade {
		t.Fatalf("expected crowd-overpay and downgrade, got %#v", recs)
	}
}

func TestOverlapCategory(t *testing.T) {
	recs := BuildRecommendations(Input{
		Subscriptions: []Subscription{
			{ID: 1, Name: "Netflix", Price: 400, SharePrice: 400, SharePercent: 100, Currency: "RUB", Period: "monthly", IncludeInAnalytics: true},
			{ID: 2, Name: "IVI", Price: 300, SharePrice: 300, SharePercent: 100, Currency: "RUB", Period: "monthly", IncludeInAnalytics: true},
		},
		Categories:      map[uint64][]string{1: {"streaming"}, 2: {"streaming"}},
		DisplayCurrency: "RUB",
		Convert:         identityConvert,
		Now:             time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC),
	})

	found := false
	for _, rec := range recs {
		if rec.Type == "overlap" && rec.DescValues["amount"] == 400.0 {
			found = true
		}
	}

	if !found {
		t.Fatalf("expected overlap savings 400, got %#v", recs)
	}
}

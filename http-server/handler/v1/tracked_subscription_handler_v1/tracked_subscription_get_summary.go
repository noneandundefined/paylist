package tracked_subscription_handler_v1

import (
	"context"
	"math"
	"net/http"
	"strings"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/currency"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/premium"
	"paylist.server/types"
)

func resolveDominantCurrency(ctx context.Context, subscriptions []models.TrackedSubscription) string {
	if len(subscriptions) == 0 {
		return "USD"
	}

	type currencyStats struct {
		count       int
		totalNative float64
		totalUSD    float64
		hasUSDTotal bool
	}

	stats := map[string]currencyStats{}

	for _, sub := range subscriptions {
		if !sub.IncludeInAnalytics {
			continue
		}

		cur := strings.ToUpper(sub.Currency)
		if cur == "" {
			cur = "USD"
		}

		monthlyAmount := currency.GetMonthlyAmount(sub.SharePrice, sub.Period)

		entry := stats[cur]
		entry.count++
		entry.totalNative += monthlyAmount

		converted, err := currency.Convert(ctx, cur, "USD", monthlyAmount)
		if err == nil {
			entry.totalUSD += converted
			entry.hasUSDTotal = true
		}

		stats[cur] = entry
	}

	if len(stats) == 0 {
		return "USD"
	}

	displayCurrency := "USD"
	bestUSD := -1.0
	bestCount := -1
	bestNative := -1.0
	useUSDTotals := false

	for _, entry := range stats {
		if entry.hasUSDTotal {
			useUSDTotals = true
			break
		}
	}

	for cur, entry := range stats {
		if useUSDTotals {
			if entry.totalUSD > bestUSD ||
				(entry.totalUSD == bestUSD && entry.count > bestCount) ||
				(entry.totalUSD == bestUSD && entry.count == bestCount && entry.totalNative > bestNative) {
				displayCurrency = cur
				bestUSD = entry.totalUSD
				bestCount = entry.count
				bestNative = entry.totalNative
			}

			continue
		}

		if entry.count > bestCount || (entry.count == bestCount && entry.totalNative > bestNative) {
			displayCurrency = cur
			bestCount = entry.count
			bestNative = entry.totalNative
		}
	}

	return displayCurrency
}

func (h *Handler) GetSubscriptionSummaryHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	subscriptions, err := h.Store.TrackedSubscriptions.Get_AllSubscriptionsByUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	displayCurrency := "USD"
	if premium.IsPremiumPlan(authToken) {
		displayCurrency = resolveDominantCurrency(ctx, subscriptions)

		settings, err := h.Store.Users.Get_UserSettingsByUserUuid(ctx, authToken.User.UserUUID)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if settings.DisplayCurrency != nil {
			currencyCode := strings.ToUpper(strings.TrimSpace(*settings.DisplayCurrency))
			if currencyCode != "" {
				displayCurrency = currencyCode
			}
		}
	}

	previewIDs := make([]uint64, 0, 3)
	for i, sub := range subscriptions {
		if i >= 3 {
			break
		}

		previewIDs = append(previewIDs, sub.ID)
	}

	var totalAmount float64

	for _, sub := range subscriptions {
		if !sub.IncludeInAnalytics {
			continue
		}

		monthlyAmount := currency.GetMonthlyAmount(sub.SharePrice, sub.Period)
		subCurrency := strings.ToUpper(sub.Currency)
		if subCurrency == "" {
			subCurrency = "USD"
		}

		converted, err := currency.Convert(ctx, subCurrency, displayCurrency, monthlyAmount)
		if err != nil {
			return httperr.BadRequest(err.Error())
		}

		totalAmount += converted
	}

	totalAmount = math.Round(totalAmount*100) / 100

	response := SubscriptionSummary{
		DisplayCurrency:        displayCurrency,
		TotalAmount:            totalAmount,
		ActiveCount:            len(subscriptions),
		PreviewSubscriptionIDs: previewIDs,
		ComparisonPercent:      0,
		ComparisonDirection:    "less",
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, response)
	return nil
}

package tracked_subscription_handler_v1

import (
	"net/http"
	"strings"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/analytics"
	"paylist.server/pkg/currency"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/premium"
	"paylist.server/types"
)

type analyticsRecommendationsResponse struct {
	Recommendations []analytics.Recommendation `json:"recommendations"`
}

func (h *Handler) GetSubscriptionAnalyticsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	subscriptions, err := h.Store.TrackedSubscriptions.Get_AllSubscriptionsByUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	categoryMap, err := h.Store.TrackedSubscriptions.Get_CategorySlugsMapByUserUUID(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	services, err := h.Store.TrackedSubscriptions.Get_Services(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	crowdRows, err := h.Store.TrackedSubscriptions.Get_CrowdSubscriptionPrices(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	settings, err := h.Store.Users.Get_UserSettingsByUserUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	displayCurrency := "USD"
	if premium.IsPremiumPlan(authToken) {
		displayCurrency = resolveDominantCurrency(ctx, subscriptions)

		if settings.DisplayCurrency != nil {
			currencyCode := strings.ToUpper(strings.TrimSpace(*settings.DisplayCurrency))
			if currencyCode != "" {
				displayCurrency = currencyCode
			}
		}
	}

	country := ""
	if settings.Country != nil {
		country = strings.ToUpper(strings.TrimSpace(*settings.Country))
	}

	input := analytics.Input{
		Subscriptions:   mapAnalyticsSubscriptions(subscriptions),
		Categories:      categoryMap,
		Services:        mapAnalyticsServices(services),
		Crowd:           mapCrowdPrices(crowdRows),
		DisplayCurrency: displayCurrency,
		Country:         country,
		Convert: func(amount float64, from string) float64 {
			from = strings.ToUpper(strings.TrimSpace(from))
			if from == "" {
				from = "USD"
			}

			if from == displayCurrency {
				return amount
			}

			converted, convErr := currency.Convert(ctx, from, displayCurrency, amount)
			if convErr != nil {
				return amount
			}

			return converted
		},
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, analyticsRecommendationsResponse{
		Recommendations: analytics.BuildRecommendations(input),
	})
	return nil
}

func mapAnalyticsSubscriptions(items []models.TrackedSubscription) []analytics.Subscription {
	out := make([]analytics.Subscription, 0, len(items))
	for _, item := range items {
		out = append(out, analytics.Subscription{
			ID:                 item.ID,
			Name:               item.Name,
			Tariff:             item.Tariff,
			Price:              item.Price,
			SharePrice:         item.SharePrice,
			SharePercent:       item.SharePercent,
			Currency:           item.Currency,
			Period:             item.Period,
			DatePay:            item.DatePay,
			IncludeInAnalytics: item.IncludeInAnalytics,
			IsOwner:            item.IsOwner,
		})
	}

	return out
}

func mapAnalyticsServices(items []models.Service) []analytics.Service {
	out := make([]analytics.Service, 0, len(items))
	for _, item := range items {
		out = append(out, analytics.Service{
			Slug:     item.Slug,
			Name:     item.Name,
			Category: item.Category,
			Aliases:  item.Aliases,
		})
	}

	return out
}

func mapCrowdPrices(items []models.CrowdSubscriptionPrice) []analytics.CrowdPrice {
	out := make([]analytics.CrowdPrice, 0, len(items))
	for _, item := range items {
		country := ""
		if item.Country != nil {
			country = *item.Country
		}

		out = append(out, analytics.CrowdPrice{
			Name:     item.Name,
			Tariff:   item.Tariff,
			Price:    item.Price,
			Currency: item.Currency,
			Period:   item.Period,
			Country:  country,
		})
	}

	return out
}

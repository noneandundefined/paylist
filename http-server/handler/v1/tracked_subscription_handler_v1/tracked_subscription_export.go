package tracked_subscription_handler_v1

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"paylist.server/infra/constants"
	"paylist.server/infra/locale"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
)

func (h *Handler) ExportSubscriptionsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	subs, err := h.Store.TrackedSubscriptions.Get_SubscriptionsByUuid(ctx, authToken.User.UserUUID, "")
	if err != nil {
		return httperr.Db(ctx, err)
	}

	categoryMap, err := h.Store.TrackedSubscriptions.Get_CategorySlugsMapByUserUUID(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	list := []models.TrackedSubscription{}
	if subs != nil {
		list = *subs
	}

	payload, err := buildSubscriptionsCSV(tr, list, categoryMap)
	if err != nil {
		return httperr.InternalServerError(tr.TErr("error.server-error"))
	}

	filename := fmt.Sprintf("paylist-subscriptions-%s.csv", time.Now().UTC().Format(time.DateOnly))
	httpx.HttpFileResponse(w, r, filename, payload, "text/csv; charset=utf-8")
	return nil
}

func buildSubscriptionsCSV(tr locale.Translator, subs []models.TrackedSubscription, categoryMap map[uint64][]string) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{
		tr.T("csv.name"),
		tr.T("csv.tariff"),
		tr.T("csv.price"),
		tr.T("csv.currency"),
		tr.T("csv.period"),
		tr.T("csv.date-pay"),
		tr.T("csv.auto-renewal"),
		tr.T("csv.notification"),
		tr.T("csv.analytics"),
		tr.T("csv.share-percent"),
		tr.T("csv.share-price"),
		tr.T("csv.is-owner"),
		tr.T("csv.categories"),
		tr.T("csv.note"),
	}); err != nil {
		return nil, err
	}

	if categoryMap == nil {
		categoryMap = map[uint64][]string{}
	}

	for _, sub := range subs {
		categories := categoryMap[sub.ID]
		note := ""
		if sub.Note != nil {
			note = *sub.Note
		}

		if err := writer.Write([]string{
			sub.Name,
			localizeCSVTariff(tr, sub.Tariff),
			formatCSVNumber(sub.Price),
			sub.Currency,
			localizeCSVPeriod(tr, sub.Period),
			sub.DatePay.UTC().Format(time.DateOnly),
			formatCSVBool(tr, sub.AutoRenewal),
			formatCSVBool(tr, sub.Notification),
			formatCSVBool(tr, sub.IncludeInAnalytics),
			formatCSVNumber(sub.SharePercent),
			formatCSVNumber(sub.SharePrice),
			formatCSVBool(tr, sub.IsOwner),
			strings.Join(categories, "; "),
			note,
		}); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func formatCSVNumber(value float64) string {
	formatted := strconv.FormatFloat(value, 'f', 3, 64)
	formatted = strings.TrimRight(formatted, "0")
	return strings.TrimRight(formatted, ".")
}

func formatCSVBool(tr locale.Translator, value bool) string {
	if value {
		return tr.T("csv.yes")
	}

	return tr.T("csv.no")
}

func localizeCSVPeriod(tr locale.Translator, period string) string {
	if period == "yearly" {
		return tr.T("csv.period-yearly")
	}

	return tr.T("csv.period-monthly")
}

func localizeCSVTariff(tr locale.Translator, tariff string) string {
	normalized := constants.NormalizeTariff(tariff)

	return tr.T("csv.tariff-" + normalized)
}

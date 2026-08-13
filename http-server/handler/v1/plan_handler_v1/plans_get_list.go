package plan_handler_v1

import (
	"net/http"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
)

const monthlyPlanMaxDurationDays = 31

func (h *Handler) GetPlansHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	subs, err := h.Store.Subscriptions.Get_Subscriptions(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	plans := make([]models.Subscription, 0, len(subs))

	for _, sub := range subs {
		if sub.Amount <= 0 {
			continue
		}

		if sub.DurationDays > monthlyPlanMaxDurationDays {
			continue
		}

		plans = append(plans, sub)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, plans)
	return nil
}

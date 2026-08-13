package tracked_subscription_handler_v1

import (
	"net/http"
	"strconv"

	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"

	"github.com/gorilla/mux"
)

func (h *Handler) GetSubscriptionByIdHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-not-found"))
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return httperr.BadRequest(tr.TErr("error.tracked-subscription-invalid-id"))
	}

	sub, err := h.Store.TrackedSubscriptions.Get_SubscriptionById(ctx, uint64(id), authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if sub == nil {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-not-found"))
	}

	categories, err := h.Store.TrackedSubscriptions.Get_CategorySlugsBySubscriptionID(ctx, uint64(id))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, TrackedSubscriptionDetailResponse{
		TrackedSubscription: *sub,
		Categories:          categories,
	})
	return nil
}

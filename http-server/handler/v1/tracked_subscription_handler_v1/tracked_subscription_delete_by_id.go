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

func (h *Handler) DeleteSubscriptionHandler_V1(w http.ResponseWriter, r *http.Request) error {
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

	if err := h.Store.TrackedSubscriptions.Delete_SubscriptionById(ctx, id, authToken.User.UserUUID); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.tracked-subscription-deleted"))
	return nil
}

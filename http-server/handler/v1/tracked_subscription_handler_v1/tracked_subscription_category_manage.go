package tracked_subscription_handler_v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"paylist.server/middleware"
	"paylist.server/pkg/categoryslug"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/premium"
	"paylist.server/pkg/profanity"
	"paylist.server/types"
)

func (h *Handler) CreateSubscriptionCategoryHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !premium.IsPremiumPlan(authToken) {
		return httperr.Forbidden(tr.TErr("error.premium-required"))
	}

	var payload *SubscriptionCategoryCreatePayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	label := strings.TrimSpace(payload.Label)
	slug := categoryslug.FromLabel(label)

	if slug == "" {
		return httperr.BadRequest(tr.TErr("error.category-label-invalid"))
	}

	if err := profanity.Reject(ctx, tr, "category-label", label); err != nil {
		return err
	}

	category, err := h.Store.TrackedSubscriptions.Create_UserSubscriptionCategory(ctx, authToken.User.UserUUID, slug, label)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return httperr.BadRequest(tr.TErr("error.category-slug-exists"))
		}

		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusCreated, mapSubscriptionCategoryResponse(*category))
	return nil
}

func (h *Handler) DeleteSubscriptionCategoryHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !premium.IsPremiumPlan(authToken) {
		return httperr.Forbidden(tr.TErr("error.premium-required"))
	}

	idParam := mux.Vars(r)["categoryId"]
	if idParam == "" {
		return httperr.NotFound(tr.TErr("error.category-not-found"))
	}

	categoryID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || categoryID == 0 {
		return httperr.BadRequest(tr.TErr("error.category-not-found"))
	}

	if err := h.Store.TrackedSubscriptions.Delete_UserSubscriptionCategory(ctx, authToken.User.UserUUID, categoryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return httperr.NotFound(tr.TErr("error.category-not-found"))
		}

		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("success.category-deleted"))
	return nil
}

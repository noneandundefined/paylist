package payment_handler_v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"paylist.server/infra/constants"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/yookassa"
	"paylist.server/types"
)

func (h *Handler) PostCheckoutHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if h.Yookassa == nil || !h.Yookassa.Configured() {
		return httperr.ServiceUnavailable(tr.TErr("error.payment-not-configured"))
	}

	if h.Yookassa.ReturnURL() == "" {
		return httperr.ServiceUnavailable(tr.TErr("error.payment-return-url-missing"))
	}

	var payload CheckoutPayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	count, err := h.Store.Payments.Get_PaymentActiveCount(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if count >= constants.MAX_PAYMENTS_COUNT {
		return httperr.Conflict(tr.TErr("error.too-many-active-payments"))
	}

	planName := strings.TrimSpace(payload.PlanName)

	plan, err := h.Store.Subscriptions.Get_SubscriptionByPlanName(ctx, planName)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if plan == nil || plan.Amount <= 0 {
		return httperr.BadRequest(tr.TErr("error.plan-not-available"))
	}

	loginState, err := h.Store.Users.Get_UserLoginStateByUserUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if loginState != nil && strings.EqualFold(loginState.PlanName, plan.PlanName) {
		return httperr.Conflict(tr.TErr("error.already-on-plan"))
	}

	amountValue := fmt.Sprintf("%.2f", plan.Amount)
	idempotenceKey := uuid.NewString()

	returnURL := fmt.Sprintf("%s?payment_id=%s", h.Yookassa.ReturnURL(), idempotenceKey)

	ykPayment, err := h.Yookassa.CreatePayment(ctx, idempotenceKey, yookassa.CreatePaymentRequest{
		Amount: yookassa.Amount{
			Value:    amountValue,
			Currency: plan.Currency,
		},
		Capture: true,
		Confirmation: yookassa.Confirmation{
			Type:      "redirect",
			ReturnURL: returnURL,
		},
		Description: fmt.Sprintf("Subs paylist.app - %s", plan.PlanName),
		Metadata: map[string]string{
			"user_uuid": authToken.User.UserUUID,
			"plan_name": plan.PlanName,
		},
		SavePaymentMethod: false,
	})
	if err != nil {
		logger.Error("PostCheckoutHandler_V1: %s", err.Error())
		return httperr.InternalServerError(tr.TErr("error.payment-create-failed"))
	}

	description := ykPayment.Description
	_, err = h.Store.Payments.Create_PaymentHistory(ctx, &models.PaymentHistory{
		UserUUID:          authToken.User.UserUUID,
		PlanName:          plan.PlanName,
		YookassaPaymentID: ykPayment.ID,
		Amount:            plan.Amount,
		Currency:          plan.Currency,
		Status:            mapYookassaStatus(ykPayment.Status, ykPayment.Paid),
		PaymentKind:       models.PaymentKindInitial,
		Description:       &description,
	})
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if ykPayment.Confirmation.ConfirmationURL == "" {
		logger.Error("PostCheckoutHandler_V1: ykPayment.Confirmation.ConfirmationURL is Empty")
		return httperr.InternalServerError(tr.TErr("error.payment-create-failed"))
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, CheckoutResponse{
		PaymentID:       ykPayment.ID,
		ConfirmationURL: ykPayment.Confirmation.ConfirmationURL,
	})
	return nil
}

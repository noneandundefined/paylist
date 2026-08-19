package auth_handler_v1

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"paylist.server/infra/constants"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/profanity"
	"paylist.server/pkg/security"
	"paylist.server/util"
)

/* Neosync HTTPx V1 */
/* Handler: создание пользователя */

func (h *Handler) AuthSignupHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	var payload *AuthSignupPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	valid := regexp.MustCompile(`^[\p{L}0-9_]+$`)

	var firstNameStr, lastNameStr string
	if payload.FirstName != nil {
		firstNameStr = strings.TrimSpace(*payload.FirstName)
		payload.FirstName = &firstNameStr

		if !valid.MatchString(firstNameStr) {
			return httperr.BadRequest(tr.TErr("error.invalid-characters-firstname"))
		}
	}

	if payload.LastName != nil {
		lastNameStr = strings.TrimSpace(*payload.LastName)
		payload.LastName = &lastNameStr

		if !valid.MatchString(lastNameStr) {
			return httperr.BadRequest(tr.TErr("error.invalid-characters-lastname"))
		}
	}

	if (payload.FirstName == nil || *payload.FirstName == "") && (payload.LastName == nil || *payload.LastName == "") {
		return httperr.BadRequest(tr.TErr("error.at-least-one-name-required"))
	}

	if err := profanity.Reject(ctx, tr, "signup-name", profanity.Pointers(payload.FirstName, payload.LastName)...); err != nil {
		return err
	}

	password := strings.TrimSpace(payload.Password)

	if _, chPass := constants.CheckSimplePasswords[strings.ToLower(password)]; chPass {
		return httperr.BadRequest(tr.TErr("error.simple-password"))
	}

	normalizedEmail := util.NormalizeEmail(payload.Email)
	if security.PasswordEqualsEmail(password, normalizedEmail) {
		return httperr.BadRequest(tr.TErr("error.password-equals-email"))
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("AuthSignupHandler_V1 req={%s}: Failed start tx: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, httperr.Err_DbNetwork)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	passwordHashed, err := security.HashPassword(payload.Password)
	if err != nil {
		logger.Error("AuthSignupHandler_V1 req={%s}: Failed hash password: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.InternalServerError(err.Error())
	}

	uuid := uuid.NewString()

	userCore := &models.UserCore{
		UserUUID:  uuid,
		Email:     normalizedEmail,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  passwordHashed,
	}

	if err := h.Store.Users.Create_UserCore(ctx, tx, userCore); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Users.Create_UserSubscription(ctx, tx, uuid); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		return httperr.Conflict(tr.TErr("error.failed-to-save-data"))
	}

	if _, err := h.Store.Referrals.Ensure_ReferralCode(ctx, uuid); err != nil {
		logger.Error("AuthSignupHandler_V1 req={%s}: Failed to create referral code: %s", ctx.Value("XREQID").(string), err.Error())
	}

	if err := h.Store.Referrals.Attach_Referral(ctx, payload.ReferralCode, uuid); err != nil {
		logger.Error("AuthSignupHandler_V1 req={%s}: Failed to attach referral: %s", ctx.Value("XREQID").(string), err.Error())
	}

	/* Send confirm to email */
	if err := h.sendConfirmEmail(ctx, userCore.Email, uuid, ctx.Value("XREQID").(string)); err != nil {
		return err
	}

	httpx.HttpResponseWithETag(w, r, http.StatusCreated, tr.T("confirm-link-send-to-email"))
	return nil
}

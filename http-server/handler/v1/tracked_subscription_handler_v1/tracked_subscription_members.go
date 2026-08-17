package tracked_subscription_handler_v1

import (
	"errors"
	"net/http"

	"paylist.server/infra/constants"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
	"paylist.server/util"

	"github.com/go-playground/validator"
)

func (h *Handler) GetSubscriptionMembersHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	id, err := parseSubscriptionID(tr, r)
	if err != nil {
		return err
	}

	if _, _, err := h.requireAcceptedMember(ctx, tr, uint64(id), authToken.User.UserUUID); err != nil {
		return err
	}

	members, err := h.Store.TrackedSubscriptions.Get_MembersBySubscriptionID(ctx, uint64(id))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	proposal, err := h.buildPendingProposalResponse(ctx, uint64(id), authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, TrackedSubscriptionMembersResponse{
		Members:         members,
		PendingProposal: proposal,
	})
	return nil
}

func (h *Handler) InviteSubscriptionMemberHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)
	reqID, _ := ctx.Value("XREQID").(string)

	id, err := parseSubscriptionID(tr, r)
	if err != nil {
		return err
	}

	var payload TrackedSubscriptionInvitePayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	role := constants.MemberRoleMember
	sharePercent := payload.SharePercent
	if payload.Role == constants.MemberRoleObserver {
		role = constants.MemberRoleObserver
		sharePercent = 0
	} else if err := validateSharePercent(tr, sharePercent); err != nil {
		return err
	} else if sharePercent <= 0 {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-share-invalid"))
	}

	sub, owner, err := h.requireOwnerMember(ctx, tr, uint64(id), authToken.User.UserUUID)
	if err != nil {
		return err
	}

	email := util.NormalizeEmail(payload.Email)
	if email == util.NormalizeEmail(authToken.User.Email) {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-invite-self"))
	}

	existing, err := h.Store.TrackedSubscriptions.Get_MemberByEmail(ctx, uint64(id), email)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if existing != nil && existing.Status == "accepted" {
		return httperr.Conflict(tr.TErr("error.shared-subscription-already-member"))
	}

	activeCount, err := h.Store.TrackedSubscriptions.Count_ActiveMembers(ctx, uint64(id))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if existing == nil || existing.Status == "declined" {
		if !canInviteMoreMembers(authToken, activeCount) {
			return httperr.BadRequest(tr.TErr("error.shared-subscription-member-limit"))
		}
	}

	token, expires := newInviteToken()
	ownerName := memberDisplayName(authToken.User.FirstName, authToken.User.LastName, authToken.User.Email)

	if existing != nil && existing.Status != "accepted" {
		existing.InviteToken = &token
		existing.InviteExpiresAt = &expires
		existing.Role = role
		if err := h.Store.TrackedSubscriptions.Refresh_MemberInvite(ctx, existing, sharePercent); err != nil {
			if errors.Is(err, httperr.Err_NotUpdated) {
				return httperr.BadRequest(tr.TErr("error.shared-subscription-share-exceeds-owner"))
			}

			return httperr.Db(ctx, err)
		}

		h.sendShareInviteEmail(tr, email, ownerName, sub.Name, role, sharePercent, token, reqID)
		httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.shared-subscription-invited"))
		return nil
	}

	if isPayingRole(role) && owner.SharePercent-sharePercent < 0 {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-share-exceeds-owner"))
	}

	member := &models.TrackedSubscriptionMember{
		TrackedSubscriptionID: uint64(id),
		Email:                 email,
		Role:                  role,
		SharePercent:          sharePercent,
		IncludeInAnalytics:    isPayingRole(role),
		InviteToken:           &token,
		InviteExpiresAt:       &expires,
	}

	if err := h.Store.TrackedSubscriptions.Create_MemberInvite(ctx, member); err != nil {
		if errors.Is(err, httperr.Err_NotUpdated) {
			return httperr.BadRequest(tr.TErr("error.shared-subscription-share-exceeds-owner"))
		}

		return httperr.Db(ctx, err)
	}

	h.sendShareInviteEmail(tr, email, ownerName, sub.Name, role, sharePercent, token, reqID)
	httpx.HttpResponse(w, r, http.StatusCreated, tr.T("success.shared-subscription-invited"))
	return nil
}

func (h *Handler) DeleteSubscriptionMemberHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	id, err := parseSubscriptionID(tr, r)
	if err != nil {
		return err
	}

	memberID, err := parseMuxID(tr, r, "memberId", "error.shared-subscription-member-invalid")
	if err != nil {
		return err
	}

	if _, _, err := h.requireOwnerMember(ctx, tr, uint64(id), authToken.User.UserUUID); err != nil {
		return err
	}

	member, err := h.Store.TrackedSubscriptions.Get_MemberByID(ctx, uint64(id), memberID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if member == nil {
		return httperr.NotFound(tr.TErr("error.shared-subscription-member-not-found"))
	}

	if member.Role == constants.MemberRoleOwner {
		return httperr.Forbidden(tr.TErr("error.shared-subscription-owner-only"))
	}

	if err := h.Store.TrackedSubscriptions.Delete_MemberAndReturnShare(ctx, uint64(id), member.ID, member.SharePercent); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.shared-subscription-member-removed"))
	return nil
}

func (h *Handler) LeaveSubscriptionHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	id, err := parseSubscriptionID(tr, r)
	if err != nil {
		return err
	}

	_, member, err := h.requireAcceptedMember(ctx, tr, uint64(id), authToken.User.UserUUID)
	if err != nil {
		return err
	}

	if member.Role == constants.MemberRoleOwner {
		return httperr.Forbidden(tr.TErr("error.shared-subscription-owner-cannot-leave"))
	}

	if err := h.Store.TrackedSubscriptions.Delete_MemberAndReturnShare(ctx, uint64(id), member.ID, member.SharePercent); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.shared-subscription-left"))
	return nil
}

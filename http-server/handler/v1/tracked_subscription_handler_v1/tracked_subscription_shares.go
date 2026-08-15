package tracked_subscription_handler_v1

import (
	"net/http"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"

	"github.com/go-playground/validator"
)

func (h *Handler) ProposeSubscriptionSharesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	id, err := parseSubscriptionID(tr, r)
	if err != nil {
		return err
	}

	var payload TrackedSubscriptionShareProposalPayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	if _, _, err := h.requireAcceptedMember(ctx, tr, uint64(id), authToken.User.UserUUID); err != nil {
		return err
	}

	members, err := h.Store.TrackedSubscriptions.Get_MembersBySubscriptionID(ctx, uint64(id))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	accepted := make([]models.TrackedSubscriptionMember, 0)
	acceptedIDs := map[uint64]struct{}{}
	for _, member := range members {
		if member.Status == "pending" {
			return httperr.Conflict(tr.TErr("error.shared-subscription-pending-invite"))
		}

		if member.Status != "accepted" {
			continue
		}

		accepted = append(accepted, member)
		acceptedIDs[member.ID] = struct{}{}
	}

	if len(payload.Shares) != len(accepted) {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-shares-incomplete"))
	}

	items := make([]models.TrackedSubscriptionShareProposalItem, 0, len(payload.Shares))
	seen := map[uint64]struct{}{}
	for _, share := range payload.Shares {
		if err := validateSharePercent(tr, share.SharePercent); err != nil {
			return err
		}

		if _, ok := acceptedIDs[share.MemberID]; !ok {
			return httperr.BadRequest(tr.TErr("error.shared-subscription-shares-incomplete"))
		}

		if _, ok := seen[share.MemberID]; ok {
			return httperr.BadRequest(tr.TErr("error.shared-subscription-shares-incomplete"))
		}

		seen[share.MemberID] = struct{}{}
		items = append(items, models.TrackedSubscriptionShareProposalItem{
			MemberID:     share.MemberID,
			SharePercent: share.SharePercent,
		})
	}

	if !sharesSumToHundred(items) {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-shares-sum"))
	}

	if len(accepted) == 1 {
		proposal, err := h.Store.TrackedSubscriptions.Create_ShareProposal(ctx, uint64(id), authToken.User.UserUUID, items)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if err := h.Store.TrackedSubscriptions.Apply_ShareProposal(ctx, proposal.ID, items); err != nil {
			return httperr.Db(ctx, err)
		}

		httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.shared-subscription-shares-updated"))
		return nil
	}

	if _, err := h.Store.TrackedSubscriptions.Create_ShareProposal(ctx, uint64(id), authToken.User.UserUUID, items); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.shared-subscription-shares-proposed"))
	return nil
}

func (h *Handler) VoteSubscriptionSharesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	id, err := parseSubscriptionID(tr, r)
	if err != nil {
		return err
	}

	proposalID, err := parseMuxID(tr, r, "proposalId", "error.shared-subscription-proposal-invalid")
	if err != nil {
		return err
	}

	var payload TrackedSubscriptionShareVotePayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if _, _, err := h.requireAcceptedMember(ctx, tr, uint64(id), authToken.User.UserUUID); err != nil {
		return err
	}

	proposal, err := h.Store.TrackedSubscriptions.Get_ShareProposalByID(ctx, uint64(id), proposalID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if proposal == nil || proposal.Status != "pending" {
		return httperr.NotFound(tr.TErr("error.shared-subscription-proposal-not-found"))
	}

	if !payload.Accept {
		if err := h.Store.TrackedSubscriptions.Reject_ShareProposal(ctx, proposal.ID); err != nil {
			return httperr.Db(ctx, err)
		}

		httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.shared-subscription-shares-rejected"))
		return nil
	}

	if err := h.Store.TrackedSubscriptions.Upsert_ShareVote(ctx, proposal.ID, authToken.User.UserUUID, true); err != nil {
		return httperr.Db(ctx, err)
	}

	members, err := h.Store.TrackedSubscriptions.Get_MembersBySubscriptionID(ctx, uint64(id))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	votes, err := h.Store.TrackedSubscriptions.Get_ShareProposalVotes(ctx, proposal.ID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if !allAcceptedMembersAgreed(members, votes) {
		httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.shared-subscription-shares-voted"))
		return nil
	}

	items, err := h.Store.TrackedSubscriptions.Get_ShareProposalItems(ctx, proposal.ID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.TrackedSubscriptions.Apply_ShareProposal(ctx, proposal.ID, items); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("success.shared-subscription-shares-updated"))
	return nil
}

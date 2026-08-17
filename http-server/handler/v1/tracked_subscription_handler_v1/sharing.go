package tracked_subscription_handler_v1

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"paylist.server/infra/constants"
	"paylist.server/infra/locale"
	"paylist.server/infra/logger"
	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/premium"
	"paylist.server/types"

	"github.com/google/uuid"
)

func (h *Handler) requireAcceptedMember(ctx context.Context, tr locale.Translator, subscriptionID uint64, userUUID string) (*models.TrackedSubscription, *models.TrackedSubscriptionMember, error) {
	sub, err := h.Store.TrackedSubscriptions.Get_SubscriptionById(ctx, subscriptionID, userUUID)
	if err != nil {
		return nil, nil, httperr.Db(ctx, err)
	}

	if sub == nil {
		return nil, nil, httperr.NotFound(tr.TErr("error.tracked-subscription-not-found"))
	}

	member, err := h.Store.TrackedSubscriptions.Get_AcceptedMember(ctx, subscriptionID, userUUID)
	if err != nil {
		return nil, nil, httperr.Db(ctx, err)
	}

	if member == nil {
		return nil, nil, httperr.NotFound(tr.TErr("error.tracked-subscription-not-found"))
	}

	return sub, member, nil
}

func (h *Handler) requireOwnerMember(ctx context.Context, tr locale.Translator, subscriptionID uint64, userUUID string) (*models.TrackedSubscription, *models.TrackedSubscriptionMember, error) {
	sub, member, err := h.requireAcceptedMember(ctx, tr, subscriptionID, userUUID)
	if err != nil {
		return nil, nil, err
	}

	if member.Role != constants.MemberRoleOwner {
		return nil, nil, httperr.Forbidden(tr.TErr("error.shared-subscription-owner-only"))
	}

	return sub, member, nil
}

func isPayingRole(role string) bool {
	return role != constants.MemberRoleObserver
}

func canInviteMoreMembers(authToken *types.AuthToken, activeCount int) bool {
	if premium.IsPremiumPlan(authToken) {
		return true
	}

	return activeCount < constants.FreeMaxSharedMembers
}

func validateSharePercent(tr locale.Translator, value float64) error {
	if value < 0 || value > 100 || math.IsNaN(value) || math.IsInf(value, 0) {
		return httperr.BadRequest(tr.TErr("error.shared-subscription-share-invalid"))
	}

	return nil
}

func sharesSumToHundred(items []models.TrackedSubscriptionShareProposalItem) bool {
	var total float64
	for _, item := range items {
		total += item.SharePercent
	}

	return math.Abs(total-100) < 0.01
}

func formatSharePercent(value float64) string {
	if math.Abs(value-math.Round(value)) < 0.0005 {
		return fmt.Sprintf("%.0f%%", math.Round(value))
	}

	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.3f", value), "0"), ".") + "%"
}

func memberDisplayName(firstName, lastName *string, email string) string {
	parts := make([]string, 0, 2)
	if firstName != nil && strings.TrimSpace(*firstName) != "" {
		parts = append(parts, strings.TrimSpace(*firstName))
	}

	if lastName != nil && strings.TrimSpace(*lastName) != "" {
		parts = append(parts, strings.TrimSpace(*lastName))
	}

	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}

	return email
}

func newInviteToken() (string, time.Time) {
	expires := time.Now().UTC().Add(constants.ShareInviteTTL)
	return uuid.NewString(), expires
}

func inviteLink(token string) string {
	clientURL := strings.TrimRight(strings.TrimSpace(os.Getenv("CLIENT_URL")), "/")
	return fmt.Sprintf("%s/paylist-subscription-invite?token=%s", clientURL, token)
}

func (h *Handler) sendShareInviteEmail(tr locale.Translator, to, ownerName, subscriptionName, role string, sharePercent float64, token, reqID string) {
	go func() {
		link := inviteLink(token)
		intro := fmt.Sprintf(tr.T("invite-email-intro-observer"), ownerName, subscriptionName)
		if isPayingRole(role) {
			intro = fmt.Sprintf(tr.T("invite-email-intro"), ownerName, subscriptionName, formatSharePercent(sharePercent))
		}

		if err := pkg.SendEmail(to, tr.T("invite-email-title"), fmt.Sprintf(`
				<p>%s</p>

				<div style="border-left:4px solid #0085FF; padding:4px 0 4px 16px; margin:16px 0;">
					<a href="%s" style="color:#0085FF; text-decoration:underline; font-weight:bold;">%s</a>
				</div>

				<p>%s</p>
			`,
			intro,
			link,
			tr.T("invite-email-cta"),
			tr.T("invite-email-expiry"),
		), tr); err != nil {
			logger.Error("sendShareInviteEmail req={%s}: %s", reqID, err.Error())
		}
	}()
}

func (h *Handler) buildPendingProposalResponse(ctx context.Context, subscriptionID uint64, currentUserUUID string) (*TrackedSubscriptionShareProposalResponse, error) {
	proposal, err := h.Store.TrackedSubscriptions.Get_PendingShareProposal(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	if proposal == nil {
		return nil, nil
	}

	items, err := h.Store.TrackedSubscriptions.Get_ShareProposalItems(ctx, proposal.ID)
	if err != nil {
		return nil, err
	}

	votes, err := h.Store.TrackedSubscriptions.Get_ShareProposalVotes(ctx, proposal.ID)
	if err != nil {
		return nil, err
	}

	var myVote *bool
	for _, vote := range votes {
		if vote.UserUUID == currentUserUUID {
			accepted := vote.Accepted
			myVote = &accepted
			break
		}
	}

	return &TrackedSubscriptionShareProposalResponse{
		ID:                 proposal.ID,
		ProposedByUserUUID: proposal.ProposedByUserUUID,
		Status:             proposal.Status,
		Items:              items,
		Votes:              votes,
		MyVote:             myVote,
	}, nil
}

func allAcceptedMembersAgreed(members []models.TrackedSubscriptionMember, votes []models.TrackedSubscriptionShareVote) bool {
	acceptedBy := map[string]bool{}
	for _, vote := range votes {
		if vote.Accepted {
			acceptedBy[vote.UserUUID] = true
		}
	}

	for _, member := range members {
		if member.Status != "accepted" || member.UserUUID == nil || !isPayingRole(member.Role) {
			continue
		}

		if !acceptedBy[*member.UserUUID] {
			return false
		}
	}

	return true
}

package user_handler_v1

import (
	"net/http"
	"os"
	"strings"

	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/pkg/referral"
	"paylist.server/types"
)

func (h *Handler) UserReferralGetHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	code, err := h.Store.Referrals.Ensure_ReferralCode(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	count, err := h.Store.Referrals.Count_ConvertedReferrals(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	current := referral.RankForCount(count)
	ranks := make([]UserReferralRankResponse, 0, len(referral.Ranks))
	for _, rank := range referral.Ranks {
		ranks = append(ranks, UserReferralRankResponse{
			Level:      rank.Level,
			Key:        rank.Key,
			MinCount:   rank.MinCount,
			MaxCount:   rank.MaxCount,
			RewardDays: rank.RewardDays,
			Current:    rank.Level == current.Level,
		})
	}

	appURL := strings.TrimRight(strings.TrimSpace(os.Getenv("CLIENT_URL")), "/")
	if !strings.HasPrefix(strings.ToLower(appURL), "https://") {
		appURL = "https://paylist.site"
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, UserReferralResponse{
		Code:          code,
		SiteURL:       referral.SiteURL(appURL, code),
		BotURL:        referral.BotURL(os.Getenv("TELEGRAM_BOT_USERNAME"), code),
		ReferralCount: count,
		Rank:          current.Level,
		Ranks:         ranks,
	})
	return nil
}

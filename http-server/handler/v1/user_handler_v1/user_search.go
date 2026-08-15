package user_handler_v1

import (
	"net/http"
	"strings"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	"paylist.server/types"
	"paylist.server/util"
)

func (h *Handler) UserSearchHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	emailQuery := util.NormalizeEmail(r.URL.Query().Get("email"))
	if !strings.Contains(emailQuery, "@") || !strings.Contains(emailQuery, ".") {
		httpx.HttpResponseWithETag(w, r, http.StatusOK, []models.UserPublicProfile{})
		return nil
	}

	users, err := h.Store.Users.Search_UserPublicProfilesByEmail(ctx, emailQuery, authToken.User.UserUUID, 1)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if users == nil {
		users = []models.UserPublicProfile{}
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, users)
	return nil
}

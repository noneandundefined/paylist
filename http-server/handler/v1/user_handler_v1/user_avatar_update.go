package user_handler_v1

import (
	"bytes"
	"io"
	"net/http"
	"path"

	"github.com/google/uuid"
	"paylist.server/infra/logger"
	"paylist.server/middleware"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
	pkgs3 "paylist.server/pkg/s3"
	"paylist.server/types"
)

const (
	avatarFormField = "avatar"
	avatarMaxBytes  = 2 << 20
)

var avatarContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

func (h *Handler) UserAvatarUpdateHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if err := r.ParseMultipartForm(avatarMaxBytes); err != nil {
		return httperr.BadRequest(tr.TErr("error.avatar-too-large"))
	}

	file, header, err := r.FormFile(avatarFormField)
	if err != nil {
		return httperr.BadRequest(tr.TErr("error.avatar-invalid"))
	}
	defer file.Close()

	if header.Size > avatarMaxBytes {
		return httperr.BadRequest(tr.TErr("error.avatar-too-large"))
	}

	payload, err := io.ReadAll(io.LimitReader(file, avatarMaxBytes+1))
	if err != nil {
		return httperr.BadRequest(tr.TErr("error.avatar-invalid"))
	}

	if len(payload) == 0 || int64(len(payload)) > avatarMaxBytes {
		return httperr.BadRequest(tr.TErr("error.avatar-too-large"))
	}

	contentType := http.DetectContentType(payload)
	ext, ok := avatarContentTypes[contentType]
	if !ok {
		return httperr.BadRequest(tr.TErr("error.avatar-invalid"))
	}

	s3Client, err := pkgs3.New()
	if err != nil {
		logger.Error("UserAvatarUpdateHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.InternalServerError(tr.TErr("error.s3-upload-failed"))
	}

	user, err := h.Store.Users.Get_UserCoreByUserUuid(ctx, authToken.User.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	key := path.Join("avatars", authToken.User.UserUUID, uuid.NewString()+ext)
	publicURL, err := s3Client.Upload(ctx, key, contentType, bytes.NewReader(payload))
	if err != nil {
		logger.Error("UserAvatarUpdateHandler_V1 req={%s}: upload failed: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.InternalServerError(tr.TErr("error.s3-upload-failed"))
	}

	if err := h.Store.Users.Update_UserAvatar(ctx, authToken.User.UserUUID, publicURL); err != nil {
		_ = s3Client.Delete(ctx, key)
		return httperr.Db(ctx, err)
	}

	if user != nil && user.Avatars != nil {
		if oldKey := s3Client.KeyFromURL(*user.Avatars); oldKey != "" && oldKey != key {
			_ = s3Client.Delete(ctx, oldKey)
		}
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("avatar-updated"))
	return nil
}

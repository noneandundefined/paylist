package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"paylist.server/infra/constants"
)

func GenerateVerificationEmailLink(userUuid string) string {
	expires := time.Now().Add(constants.EMAIL_LINK_TIMEOUT).Unix()

	payload := fmt.Sprintf("%s:%d", userUuid, expires)

	mac := hmac.New(sha256.New, []byte(os.Getenv("SUPER_SECRET_KEY")))
	mac.Write([]byte(payload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	values := url.Values{}
	values.Set("uuid", userUuid)
	values.Set("exp", fmt.Sprintf("%d", expires))
	values.Set("sig", signature)

	return fmt.Sprintf("%s/paylist-confirm-email?%s", os.Getenv("CLIENT_URL"), values.Encode())
}

func ValidateEmailConfirmLink(userUuid, expStr, sig string) bool {
	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return false
	}

	if time.Now().Unix() > exp {
		return false
	}

	payload := fmt.Sprintf("%s:%d", userUuid, exp)

	mac := hmac.New(sha256.New, []byte(os.Getenv("SUPER_SECRET_KEY")))
	mac.Write([]byte(payload))

	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(sig), []byte(expectedSig))
}

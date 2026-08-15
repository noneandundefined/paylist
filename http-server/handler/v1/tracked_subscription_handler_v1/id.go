package tracked_subscription_handler_v1

import (
	"net/http"
	"strconv"

	"paylist.server/infra/locale"
	"paylist.server/pkg/httpx/httperr"

	"github.com/gorilla/mux"
)

func parseSubscriptionID(tr locale.Translator, r *http.Request) (int, error) {
	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		return 0, httperr.NotFound(tr.TErr("error.tracked-subscription-not-found"))
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return 0, httperr.BadRequest(tr.TErr("error.tracked-subscription-invalid-id"))
	}

	return id, nil
}

func parseMuxID(tr locale.Translator, r *http.Request, key, invalidKey string) (uint64, error) {
	raw := mux.Vars(r)[key]
	if raw == "" {
		return 0, httperr.BadRequest(tr.TErr(invalidKey))
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		return 0, httperr.BadRequest(tr.TErr(invalidKey))
	}

	return value, nil
}

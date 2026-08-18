package country_handler_v1

import (
	"net/http"
	"strings"

	"paylist.server/middleware"
	"paylist.server/pkg/country"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
)

func (h *Handler) GetCountriesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	items := country.GetCountries()

	httpx.HttpCache(w, 86400)
	httpx.HttpResponseWithETag(w, r, http.StatusOK, items)
	return nil
}

func (h *Handler) GetCountryInflationHandler_V1(w http.ResponseWriter, r *http.Request) error {
	tr := middleware.TranslatorFromContext(r.Context())
	code := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("code")))

	if code == "" {
		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	item, found := country.FindCountry(code)
	if !found {
		httpx.HttpResponse(w, r, http.StatusOK, map[string]any{
			"code":           code,
			"inflation_rate": country.GetInflationRate(code),
			"estimated":      true,
		})
		return nil
	}

	httpx.HttpResponse(w, r, http.StatusOK, item)
	return nil
}

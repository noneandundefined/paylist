package country_handler_v1

import (
	"net/http"

	"paylist.server/pkg/httpx"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	countryRouter := router.PathPrefix("/country").Subrouter()

	countryRouter.Handle("/countries", httpx.ErrorHandler(h.GetCountriesHandler_V1)).Methods(http.MethodGet)
	countryRouter.Handle("/inflation", httpx.ErrorHandler(h.GetCountryInflationHandler_V1)).Methods(http.MethodGet)
}

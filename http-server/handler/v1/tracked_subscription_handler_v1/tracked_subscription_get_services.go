package tracked_subscription_handler_v1

import (
	"net/http"

	"paylist.server/infra/store/postgres/models"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
)

type ServiceResponse struct {
	ID       uint64   `json:"id"`
	Slug     string   `json:"slug"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Aliases  []string `json:"aliases"`
}

func (h *Handler) GetServicesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	services, err := h.Store.TrackedSubscriptions.Get_Services(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, mapServicesResponse(services))
	return nil
}

func mapServicesResponse(services []models.Service) []ServiceResponse {
	response := make([]ServiceResponse, 0, len(services))

	for _, service := range services {
		aliases := []string(service.Aliases)
		if aliases == nil {
			aliases = []string{}
		}

		response = append(response, ServiceResponse{
			ID:       service.ID,
			Slug:     service.Slug,
			Name:     service.Name,
			Category: service.Category,
			Aliases:  aliases,
		})
	}

	return response
}

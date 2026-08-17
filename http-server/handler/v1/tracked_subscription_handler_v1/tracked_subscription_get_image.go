package tracked_subscription_handler_v1

import (
	"context"
	"net/http"
	"os"

	"paylist.server/middleware"
	"paylist.server/pkg/analytics"
	"paylist.server/pkg/httpx/httperr"
)

var catalogImageExts = []string{".png", ".jpg", ".jpeg", ".webp"}

func (h *Handler) GetSubscriptionImageHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	name := r.URL.Query().Get("name")
	if name == "" {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-image-not-found"))
	}

	img := h.lookupCatalogImage(ctx, name)
	if img == "" {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-image-not-found"))
	}

	filePath := "media/subscriptions/" + img
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return httperr.NotFound(tr.TErr("error.tracked-subscription-image-not-found"))
	}

	http.ServeFile(w, r, filePath)
	return nil
}

func (h *Handler) lookupCatalogImage(ctx context.Context, name string) string {
	services, err := h.Store.TrackedSubscriptions.Get_Services(ctx)
	if err != nil || len(services) == 0 {
		return ""
	}

	catalog := make([]analytics.Service, 0, len(services))
	for _, service := range services {
		catalog = append(catalog, analytics.Service{
			Slug:     service.Slug,
			Name:     service.Name,
			Category: service.Category,
			Aliases:  service.Aliases,
		})
	}

	matched, ok := analytics.MatchService(catalog, name)
	if !ok || matched.Slug == "" {
		return ""
	}

	for _, ext := range catalogImageExts {
		filename := matched.Slug + ext
		if _, err := os.Stat("media/subscriptions/" + filename); err == nil {
			return filename
		}
	}

	return ""
}

package api

import (
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store"
	"github.com/go-chi/chi/v5"
)

type GatewayAPI struct {
	store store.Store
}

func NewGatewayAPI() (*GatewayAPI, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}
	return &GatewayAPI{
		store: store,
	}, nil
}

func (gapi *GatewayAPI) Bind(root *chi.Mux) error {
	root.Route("/gateways", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/owned/{owner:(?i)(0x)?[0-9a-f]{40}}", gapi.OwnedGateways)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}/", gapi.GatewayDetailsByID)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}/events", gapi.GatewayEventsByID)
			r.Route("/events", func(r chi.Router) {
				r.Post("/owner/{owner:(?i)(0x)?[0-9a-f]{40}}/pending", gapi.PendingGatewayEvents)
			})
			r.Get("/frequencyplans", gapi.SupportedFrequencyPlans)
			r.Get("/map/res0", gapi.GatewayMapRes0)
			r.Get("/map/{hex:(?i)[0-9a-f]{15}}", gapi.GatewayMap)
		})
	})

	return nil
}

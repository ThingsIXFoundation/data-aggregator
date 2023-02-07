package api

import (
	"github.com/ThingsIXFoundation/data-aggregator/mapper/store"
	"github.com/go-chi/chi/v5"
)

type MapperAPI struct {
	store store.Store
}

func NewMapperAPI() (*MapperAPI, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}
	return &MapperAPI{
		store: store,
	}, nil
}

func (mapi *MapperAPI) Bind(root *chi.Mux) error {
	root.Route("/mappers", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/owned/{owner:(?i)(0x)?[0-9a-f]{40}}", mapi.OwnedMappers)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}/", mapi.MapperDetailsByID)
			r.Get("/{id:(?i)(0x)?[0-9a-f]{64}}/events", mapi.MapperEventsByID)
			r.Route("/events", func(r chi.Router) {
				r.Post("/owner/{owner:(?i)(0x)?[0-9a-f]{40}}/pending", mapi.PendingMapperEvents)
			})
		})
	})

	return nil
}

package api

import (
	"github.com/ThingsIXFoundation/data-aggregator/router/store"
	"github.com/go-chi/chi/v5"
)

type RouterAPI struct {
	store store.Store
}

func NewRouterAPI() (*RouterAPI, error) {
	store, err := store.NewStore()
	if err != nil {
		return nil, err
	}
	return &RouterAPI{
		store: store,
	}, nil
}

func (rapi *RouterAPI) Bind(root *chi.Mux) error {
	root.Route("/routers", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/snapshot", rapi.Snapshot)
		})
	})

	return nil
}

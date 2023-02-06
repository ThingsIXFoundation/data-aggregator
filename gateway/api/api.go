package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/gateway/store"
	httputils "github.com/ThingsIXFoundation/http-utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func (gapi *GatewayAPI) Serve(ctx context.Context) chan error {
	root := chi.NewRouter()

	httputils.BindStandardMiddleware(root)
	//root.Use(cache.DisableCacheOnGetRequests)

	root.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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

	srv := http.Server{
		Handler:      root,
		Addr:         viper.GetString(config.CONFIG_GATEWAY_API_HTTP_LISTEN_ADDRESS),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	stopped := make(chan error)
	go func() {
		logrus.WithField("addr", srv.Addr).Info("start HTTP API service")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Fatal("HTTP service crashed")
		}
	}()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		stopped <- srv.Shutdown(ctx)
	}()

	return stopped
}

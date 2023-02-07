package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	gatewayapi "github.com/ThingsIXFoundation/data-aggregator/gateway/api"
	mapperapi "github.com/ThingsIXFoundation/data-aggregator/mapper/api"
	routerapi "github.com/ThingsIXFoundation/data-aggregator/router/api"
	httputils "github.com/ThingsIXFoundation/http-utils"
	"github.com/ThingsIXFoundation/http-utils/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type API struct {
	gatewayAPI *gatewayapi.GatewayAPI
	routerAPI  *routerapi.RouterAPI
	mapperAPI  *mapperapi.MapperAPI
}

func NewAPI() (*API, error) {
	gatewayAPI, err := gatewayapi.NewGatewayAPI()
	if err != nil {
		return nil, err
	}

	mapperAPI, err := mapperapi.NewMapperAPI()
	if err != nil {
		return nil, err
	}

	routerAPI, err := routerapi.NewRouterAPI()
	if err != nil {
		return nil, err
	}

	return &API{
		gatewayAPI: gatewayAPI,
		routerAPI:  routerAPI,
		mapperAPI:  mapperAPI,
	}, nil

}

func (a *API) Serve(ctx context.Context) chan error {
	root := chi.NewRouter()

	httputils.BindStandardMiddleware(root)
	root.Use(cache.DisableCacheOnGetRequests)

	root.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	srv := http.Server{
		Handler:      root,
		Addr:         viper.GetString(config.CONFIG_API_HTTP_LISTEN_ADDRESS),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.gatewayAPI.Bind(root)
	a.routerAPI.Bind(root)
	a.mapperAPI.Bind(root)

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

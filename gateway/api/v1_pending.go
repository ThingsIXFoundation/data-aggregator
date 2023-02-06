package api

import (
	"context"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
)

func (gapi *GatewayAPI) PendingGatewayEvents(w http.ResponseWriter, r *http.Request) {
	var (
		log           = logging.WithContext(r.Context())
		ctx, cancel   = context.WithTimeout(r.Context(), 15*time.Second)
		owner         = common.HexToAddress(chi.URLParam(r, "owner"))
		filterRequest []types.GatewayEventType
	)
	defer cancel()

	if r.ContentLength > 0 {
		if err := encoding.DecodeHTTPJSONBody(w, r, &filterRequest); err != nil {
			log.WithError(err).Error("unable to decode search gateway request")
			http.Error(w, err.Msg, err.Status)
			return
		}
	}

	events, err := gapi.store.PendingEventsForOwner(ctx, owner)
	if err != nil {
		log.WithError(err).Error("error while getting pending events for owner")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	events = utils.Filter(events, func(event *types.GatewayEvent) bool {
		if len(filterRequest) > 0 && utils.In(filterRequest, event.Type) {
			return true
		} else {
			return true
		}
	})

	syncedTo, err := gapi.store.CurrentBlock(ctx, "GatewayIngestor")
	if err != nil {
		log.WithError(err).Error("error while getting pending events for owner")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, &PendingGatewayEventsResponse{
		Confirmations: viper.GetUint64(config.CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS),
		SyncedTo:      syncedTo,
		Events:        events,
	})
}

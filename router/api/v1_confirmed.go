package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/spf13/viper"
)

// Snapshot returns the registed routers from cache.
func (rapi *RouterAPI) Snapshot(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
	)
	defer cancel()

	routers, err := rapi.store.GetAll(ctx)
	if err != nil {
		log.WithError(err).Error("error while getting routers")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	currentBlock, err := rapi.store.CurrentBlock(ctx, "RouterAggregator")
	if err != nil {
		log.WithError(err).Error("error while sync state")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// got router info, cache it for fast returning
	reply, err := json.Marshal(map[string]interface{}{
		"blockNumber": currentBlock,
		"chainId":     viper.GetUint64(config.CONFIG_CHAINSYNC_CHAINID),
		"routers":     routers,
	})
	if err != nil {
		log.WithError(err).Error("error while getting routers")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "public, max-age=900")
	w.WriteHeader(http.StatusOK)
	w.Write(reply)

}

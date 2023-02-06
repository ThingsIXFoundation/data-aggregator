package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
)

func replyGatewaysPaged(gateways []*types.Gateway, pageSize int, w http.ResponseWriter, r *http.Request) {
	if gateways == nil {
		gateways = make([]*types.Gateway, 0) // prevent null in reply
	}

	// trim results if there are more than requested
	moreAvailable := len(gateways) > pageSize
	if moreAvailable {
		gateways = gateways[:pageSize]
	}

	encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
		"moreAvailable": moreAvailable,
		"gateways":      gateways,
	})
}

func replyEventsPaged(events []*types.GatewayEvent, pageSize int, w http.ResponseWriter, r *http.Request) {
	if events == nil {
		events = make([]*types.GatewayEvent, 0) // prevent null in reply
	}

	// trim results if there are more than requested
	moreAvailable := len(events) > pageSize
	if moreAvailable {
		events = events[:pageSize]
	}

	encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
		"moreAvailable": moreAvailable,
		"events":        events,
	})
}

func (gapi *GatewayAPI) OwnedGateways(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		page, _     = strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		owner       = common.HexToAddress(chi.URLParam(r, "owner"))
	)
	defer cancel()

	if pageSize == 0 {
		pageSize = 15
	}

	gateways, err := gapi.store.GetByOwner(ctx, owner)
	if err != nil {
		log.WithError(err).Error("unable to retrieve gateways from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	gateways = utils.Paginate(gateways, page, pageSize, 1)

	replyGatewaysPaged(gateways, pageSize, w, r)
}

func (gapi *GatewayAPI) GatewayDetailsByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		gatewayID   = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	gateway, err := gapi.store.Get(ctx, gatewayID)
	if err != nil {
		log.WithError(err).Error("error while getting gateway details")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if gateway == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, gateway)
}

func (gapi *GatewayAPI) GatewayEventsByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		page, _     = strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		gatewayID   = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	if pageSize == 0 {
		pageSize = 15
	}

	events, err := gapi.store.GetEvents(ctx, gatewayID)
	if err != nil {
		log.WithError(err).Error("error while getting gateway events")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	events = utils.Paginate(events, page, pageSize, 1)

	replyEventsPaged(events, pageSize, w, r)
}

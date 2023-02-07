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

func replyMappersPaged(mappers []*types.Mapper, pageSize int, w http.ResponseWriter, r *http.Request) {
	if mappers == nil {
		mappers = make([]*types.Mapper, 0) // prevent null in reply
	}

	// trim results if there are more than requested
	moreAvailable := len(mappers) > pageSize
	if moreAvailable {
		mappers = mappers[:pageSize]
	}

	encoding.ReplyJSON(w, r, http.StatusOK, map[string]interface{}{
		"moreAvailable": moreAvailable,
		"mappers":       mappers,
	})
}

func replyEventsPaged(events []*types.MapperEvent, pageSize int, w http.ResponseWriter, r *http.Request) {
	if events == nil {
		events = make([]*types.MapperEvent, 0) // prevent null in reply
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

func (mapi *MapperAPI) OwnedMappers(w http.ResponseWriter, r *http.Request) {
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

	mappers, err := mapi.store.GetByOwner(ctx, owner)
	if err != nil {
		log.WithError(err).Error("unable to retrieve mappers from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	mappers = utils.Paginate(mappers, page, pageSize, 1)

	replyMappersPaged(mappers, pageSize, w, r)
}

func (mapi *MapperAPI) MapperDetailsByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mapperID    = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	mapper, err := mapi.store.Get(ctx, mapperID)
	if err != nil {
		log.WithError(err).Error("error while getting mapper details")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if mapper == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, mapper)
}

func (mapi *MapperAPI) MapperEventsByID(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		page, _     = strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
		mapperID    = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	if pageSize == 0 {
		pageSize = 15
	}

	events, err := mapi.store.GetEvents(ctx, mapperID)
	if err != nil {
		log.WithError(err).Error("error while getting mapper events")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	events = utils.Paginate(events, page, pageSize, 1)

	replyEventsPaged(events, pageSize, w, r)
}

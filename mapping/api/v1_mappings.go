package api

import (
	"context"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
)

func (mapi *MappingAPI) GetMappingById(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mappingID   = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	mappingRecord, err := mapi.store.GetMapping(ctx, mappingID)
	if err != nil {
		log.WithError(err).Error("unable to retrieve mapping-record from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, mappingRecord)
	return
}

func (mapi *MappingAPI) GetRecentMappingsForMapper(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mapperID    = utils.IDFromRequest(r, "id")
		since       = 24 * time.Hour
	)
	defer cancel()

	recentMappingRecords, err := mapi.store.GetRecentMappingsForMapper(ctx, mapperID, since)
	if err != nil {
		log.WithError(err).Error("unable to retrieve mapping-record from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, recentMappingRecords)
	return
}

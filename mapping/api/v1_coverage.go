package api

import (
	"context"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
)

func (mapi *MappingAPI) MinMaxCoverageDates(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
	)
	defer cancel()

	min, max, err := mapi.store.GetMinMaxCoverageDates(ctx)
	if err != nil {
		log.WithError(err).Error("error while getting min max coverage date")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ret := &MinMaxCoverageDates{
		Min: min.Format(time.DateOnly),
		Max: max.Format(time.DateOnly),
	}

	encoding.ReplyJSON(w, r, http.StatusOK, ret)
}

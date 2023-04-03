package api

import (
	"context"
	"net/http"
	"time"

	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/go-chi/chi/v5"
	"github.com/uber/h3-go/v4"
)

const MAP_ASSUMED_COVERAGE_MAX_RES = 8
const MAP_COVERAGE_MIN_RES = 6
const MAP_COVERAGE_MAX_RES = 8

func (mapi *MappingAPI) AssumedCoverageMapRes0(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 5*time.Minute)
		date        = chi.URLParam(r, "date")
	)
	defer cancel()

	at, err := time.Parse(time.DateOnly, date)
	if err != nil {
		log.Warnf("invalid date provided: %s", date)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	var coverageLocations []h3light.Cell

	for _, res0 := range h3.Res0Cells() {
		newLocations, err := mapi.store.GetAssumedCoverageLocationsInRegionAtWithRes(ctx, h3light.Cell(res0), at, 6)
		if err != nil {
			log.WithError(err).Error("error while getting coverage locations")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		coverageLocations = append(coverageLocations, newLocations...)
	}

	ret := &AssumedCoverageHexContainer{Hexes: coverageLocations}
	encoding.ReplyJSON(w, r, http.StatusOK, ret)
}

func (mapi *MappingAPI) AssumedCoverageMap(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 1*time.Minute)
		date        = chi.URLParam(r, "date")
		hex         = chi.URLParam(r, "hex")
	)
	defer cancel()

	at, err := time.Parse(time.DateOnly, date)
	if err != nil {
		log.Warnf("invalid date provided: %s", date)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	hexCell, err := h3light.CellFromString(hex)
	if err != nil {
		log.Warnf("invalid h3 index provided: %s", hex)
		http.Error(w, "invalid h3 index", http.StatusBadRequest)
		return
	}

	res := hexCell.Resolution()
	if res > MAP_ASSUMED_COVERAGE_MAX_RES {
		log.Warnf("invalid h3 resolution: %d", res)
		http.Error(w, "invalid h3 resolution", http.StatusBadRequest)
		return
	}

	//var coverageLocations []h3.Cell
	coverageLocations, err := mapi.store.GetAssumedCoverageLocationsInRegionAtWithRes(ctx, hexCell, at, 8)
	if err != nil {
		log.WithError(err).Error("error while getting coverage locations")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ret := &AssumedCoverageHexContainer{Hexes: coverageLocations}
	encoding.ReplyJSON(w, r, http.StatusOK, ret)

}

func (mapi *MappingAPI) CoverageMap(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 1*time.Minute)
		date        = chi.URLParam(r, "date")
		hex         = chi.URLParam(r, "hex")
	)
	defer cancel()

	at, err := time.Parse(time.DateOnly, date)
	if err != nil {
		log.Warnf("invalid date provided: %s", date)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	hexCell, err := h3light.CellFromString(hex)
	if err != nil {
		log.Warnf("invalid h3 index provided: %s", hex)
		http.Error(w, "invalid h3 index", http.StatusBadRequest)
		return
	}

	res := hexCell.Resolution()
	if res > MAP_COVERAGE_MAX_RES || res < MAP_COVERAGE_MIN_RES {
		log.Warnf("invalid h3 resolution: %d", res)
		http.Error(w, "invalid h3 resolution", http.StatusBadRequest)
		return
	}

	chs, err := mapi.store.GetCoverageInRegionAt(ctx, hexCell, at)
	if err != nil {
		log.WithError(err).Error("error while getting coverage locations")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	chc := &CoverageHexContainer{
		Hexes: chs,
	}

	encoding.ReplyJSON(w, r, http.StatusOK, chc)
}

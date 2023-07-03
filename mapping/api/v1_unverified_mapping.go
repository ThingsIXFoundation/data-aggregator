package api

import (
	"context"
	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/utils"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/types"
	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"google.golang.org/protobuf/encoding/protojson"
)

func (mapi *MappingAPI) StoreUnverifiedMapping(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		event       = r.URL.Query().Get("event")
	)
	defer cancel()
	defer r.Body.Close()

	if event != "up" {
		// Silently drop anything that's not a up-event
		w.WriteHeader(http.StatusOK)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warn("could not read body")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var up integration.UplinkEvent

	err = protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(b, &up)
	if err != nil {
		log.WithError(err).Warn("could not read body")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	record := types.UnverifiedMappingRecord{}
	record.ID = types.IDFromRandom()

	if up.DeviceInfo == nil {
		log.WithError(err).Warn("invalid request")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	mapperID := types.ID(sha256.Sum256([]byte(up.DeviceInfo.GetDevEui())))
	record.MapperID = mapperID.String()

	if up.TxInfo == nil {
		log.WithError(err).Warn("invalid request")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	record.Frequency = up.TxInfo.GetFrequency()
	record.Bandwidth = up.TxInfo.GetModulation().GetLora().GetBandwidth()
	record.CodeRate = up.TxInfo.GetModulation().GetLora().GetCodeRate().String()
	record.SpreadingFactor = up.TxInfo.GetModulation().GetLora().GetSpreadingFactor()

	if up.Object == nil {
		log.WithError(err).Warn("invalid request")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	latValue, ok := up.Object.Fields["latitude"]
	if !ok {
		http.Error(w, "latitude value missing", http.StatusBadRequest)
		return
	}
	record.MapperLat = latValue.GetNumberValue()

	lonValue, ok := up.Object.Fields["longitude"]
	if !ok {
		http.Error(w, "longitude value missing", http.StatusBadRequest)
		return
	}
	record.MapperLon = lonValue.GetNumberValue()
	record.MapperLocation = h3light.LatLonToCell(record.MapperLat, record.MapperLon, 10)

	accuracyValue, ok := up.Object.Fields["accuracy"]
	if !ok {
		http.Error(w, "accuracy value missing", http.StatusBadRequest)
		return
	}
	record.MapperAccuracy = accuracyValue.GetNumberValue()

	var bestGatewayID *types.ID
	var bestGatewayRssi *int32
	var bestGatewaySnr *float64
	var bestGatewayLocation *h3light.Cell

	for _, rxInfo := range up.RxInfo {
		gatewayIDStr, ok := rxInfo.Metadata["thingsix_gateway_id"]
		if !ok {
			continue
		}
		gatewayID := types.IDFromString(gatewayIDStr)

		var gatewayLocation *h3light.Cell
		gatewayLocationStr, ok := rxInfo.Metadata["thingsix_location_hex"]
		if ok {
			gatewayLocationCell, err := h3light.CellFromString(gatewayLocationStr)
			if err == nil {
				gatewayLocation = &gatewayLocationCell
			}
		}

		gatewayRecord := &types.UnverifiedMappingGatewayRecord{
			MappingID:       record.ID,
			Rssi:            rxInfo.GetRssi(),
			Snr:             float64(rxInfo.GetSnr()),
			GatewayID:       gatewayID,
			GatewayLocation: gatewayLocation,
			GatewayTime:     rxInfo.GetTime().AsTime(),
		}

		if gatewayLocation != nil && (bestGatewayID == nil || *bestGatewayRssi < rxInfo.GetRssi()) {
			bestGatewayID = &gatewayID
			bestGatewayLocation = gatewayLocation
			bestGatewayRssi = utils.Ptr(rxInfo.GetRssi())
			bestGatewaySnr = utils.Ptr(float64(rxInfo.GetSnr()))
		}

		record.GatewayRecords = append(record.GatewayRecords, gatewayRecord)
	}

	record.BestGatewayID = bestGatewayID
	record.BestGatewayLocation = bestGatewayLocation
	record.BestGatewayRssi = bestGatewayRssi
	record.BestGatewaySnr = bestGatewaySnr

	err = mapi.store.StoreUnverifiedMappingRecord(ctx, &record)
	if err != nil {
		log.WithError(err).Warn("could not store unverified mapping")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	assumedCoverage, err := mapi.store.GetAssumedUnverifiedCoverageByLocation(ctx, record.MapperLocation.Parent(8))
	if err != nil {
		log.WithError(err).Warn("could not store unverified mapping")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if assumedCoverage == nil || assumedCoverage != nil && time.Since(assumedCoverage.LatestUpdate) > 24*time.Hour {
		assumedCoverage = &types.AssumedUnverifiedCoverage{
			Location:     record.MapperLocation.Parent(8),
			LatestUpdate: time.Now(),
		}
		err := mapi.store.StoreAssumedUnverifiedCoverage(ctx, assumedCoverage)
		if err != nil {
			log.WithError(err).Warn("could not store unverified mapping")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func (mapi *MappingAPI) GetUnverifiedMappingById(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 15*time.Second)
		mappingID   = utils.IDFromRequest(r, "id")
	)
	defer cancel()

	mappingRecord, err := mapi.store.GetUnverifiedMappingRecord(ctx, mappingID)
	if err != nil {
		log.WithError(err).Error("unable to retrieve mapping-record from DB")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoding.ReplyJSON(w, r, http.StatusOK, mappingRecord)
}

package api

import (
	"net/http"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/types"
)

var (
	supportedFrequencyPlans []types.FrequencyPlan
)

func init() {
	supportedFrequencyPlans = make([]types.FrequencyPlan, len(frequency_plan.AllBands))
	for i, band := range frequency_plan.AllBands {
		supportedFrequencyPlans[i] = types.FrequencyPlan{
			ID:   uint8(band.ToBlockchain()),
			Plan: string(band),
		}
	}

}

func (gapi *GatewayAPI) SupportedFrequencyPlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	encoding.ReplyJSON(w, r, http.StatusOK, supportedFrequencyPlans)
}

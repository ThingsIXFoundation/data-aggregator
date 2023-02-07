package types

import (
	"time"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ethereum/go-ethereum/common"
)

type GatewayHistory struct {
	// ID is the ThingsIX compressed public key for this gateway
	ID              ID                       `json:"id"`
	ContractAddress common.Address           `json:"contract"`
	Version         uint8                    `json:"version"`
	Owner           *common.Address          `json:"owner"`
	AntennaGain     *float32                 `json:"antennaGain,omitempty"`
	FrequencyPlan   *frequency_plan.BandName `json:"frequencyPlan,omitempty"`
	Location        *h3light.Cell            `json:"location,omitempty"`
	Altitude        *uint                    `json:"altitude,omitempty"`

	Time        time.Time
	BlockNumber uint64
	Block       common.Hash
	Transaction common.Hash
}

func (gh *GatewayHistory) Gateway() *Gateway {
	return &Gateway{
		ID:              gh.ID,
		ContractAddress: gh.ContractAddress,
		Version:         gh.Version,
		Owner:           *gh.Owner,
		AntennaGain:     gh.AntennaGain,
		FrequencyPlan:   gh.FrequencyPlan,
		Location:        gh.Location,
		Altitude:        gh.Altitude,
	}
}

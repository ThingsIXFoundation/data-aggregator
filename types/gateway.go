package types

import (
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ethereum/go-ethereum/common"
)

type Gateway struct {
	// ID is the ThingsIX compressed public key for this gateway
	ID              ID                       `json:"id" gorm:"primaryKey;type:bytea"`
	ContractAddress common.Address           `json:"contract"`
	Version         uint8                    `json:"version"`
	Owner           common.Address           `json:"owner" gorm:"index;type:bytea"`
	AntennaGain     *float32                 `json:"antennaGain,omitempty"`
	FrequencyPlan   *frequency_plan.BandName `json:"frequencyPlan,omitempty"`
	Location        *h3light.Cell            `json:"location,omitempty"`
	Altitude        *uint                    `json:"altitude,omitempty"`
}

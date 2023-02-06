package types

import (
	"time"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ethereum/go-ethereum/common"
)

// GatewayEvent represents a log emitted by the gateway registry that is related
// to a gateway.
type GatewayEvent struct {
	ContractAddress  common.Address `json:"contract"`
	Block            common.Hash    `json:"blockHash"`
	Transaction      common.Hash    `json:"transaction"`
	BlockNumber      uint64         `json:"blockNumber"`
	TransactionIndex uint           `json:"transactionIndex"`
	LogIndex         uint           `json:"logIndex"`

	Type             GatewayEventType         `json:"type"`
	GatewayID        ID                       `json:"id"`
	Version          uint8                    `json:"version"`
	NewOwner         *common.Address          `json:"newOwner"`
	OldOwner         *common.Address          `json:"oldOwner"`
	NewLocation      *h3light.Cell            `json:"newLocation,omitempty"`
	OldLocation      *h3light.Cell            `json:"oldLocation,omitempty"`
	NewAltitude      *uint                    `json:"newAltitude,omitempty"`
	OldAltitude      *uint                    `json:"oldAltitude,omitempty"`
	NewFrequencyPlan *frequency_plan.BandName `json:"newFrequencyPlan,omitempty"`
	OldFrequencyPlan *frequency_plan.BandName `json:"oldFrequencyPlan,omitempty"`
	NewAntennaGain   *float32                 `json:"newAntennaGain,omitempty"`
	OldAntennaGain   *float32                 `json:"oldAntennaGain,omitempty"`
	Time             time.Time                `json:"time"`
}

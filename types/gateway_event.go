package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/ethereum/go-ethereum/common"
)

// GatewayEvent represents a log emitted by the gateway registry that is related
// to a gateway.
type GatewayEvent struct {
	ContractAddress  common.Address
	BlockNumber      uint64
	TransactionIndex uint
	LogIndex         uint

	Type      GatewayEventType
	GatewayID ID

	NewOwner *common.Address
	OldOwner *common.Address

	NewLocation *h3light.Cell
	OldLocation *h3light.Cell

	NewAltitude *uint
	OldAltitude *uint

	NewFrequencyPlan *frequency_plan.BandName
	OldFrequencyPlan *frequency_plan.BandName

	NewAntennaGain *float32
	OldAntennaGain *float32

	Block       common.Hash
	Transaction common.Hash // in case a transaction causes multiple events this isn't guaranteed unique and not suitable as a primary key

	Time time.Time
}

// MarshalJSON returns the given e into its JSON representation.
func (e GatewayEvent) MarshalJSON() ([]byte, error) {
	reply := map[string]interface{}{
		"gatewayId":   e.GatewayID,
		"type":        e.Type,
		"blockNumber": e.BlockNumber,
		"blockHash":   e.Block,
		"transaction": e.Transaction,
		"logIndex":    e.LogIndex,
	}

	switch e.Type {
	case GatewayOnboardedEvent:
		reply["owner"] = e.NewOwner
	case GatewayTransferredEvent:
		reply["oldOwner"] = e.OldOwner
		reply["newOwner"] = e.NewOwner
	case GatewayOffboardedEvent, GatewayUpdatedEvent:
		// do nothing
	default:
		return nil, fmt.Errorf("unsupported gateway event")
	}

	return json.Marshal(reply)
}

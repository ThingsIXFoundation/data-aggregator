package types

import (
	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ethereum/go-ethereum/common"
)

type Mapper struct {
	ID              ID                      `json:"id"`
	ContractAddress common.Address          `json:"contract"`
	Revision        uint16                  `json:"revision"`
	Owner           *common.Address         `json:"owner"`
	FrequencyPlan   frequency_plan.BandName `json:"frequencyPlan"`
	Active          bool                    `json:"active"`
}

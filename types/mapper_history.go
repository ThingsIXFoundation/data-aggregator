package types

import (
	"time"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ethereum/go-ethereum/common"
)

type MapperHistory struct {
	// ID is the ThingsIX compressed public key for this mapper
	ID              ID                      `json:"id"`
	ContractAddress common.Address          `json:"contract"`
	Revision        uint16                  `json:"revision"`
	Owner           *common.Address         `json:"owner"`
	FrequencyPlan   frequency_plan.BandName `json:"frequencyPlan"`
	Active          bool                    `json:"active"`

	Time        time.Time
	BlockNumber uint64
	Block       common.Hash
	Transaction common.Hash
}

func (mh *MapperHistory) Mapper() *Mapper {
	return &Mapper{
		ID:              mh.ID,
		ContractAddress: mh.ContractAddress,
		Revision:        mh.Revision,
		Owner:           mh.Owner,
		FrequencyPlan:   mh.FrequencyPlan,
		Active:          mh.Active,
	}
}

package models

import (
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/types"
)

type DBGatewayRewardHistory struct {

	// ID of the gateway
	GatewayID string

	// Date these rewards where issued
	Date time.Time

	// The total amount of Coverage Share Units this gateway has a the date.
	AssumedCoverageShareUnits string

	// The reward in THIX "gweis" for this gateway
	Rewards string
}

func (m *DBGatewayRewardHistory) Entity() string {
	return "GatewayRewardHistory"
}

func (m *DBGatewayRewardHistory) Key() string {
	return fmt.Sprintf("%s.%s", m.GatewayID, m.Date.String())
}

func NewDBGatewayRewardHistory(e *types.GatewayRewardHistory) *DBGatewayRewardHistory {
	return &DBGatewayRewardHistory{
		GatewayID:                 e.GatewayID.String(),
		Date:                      e.Date,
		AssumedCoverageShareUnits: e.AssumedCoverageShareUnits.String(),
		Rewards:                   e.Rewards.String(),
	}
}

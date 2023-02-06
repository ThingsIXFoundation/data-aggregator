package api

import "github.com/ThingsIXFoundation/data-aggregator/types"

type GatewayHexInfo struct {
	Count    int             `json:"count"`
	Gateways []types.Gateway `json:"gateways,omitempty"`
}

type GatewayHex struct {
	Hexes map[string]GatewayHexInfo `json:"hexes,omitempty"`
}

type Res0GatewayHex struct {
	Hexes map[string]GatewayHex `json:"hexes,omitempty"`
}

type PendingGatewayEventsResponse struct {
	Confirmations uint64                `json:"confirmations"`
	SyncedTo      uint64                `json:"syncedTo"`
	Events        []*types.GatewayEvent `json:"events"`
}

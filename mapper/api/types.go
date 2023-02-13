package api

import "github.com/ThingsIXFoundation/types"

type PendingMapperEventsResponse struct {
	Confirmations uint64               `json:"confirmations"`
	SyncedTo      uint64               `json:"syncedTo"`
	Events        []*types.MapperEvent `json:"events"`
}

package types

import "github.com/ethereum/go-ethereum/common"

type Router struct {
	ID              ID             `json:"id"`
	ContractAddress common.Address `json:"contract"`
	Owner           common.Address `json:"owner"`
	NetID           uint32         `json:"netid"`
	Prefix          uint32         `json:"prefix"`
	Mask            uint8          `json:"mask"`
	Endpoint        string         `json:"endpoint"`
}

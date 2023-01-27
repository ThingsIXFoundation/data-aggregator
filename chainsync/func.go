package chainsync

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type SetCurrentBlockFunc func(context.Context, common.Address, uint64) error
type CurrentBlockFunc func(context.Context, common.Address) (uint64, error)

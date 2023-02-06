package chainsync

import (
	"context"
)

type SetCurrentBlockFunc func(context.Context, uint64) error
type CurrentBlockFunc func(context.Context) (uint64, error)

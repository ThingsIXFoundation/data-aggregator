package interfac

import (
	"context"

	"github.com/ThingsIXFoundation/types"
)

type MappingFunc func(context.Context, *types.MappingRecord) error

type Source interface {
	Run(context.Context) error
	SetFuncs(MappingFunc)
}

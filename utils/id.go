package utils

import (
	"github.com/ThingsIXFoundation/types"
)

func StringPtrToIDtr(hex *string) *types.ID {
	if hex == nil {
		return nil
	}

	return Ptr(types.IDFromString(*hex))
}

func IDPtrToStringPtr(id *types.ID) *string {
	if id == nil {
		return nil
	}

	return Ptr(id.String())
}

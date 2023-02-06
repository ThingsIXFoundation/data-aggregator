package utils

import (
	"encoding/hex"
	"net/http"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/go-chi/chi/v5"
)

func IDFromRequest(r *http.Request, key string) types.ID {
	val := chi.URLParam(r, key)
	if val[0] == '0' && (val[1] == 'x' || val[1] == 'X') {
		val = val[2:]
	}
	raw, _ := hex.DecodeString(val)
	var id types.ID
	copy(id[:], raw)
	return id
}

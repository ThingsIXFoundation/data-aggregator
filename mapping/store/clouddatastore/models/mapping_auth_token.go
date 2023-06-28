package models

import (
	"fmt"
	"time"
)

type DBMappingAuthToken struct {
	Owner      string
	Expiration time.Time
	Code       string
	Challenge  string
}

func (e *DBMappingAuthToken) Entity() string {
	return "MappingAuthToken"
}

func (e *DBMappingAuthToken) Key() string {
	return fmt.Sprintf("%s.%s", e.Owner, e.Challenge)
}

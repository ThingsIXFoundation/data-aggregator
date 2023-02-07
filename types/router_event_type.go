package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type RouterEventType string

const (
	// RouterRegisteredEvent is raised when a router is registered by its manager
	RouterRegisteredEvent RouterEventType = "register"
	// RegisterEvent is raised when a router details are updated
	RouterUpdatedEvent RouterEventType = "update"
	// RouterRemovedEvent is raised when router details are removed
	RouterRemovedEvent RouterEventType = "removed"

	RouterUnknownEvent RouterEventType = "unknown"
)

func (event *RouterEventType) Scan(value interface{}) error {
	*event = RouterEventType(value.([]byte))
	return nil
}

func (event RouterEventType) Value() (driver.Value, error) {
	return string(event), nil
}

func (event RouterEventType) MarshalText() ([]byte, error) {
	return []byte(event), nil
}

func (event *RouterEventType) UnmarshalText(input []byte) error {
	normalized := strings.ToLower(string(input))
	switch RouterEventType(normalized) {
	case RouterRegisteredEvent:
		*event = RouterEventType(normalized)
		return nil
	default:
		return fmt.Errorf(`invalid event type "%s"`, input)
	}
}

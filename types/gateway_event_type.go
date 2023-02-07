package types

import (
	"fmt"
	"reflect"
	"strings"
)

type GatewayEventType string

const (
	GatewayOnboardedEvent   GatewayEventType = "onboard"
	GatewayOffboardedEvent  GatewayEventType = "offboard"
	GatewayUpdatedEvent     GatewayEventType = "update"
	GatewayTransferredEvent GatewayEventType = "transfer"
	GatewayUnknownEvent     GatewayEventType = "unknown"
)

func (event *GatewayEventType) Scan(value interface{}) error {
	if str, ok := value.(string); ok {
		*event = GatewayEventType(str)
		return nil
	}

	if b, ok := value.([]byte); ok {
		*event = GatewayEventType(b)
		return nil
	}
	return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
}

func (event GatewayEventType) MarshalText() ([]byte, error) {
	return []byte(event), nil
}

func (event *GatewayEventType) UnmarshalText(input []byte) error {
	normalized := strings.ToLower(string(input))
	switch GatewayEventType(normalized) {
	case GatewayOnboardedEvent, GatewayOffboardedEvent, GatewayUpdatedEvent, GatewayTransferredEvent:
		*event = GatewayEventType(normalized)
		return nil
	default:
		return fmt.Errorf(`invalid event type "%s"`, input)
	}
}

package dynamodb

import (
	"testing"

	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ethereum/go-ethereum/common"

	dadynamo "github.com/ThingsIXFoundation/data-aggregator/dynamodb"
)

func TestNewDBGatewayEvent(t *testing.T) {
	gw := types.GatewayEvent{
		ContractAddress:  common.HexToAddress("0x43E17981fEFAE6d92926f5992D8993490C541CC5"),
		BlockNumber:      1,
		TransactionIndex: 1,
		LogIndex:         1,
		Type:             types.GatewayOnboardedEvent,
		GatewayID:        types.ID{},
		NewOwner:         utils.Ptr(common.HexToAddress("0x43E17981fEFAE6d92926f5992D8993490C541CC5")),
	}

	dbgw := NewDBGatewayEvent(&gw)

	av, err := dadynamo.Marshal(dbgw)
	if err != nil {
		t.Error(err)
	}

	t.Fatalf("%s", av)
}

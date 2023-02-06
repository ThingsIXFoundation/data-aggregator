package chainsync

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	gateway_registry "github.com/ThingsIXFoundation/gateway-registry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	syncedBlockGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "synced_block",
		Help: "Synced to blockchain block",
	})
)

func init() {
	prometheus.MustRegister(syncedBlockGauge)
}

func (cs *ChainSync) runConfirmedSync(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"registry":             cs.contractAddress,
		"poll-interval":        viper.GetDuration(config.CONFIG_GATEWAY_CHAINSYNC_POLL_INTERVAL),
		"max-block-scan-range": viper.GetUint(config.CONFIG_GATEWAY_CHAINSYNC_MAX_BLOCK_SCAN_RANGE),
		"confirmations":        viper.GetUint(config.CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS),
	}).Info("integrate gateways from smart contract")

	pollInterval := time.Duration(time.Second) // first run almost instant

	// periodically check if there is gateway data that needs to be integrated
	for {
		select {
		case <-time.After(pollInterval):
			for {
				synced, err := cs.syncConfirmed(ctx)
				if err != nil {
					logrus.WithError(err).Warn("unable to integrate gateway events")
					break
				}
				if synced {
					pollInterval = viper.GetDuration(config.CONFIG_GATEWAY_CHAINSYNC_POLL_INTERVAL)
					break
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cs *ChainSync) syncConfirmed(ctx context.Context) (bool, error) {
	// dial RPC node
	client, err := chainsync.DialRpc(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to dial RPC node: %w", err)
	}
	defer client.Close()

	// retrieve where to sync from
	syncFrom, err := chainsync.GetSyncFromBlock(ctx, client, cs.contractAddress, cs.currentBlockFunc)
	if err != nil {
		return false, fmt.Errorf("unable to determine sync from block: %w", err)
	}

	// determine to sync to
	syncTo, capped, err := chainsync.GetSyncToBlock(ctx, client, syncFrom.Uint64(), viper.GetUint64(config.CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS), viper.GetUint64(config.CONFIG_GATEWAY_CHAINSYNC_MAX_BLOCK_SCAN_RANGE))
	if err != nil {
		return false, fmt.Errorf("unable to determine sync to block: %w", err)
	}

	if syncTo == nil {
		// already synced to latest confirmed block
		return true, nil
	}

	logrus.WithFields(logrus.Fields{
		"from":     syncFrom,
		"to":       syncTo,
		"contract": cs.contractAddress,
		"synced":   !capped,
	}).Info("ingesting gateway events from blockchain")

	// retrieve logs
	events, err := cs.getEvents(ctx, client, syncFrom, syncTo)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve gateway registry logs: %w", err)
	}

	err = cs.eventsFunc(ctx, events)
	if err != nil {
		return false, err
	}

	err = cs.setCurrentBlockFunc(ctx, syncTo.Uint64())
	if err != nil {
		return false, err
	}

	return !capped, nil
}

func (cs *ChainSync) getEvents(ctx context.Context, client *ethclient.Client, from, to *big.Int) ([]*types.GatewayEvent, error) {
	logrus.WithFields(logrus.Fields{
		"fromBlock": from,
		"to":        to,
		"address":   cs.contractAddress,
	}).Trace("get events")
	logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: from,
		ToBlock:   to,
		Addresses: []common.Address{cs.contractAddress},
	})
	if err != nil {
		logrus.WithError(err).Error("error while getting gateway events")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"fromBlock": from,
		"to":        to,
		"address":   cs.contractAddress,
		"#":         len(logs),
	}).Debug("retrieve gateway registry events")

	gatewayRegistry, err := gateway_registry.NewGatewayRegistryCaller(cs.contractAddress, client)
	if err != nil {
		logrus.WithError(err).Error("error while creating gateway-registry caller")
		return nil, err
	}

	// decode logs into gateway events and filter out non gateway events
	var (
		events []*types.GatewayEvent
	)

	for _, l := range logs {
		logrus.WithFields(logrus.Fields{
			"block": l.BlockHash,
			"tx":    l.TxHash,
			"type":  l.Topics[0],
		}).Trace("event")
		if event := decodeLogToGatewayEvent(&l); event != nil {
			if event.Type == types.GatewayOnboardedEvent {
				gateway, err := gatewayDetails(gatewayRegistry, cs.contractAddress, l.BlockNumber, event.GatewayID)
				if err != nil {
					logrus.WithError(err).Error("error while getting added gateway details")
					return nil, err
				}
				event.Version = gateway.Version
			}
			// The GatewayUpdated event doesn't contain the gateway details. So fetch the gateway details before and after the update
			// and include them in the event
			if event.Type == types.GatewayUpdatedEvent {
				gatewayBefore, err := gatewayDetails(gatewayRegistry, cs.contractAddress, l.BlockNumber-1, event.GatewayID)
				if err != nil {
					logrus.WithError(err).Error("error while getting before-update gateway details")
					return nil, err
				}
				gatewayAfter, err := gatewayDetails(gatewayRegistry, cs.contractAddress, l.BlockNumber, event.GatewayID)
				if err != nil {
					logrus.WithError(err).Error("error while getting updated gateway details")
					return nil, err
				}

				event.OldOwner = utils.Ptr(gatewayBefore.Owner)
				event.OldFrequencyPlan = gatewayBefore.FrequencyPlan
				event.OldAltitude = gatewayBefore.Altitude
				event.OldLocation = gatewayBefore.Location
				event.OldAntennaGain = gatewayBefore.AntennaGain

				event.NewOwner = utils.Ptr(gatewayAfter.Owner)
				event.NewFrequencyPlan = gatewayAfter.FrequencyPlan
				event.NewAltitude = gatewayAfter.Altitude
				event.NewLocation = gatewayAfter.Location
				event.NewAntennaGain = gatewayAfter.AntennaGain
			}

			eventTime, err := chainsync.BlockTime(ctx, client, event.BlockNumber)
			if err != nil {
				logrus.WithError(err).Error("error while getting time of block")
				return nil, err
			}
			event.ContractAddress = cs.contractAddress
			event.Time = eventTime
			events = append(events, event)
		}
	}

	return events, nil
}

package chainsync

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	router_registry "github.com/ThingsIXFoundation/router-registry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (cs *ChainSync) runConfirmedSync(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"registry":             cs.contractAddress,
		"poll-interval":        viper.GetDuration(config.CONFIG_ROUTER_CHAINSYNC_POLL_INTERVAL),
		"max-block-scan-range": viper.GetUint(config.CONFIG_ROUTER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE),
		"confirmations":        viper.GetUint(config.CONFIG_ROUTER_CHAINSYNC_CONFORMATIONS),
	}).Info("integrate routers from smart contract")

	pollInterval := time.Duration(time.Second) // first run almost instant

	// periodically check if there is router data that needs to be integrated
	for {
		select {
		case <-time.After(pollInterval):
			for {
				synced, err := cs.syncConfirmed(ctx)
				if err != nil {
					logrus.WithError(err).Warn("unable to integrate router events")
					break
				}
				if synced {
					pollInterval = viper.GetDuration(config.CONFIG_ROUTER_CHAINSYNC_POLL_INTERVAL)
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
	syncTo, capped, err := chainsync.GetSyncToBlock(ctx, client, syncFrom.Uint64(), viper.GetUint64(config.CONFIG_ROUTER_CHAINSYNC_CONFORMATIONS), viper.GetUint64(config.CONFIG_ROUTER_CHAINSYNC_MAX_BLOCK_SCAN_RANGE))
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
	}).Info("ingesting router events from blockchain")

	// retrieve logs
	events, err := cs.getEvents(ctx, client, syncFrom, syncTo)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve router registry logs: %w", err)
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

func (cs *ChainSync) getEvents(ctx context.Context, client *ethclient.Client, from, to *big.Int) ([]*types.RouterEvent, error) {
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
		logrus.WithError(err).Error("error while getting router events")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"fromBlock": from,
		"to":        to,
		"address":   cs.contractAddress,
		"#":         len(logs),
	}).Debug("retrieve router registry events")

	routerRegistry, err := router_registry.NewRouterRegistryCaller(cs.contractAddress, client)
	if err != nil {
		logrus.WithError(err).Error("error while creating router-registry caller")
		return nil, err
	}

	// decode logs into router events and filter out non router events
	var (
		events []*types.RouterEvent
	)

	for _, log := range logs {
		logrus.WithFields(logrus.Fields{
			"block": log.BlockHash,
			"tx":    log.TxHash,
			"type":  log.Topics[0],
		}).Trace("event")
		event, err := decodeLogToRouterEvent(ctx, &log, client, routerRegistry, cs.contractAddress)
		if event == nil {
			continue
		}
		if err != nil {
			logrus.WithError(err).Error("error while processing router logs")
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

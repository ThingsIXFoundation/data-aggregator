package chainsync

import (
	"context"
	"fmt"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/chainsync"
	"github.com/ThingsIXFoundation/data-aggregator/config"
	"github.com/ThingsIXFoundation/data-aggregator/types"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	gateway_registry "github.com/ThingsIXFoundation/gateway-registry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	GatewayOnboardedEvent   = common.BytesToHash(crypto.Keccak256([]byte("GatewayOnboarded(bytes32,address)")))
	GatewayOffboardedEvent  = common.BytesToHash(crypto.Keccak256([]byte("GatewayOffboarded(bytes32)")))
	GatewayUpdatedEvent     = common.BytesToHash(crypto.Keccak256([]byte("GatewayUpdated(bytes32)")))
	GatewayTransferredEvent = common.BytesToHash(crypto.Keccak256([]byte("GatewayTransferred(bytes32,address,address)")))
)

func (cs *ChainSync) runPending(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"registry":      cs.contractAddress,
		"confirmations": viper.GetUint(config.CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS),
	}).Info("syncing pending gateway events from smart contract")

	if viper.GetUint(config.CONFIG_GATEWAY_CHAINSYNC_CONFORMATIONS) == 0 {
		logrus.Info("confirmations 0, don't integrate pending events")
		<-ctx.Done() // wait until the shutdown signal is given
		return nil
	}

	// periodically check if there is gateway data that needs to be integrated
	var (
		retry    = 5 * time.Second
		lastTime time.Time
	)

	for {
		select {
		case <-time.After(retry):
			lastTime = time.Now()
			if err := cs.handlePending(ctx); err != nil {
				logrus.WithError(err).Warn("integrate pending gateway events stopped")
			}
			if lastTime.Before(time.Now().Add(-time.Minute)) {
				retry = time.Millisecond
			} else {
				retry *= 2
				if retry > time.Minute {
					retry = time.Minute
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cs *ChainSync) handlePending(ctx context.Context) error {
	// dial RPC node
	client, err := chainsync.DialRpc(ctx)
	if err != nil {
		return fmt.Errorf("unable to dial RPC node: %w", err)
	}
	defer client.Close()

	// wait for new gateway related events
	var (
		q = ethereum.FilterQuery{
			Addresses: []common.Address{cs.contractAddress},
			Topics: [][]common.Hash{
				{
					GatewayOnboardedEvent,
					GatewayOffboardedEvent,
					GatewayUpdatedEvent,
					GatewayTransferredEvent,
				},
			},
		}
		logs = make(chan etypes.Log)
	)

	// retrieve new onboard logs and integrate them into the pending gateway events table
	sub, err := client.SubscribeFilterLogs(ctx, q, logs)
	if err != nil {
		return fmt.Errorf("unable to subscribe to gateway registry events: %w", err)
	}

	gatewayRegistry, err := gateway_registry.NewGatewayRegistryCaller(cs.contractAddress, client)
	if err != nil {
		logrus.WithError(err).Error("error while creating gateway-registry caller")
		return err
	}

	// begin integrating events
	for {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
			return nil
		case err, ok := <-sub.Err():
			if ok {
				return fmt.Errorf("waiting for pendings log subscription failed: %w", err)
			}
			return nil
		case l, ok := <-logs:
			if !ok {
				return fmt.Errorf("unable to retrieve pending gateway logs")
			}
			if event := decodeLogToGatewayEvent(&l); event != nil {
				// The GatewayUpdated event doesn't contain the gateway details. So fetch the gateway details before and after the update
				// and include them in the event
				if event.Type == types.GatewayUpdatedEvent {
					gatewayBefore, err := gatewayDetails(gatewayRegistry, cs.contractAddress, l.BlockNumber-1, event.GatewayID)
					if err != nil {
						logrus.WithError(err).Error("error while getting before-update gateway details")
						return err
					}
					gatewayAfter, err := gatewayDetails(gatewayRegistry, cs.contractAddress, l.BlockNumber, event.GatewayID)
					if err != nil {
						logrus.WithError(err).Error("error while getting updated gateway details")
						return err
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
					return err
				}
				event.ContractAddress = cs.contractAddress
				event.Time = eventTime

				cs.pendingEventFunc(ctx, event)

				return nil
			}
		}
	}
}

package relay

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	extTypes "gitlab.inlive7.com/crypto/ethereum-relay/pkg/types"
)

func (r *Relay) SubscribeTokenEvents(symbol string, address string, realtimeLogs chan extTypes.EventLogTransfer) error {
	tokenAddress := r.supportTokens[symbol]
	if tokenAddress == "" {
		return errors.New("token not match any of supported tokens")
	}
	var targetAddress []common.Hash
	if address != "" {
		targetAddress = []common.Hash{common.HexToHash(address)}
	}
	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(tokenAddress)},
		Topics:    [][]common.Hash{nil, nil, targetAddress}, // T0 is method, T1 is from, we only filter T2 which is "To"
	}
	sub, err := r.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}
	r.subscriptions[realtimeLogs] = &subscription{
		Output: realtimeLogs,
		Sub:    sub,
		Log:    logs,
		Enable: true,
	}
	go r.startSubscription(r.subscriptions[realtimeLogs])
	return nil
}

// func (r *Relay) UnsubscribeEvent(key chan extTypes.EventLogTransfer) {
// 	sub := r.subscriptions[key]
// 	if sub != nil {
// 		sub.Enable = false
// 		r.subscriptions[key] = nil
// 		close(key)
// 	}
// }

// func (r *Relay) UnsubscribeAllEvents() {
// 	for key, sub := range r.subscriptions {
// 		sub.Enable = false
// 		r.subscriptions[key] = nil
// 		close(key)
// 	}
// }

func (r *Relay) startSubscription(sub *subscription) {
	defer func() {
		recover()
		sub.Enable = false
		sub.Sub.Unsubscribe()
		r.subscriptions[sub.Output] = nil
	}()
	for sub.Enable {
		select {
		case e := <-sub.Sub.Err():
			if e != nil {
				fmt.Println(e)
				close(sub.Output)
			}
		case l := <-sub.Log:
			logTransferSig := []byte("Transfer(address,address,uint256)")
			logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
			// Topics[0] is always the method signature
			switch l.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				sub.Output <- extTypes.EventLogTransfer{
					Txn:    l.TxHash.String(),
					From:   common.HexToAddress(l.Topics[1].Hex()).String(),
					To:     common.HexToAddress(l.Topics[2].Hex()).String(),
					Tokens: new(big.Int).SetBytes(l.Data),
				}
			default:
				fmt.Println("event log no match")
			}
		default:
		}
	}
}

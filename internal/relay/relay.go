package relay

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"

	"gitlab.inlive7.com/crypto/ethereum-relay/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var instance *Relay

/*
	Use Shared to get same Relay struct to do the logics,
	it stores ethclient.Client to prevent alloc & dealloc multiple times
*/
func Shared(chainID config.ChainID) *Relay {

	if instance != nil && instance.currentChainInfo.ID == chainID {
		return instance
	}

	info, err := config.RetrieveChainInfo(chainID)
	if err != nil {
		log.Println(err.Error())
		return instance
	}
	if instance != nil {
		instance.destory()
		instance = nil
	}
	instance, err = createInstance(info)
	if err != nil {
		log.Fatal(err.Error())
	}
	return instance
}

/*
	This method is for hot-update usecase. If we manage to update the yml files,
	then destory instance to make it load the newer version of yml file.
*/
func Destory() {
	instance.destory()
	instance = nil
}

func (r Relay) QueryTransaction(txn string) (trans *TransactionState, isPending bool, err error) {
	tx, isPending, err := r.client.TransactionByHash(context.Background(), common.HexToHash(txn))
	if err != nil {
		return trans, isPending, err
	}
	if isPending {
		return trans, isPending, err
	}

	receipt, err := r.client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return trans, isPending, err
	}
	// Txn Type: 2 (EIP-1559)
	// receipt.Type

	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), tx.GasPrice())
	if err != nil {
		return trans, isPending, err
	}

	h, err := r.client.HeaderByHash(context.Background(), receipt.BlockHash)
	if err != nil {
		return trans, isPending, err
	}

	return &TransactionState{
		receipt.Status == 1,
		tx.Value(),
		msg.From().Hex(),
		tx.To().Hex(),
		tx.GasPrice(),
		tx.Gas(),
		h.Time,
		uint16(tx.ChainId().Uint64()),
		r.currentChainInfo.Name,
		msg.Nonce(),
	}, isPending, err
}

func (r Relay) SendTransaction(privateKey string, to string, value *big.Int) error {
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := r.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	gasLimit := uint64(21000)           // in units
	gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	toAddress := common.HexToAddress(to)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(int64(r.currentChainInfo.ID))), pk)
	if err != nil {
		return err
	}
	fmt.Println(signedTx)
	return nil
}

func (r Relay) GasPrice() (*EstimateGasInfo, error) {
	price, pErr := r.client.SuggestGasPrice(context.Background())
	tip, tErr := r.client.SuggestGasTipCap(context.Background())
	if pErr != nil || tErr != nil {
		return nil, errors.New("GasPrice failed")
	}
	return &EstimateGasInfo{price, tip}, nil
}

// ******** PRIVATE ******** //

func createInstance(c config.ChainInfo) (*Relay, error) {
	p := config.GetProviderInfo(c.ProviderFile)
	client, err := ethclient.Dial(p.URL)
	if err != nil {
		return nil, errors.New("ethclient dial failed")
	}

	return &Relay{
		c,
		client,
	}, nil
}

func (r Relay) destory() {
	r.client.Close()
}

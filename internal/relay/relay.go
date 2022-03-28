package relay

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"strings"

	"gitlab.inlive7.com/crypto/ethereum-relay/config"
	token "gitlab.inlive7.com/crypto/ethereum-relay/contracts/dist"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var instances = make(map[config.ChainID]*Relay)

/*
	Use Shared to get same Relay struct to do the logics,
	it stores ethclient.Client to prevent alloc & dealloc multiple times
*/
func Shared(chainID config.ChainID) *Relay {

	if instances[chainID] != nil {
		return instances[chainID]
	}

	info, err := config.RetrieveChainInfo(chainID)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	instance, err := createInstance(info)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	instances[chainID] = instance
	return instances[chainID]
}

/*
	This method is for hot-update usecase. If we manage to update the yml files,
	then destory instance to make it load the newer version of yml file.
*/
func Destory() {
	for _, v := range instances {
		v.destory()
	}
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

func (r Relay) SendTransaction(privateKey string, data *TransactionRaw) error {
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

	gasLimit := uint64(21000)          // standard transfer limit, see https://ethereum.org/en/developers/docs/gas/, https://eips.ethereum.org/EIPS/eip-1559
	gasPrice := data.PreferredGasPrice // usually Gwei
	toAddress := common.HexToAddress(data.To)
	cID := big.NewInt(int64(r.currentChainInfo.ID))
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   cID,
		Nonce:     nonce,
		GasFeeCap: gasPrice,
		GasTipCap: gasPrice,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     data.Value,
		Data:      nil,
	})
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(int64(r.currentChainInfo.ID))), pk)
	if err != nil {
		return err
	}
	result := r.client.SendTransaction(context.Background(), signedTx)
	return result
}

func (r Relay) GasPrice() (*EstimateGasInfo, error) {
	price, pErr := r.client.SuggestGasPrice(context.Background())
	tip, tErr := r.client.SuggestGasTipCap(context.Background())
	// BSC網路取得不到 GasTip
	if pErr != nil || (tErr != nil && r.currentChainInfo.ID != config.BscMainnet && r.currentChainInfo.ID != config.BscTestnet) {
		return nil, errors.New("GasPrice failed")
	}
	return &EstimateGasInfo{price, tip}, nil
}

func (r Relay) GetBalance(address string) (balance *big.Int, err error) {
	balance, err = r.client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func (r Relay) GetBalanceForToken(address string, symbol string) (*big.Int, uint8, error) {
	tokenAddress := r.supportTokens[strings.ToLower(symbol)]
	if tokenAddress == "" {
		return nil, 0, errors.New("can not find matched token")
	}
	tAddr := common.HexToAddress(tokenAddress)
	instance, err := token.NewToken(tAddr, r.client)
	if err != nil {
		log.Fatal(err)
	}
	balance, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		log.Fatal(err)

	}
	decimal, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)

	}
	return balance, decimal, nil
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
		p.Tokens,
		client,
	}, nil
}

func (r Relay) destory() {
	r.client.Close()
}

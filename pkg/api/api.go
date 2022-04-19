package api

import (
	"math/big"

	"gitlab.inlive7.com/crypto/ethereum-relay/config"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/manager"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/relay"
	"gitlab.inlive7.com/crypto/ethereum-relay/pkg/types"
)

type EthereumRelay interface {
	VerifySignature(hash string, signatureHex string) (publicAddress string, err error)
	CreateNewAccount() (privateKey string, publicKey string, publicAddress string)
	QueryTransaction(chainID uint16, txn string) (*types.TransactionState, bool, error)
	TransferValueUsingPrivateKey(chainID uint16, privateKey string, data *types.TransactionRaw) (hash string, err error)
	TransferTokenUsingPrivateKey(chainID uint16, privateKey string, data *types.TransactionRaw) (hash string, err error)
	GetGasPrice(chainID uint16) (*types.EstimateGasInfo, error)
	GetBalance(chainID uint16, address string) (balance *big.Int, err error)
	GetBalanceForToken(chainID uint16, address string, symbol string) (balance *big.Int, decimal uint8, err error)
	GetTokenAddress(chainID uint16, symbol string) (address string)
	InitRelay(chainIds []config.ChainID)
}

func VerifySignature(hash string, signatureHex string) (publicAddress string, err error) {
	return manager.VerifySignature(hash, signatureHex)
}

func CreateNewAccount() (privateKey string, publicKey string, publicAddress string) {
	return manager.CreateNewAccount()
}

func QueryTransaction(chainID uint16, txn string) (trans *types.TransactionState, isPending bool, err error) {
	r := relay.Shared(config.ChainID(chainID))
	return r.QueryTransaction(txn)
}

func TransferValueUsingPrivateKey(chainID uint16, privateKey string, data *types.TransactionRaw) (hash string, err error) {
	return relay.Shared(config.ChainID(chainID)).TransferValue(privateKey, data)
}

func TransferTokenUsingPrivateKey(chainID uint16, privateKey string, data *types.TransactionRaw) (hash string, err error) {
	return relay.Shared(config.ChainID(chainID)).TransferToken(privateKey, data)
}

func GetGasPrice(chainID uint16) (info *types.EstimateGasInfo, err error) {
	return relay.Shared(config.ChainID(chainID)).GasPrice()
}

func GetGasLimit(chainID uint16, symbol string, from string, to string, value *big.Int) (gasLimit uint64, err error) {
	return relay.Shared(config.ChainID(chainID)).GasLimit(symbol, from, to, value)
}

// balance in wei
func GetBalance(chainID uint16, address string) (balance *big.Int, err error) {
	r := relay.Shared(config.ChainID(chainID))
	return r.GetBalance(address)
}

func GetBalanceForToken(chainID uint16, address string, symbol string) (balance *big.Int, decimal uint8, err error) {
	r := relay.Shared(config.ChainID(chainID))
	return r.GetBalanceForToken(address, symbol)
}

func GetTokenAddress(chainID uint16, symbol string) (address string) {
	return relay.Shared(config.ChainID(chainID)).GetTokenAddress(symbol)
}

func InitRelay(chainIds []config.ChainID) {
	for i := range chainIds {
		// first time call Shared inits the instance
		relay.Shared(chainIds[i])
	}
}

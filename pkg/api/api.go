package api

import (
	"math/big"

	"gitlab.inlive7.com/crypto/ethereum-relay/config"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/manager"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/relay"
)

type EthereumRelay interface {
	VerifySignature(hash string, signatureHex string) (publicAddress string, err error)
	CreateNewAccount() (privateKey string, publicKey string, publicAddress string)
	QueryTransaction(chainID uint16, txn string) (*relay.TransactionState, bool, error)
	SendTransactionUsingPrivateKey(chainID uint16, privateKey string, data *relay.TransactionRaw) error
	GetGasPrice(chainID uint16) (*relay.EstimateGasInfo, error)
	GetBalance(chainID uint16, address string) (balance *big.Int, err error)
}

func VerifySignature(hash string, signatureHex string) (publicAddress string, err error) {
	return manager.VerifySignature(hash, signatureHex)
}

func CreateNewAccount() (privateKey string, publicKey string, publicAddress string) {
	return manager.CreateNewAccount()
}

func QueryTransaction(chainID uint16, txn string) (trans *relay.TransactionState, isPending bool, err error) {
	r := relay.Shared(config.ChainID(chainID))
	return r.QueryTransaction(txn)
}

func SendTransactionUsingPrivateKey(chainID uint16, privateKey string, data *relay.TransactionRaw) error {
	relay.Shared(config.ChainID(chainID)).SendTransaction(privateKey, data)
	return nil
}

func GetGasPrice(chainID uint16) (info *relay.EstimateGasInfo, err error) {
	return relay.Shared(config.ChainID(chainID)).GasPrice()
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

func InitRelay(chainIds []config.ChainID) {
	for i := range chainIds {
		// first time call Shared inits the instance
		relay.Shared(chainIds[i])
	}
}

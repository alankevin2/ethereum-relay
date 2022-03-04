package api

import (
	"errors"
	"ethereum-relay/config"
	"ethereum-relay/internal/manager"
	"ethereum-relay/internal/relay"
	"ethereum-relay/internal/utility"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(originMsg string, signatureHex string) (publicAddress string, err error) {
	msg := utility.SignHash(originMsg)
	signature := hexutil.MustDecode(signatureHex)

	if signature[crypto.RecoveryIDOffset] != 27 && signature[crypto.RecoveryIDOffset] != 28 {
		errMsg := "crypto.RecoveryIDOffset not 27 nor 28"
		log.Println(errMsg)
		return "", errors.New(errMsg)
	}
	signature[crypto.RecoveryIDOffset] -= 27

	pubKey, err := crypto.SigToPub(msg, signature)
	if err != nil {
		errMsg := "crypto.SigToPub failed"
		log.Println(errMsg)
		return "", errors.New(errMsg)
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	publicAddressString := recoveredAddress.Hex()

	return publicAddressString, nil
}

func CreateNewAccount() (privateKey string, publicKey string, publicAddress string) {
	return manager.CreateNewAccount()
}

func QueryTransaction(chainID uint16, txn string) (*relay.TransactionState, bool, error) {
	r := relay.Shared(config.ChainID(chainID))
	return r.QueryTransaction(txn)
}

func SendTransactionUsingPrivateKey(chainID uint16, privateKey string, to string, value *big.Int) error {
	relay.Shared(config.ChainID(chainID)).SendTransaction(privateKey, to, value)
	return nil
}

func GetGasPrice(chainID uint16) (*relay.EstimateGasInfo, error) {
	return relay.Shared(config.ChainID(chainID)).GasPrice()
}

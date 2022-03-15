package api

import (
	"encoding/json"
	"errors"
	"log"
	"math/big"

	"gitlab.inlive7.com/crypto/ethereum-relay/config"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/manager"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/relay"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/utility"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type EthereumRelay interface {
	VerifySignature(hash string, signatureHex string) (publicAddress string, err error)
	CreateNewAccount() (privateKey string, publicKey string, publicAddress string)
	QueryTransaction(chainID uint16, txn string) (*relay.TransactionState, bool, error)
	SendTransactionUsingPrivateKey(chainID uint16, privateKey string, data *relay.TransactionRaw) error
	GetGasPrice(chainID uint16) (*relay.EstimateGasInfo, error)
	GetBalance(chainID uint16, address string) (balance *big.Int, err error)
}

// This verification function is for sign_personal
// func VerifySignature(originMsg string, signatureHex string) (publicAddress string, err error) {
// 	msg := utility.SignHash(originMsg)
// 	signature := hexutil.MustDecode(signatureHex)

// 	if signature[crypto.RecoveryIDOffset] != 27 && signature[crypto.RecoveryIDOffset] != 28 {
// 		errMsg := "crypto.RecoveryIDOffset not 27 nor 28"
// 		log.Println(errMsg)
// 		return "", errors.New(errMsg)
// 	}
// 	signature[crypto.RecoveryIDOffset] -= 27

// 	pubKey, err := crypto.SigToPub(msg, signature)
// 	if err != nil {
// 		errMsg := "crypto.SigToPub failed"
// 		log.Println(errMsg)
// 		return "", errors.New(errMsg)
// 	}

// 	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
// 	publicAddressString := recoveredAddress.Hex()

// 	return publicAddressString, nil
// }

// This verification function is for sign_typedData, not for sign_personal
func VerifySignature(hash string, signatureHex string) (publicAddress string, err error) {
	var typedData apitypes.TypedData
	json.Unmarshal([]byte(hash), &typedData)
	msg, err := utility.EIP712Hash(typedData)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	signature := hexutil.MustDecode(signatureHex)

	if signature[crypto.RecoveryIDOffset] != 27 && signature[crypto.RecoveryIDOffset] != 28 {
		errMsg := "crypto.RecoveryIDOffset not 27 nor 28"
		log.Println(errMsg)
		return "", errors.New(errMsg)
	}
	signature[crypto.RecoveryIDOffset] -= 27

	recoveredAddress, _ := crypto.Ecrecover(msg, signature)

	pubKey, err := crypto.UnmarshalPubkey(recoveredAddress)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	publicAddressString := crypto.PubkeyToAddress(*pubKey)

	return publicAddressString.Hex(), nil
}

func CreateNewAccount() (privateKey string, publicKey string, publicAddress string) {
	return manager.CreateNewAccount()
}

func QueryTransaction(chainID uint16, txn string) (*relay.TransactionState, bool, error) {
	r := relay.Shared(config.ChainID(chainID))
	return r.QueryTransaction(txn)
}

func SendTransactionUsingPrivateKey(chainID uint16, privateKey string, data *relay.TransactionRaw) error {
	relay.Shared(config.ChainID(chainID)).SendTransaction(privateKey, data)
	return nil
}

func GetGasPrice(chainID uint16) (*relay.EstimateGasInfo, error) {
	return relay.Shared(config.ChainID(chainID)).GasPrice()
}

// balance in wei
func GetBalance(chainID uint16, address string) (balance *big.Int, err error) {
	r := relay.Shared(config.ChainID(chainID))
	return r.GetBalance(address)
}

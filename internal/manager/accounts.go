package manager

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"log"

	"gitlab.inlive7.com/crypto/ethereum-relay/internal/utility"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func CreateNewAccount() (privateKey string, publicKey string, publicAddress string) {
	defer func() {
		if err := recover(); err != nil {
			privateKey = ""
			publicKey = ""
		}
	}()
	pk, err := crypto.GenerateKey()
	if err != nil {
		log.Panic("crypto.GenerateKey failed")
	}
	pkECDSA := crypto.FromECDSA(pk)
	// strip the '0x'
	privateKey = hexutil.Encode(pkECDSA)[2:]

	pb := pk.Public()
	publicKeyECDSA, ok := pb.(*ecdsa.PublicKey)
	if !ok {
		log.Panic("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// strip the '0x' and '04' which is for EC prefix
	publicKey = hexutil.Encode(publicKeyBytes)[4:]

	publicAddress = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return privateKey, publicKey, publicAddress
}

// This verification function is for sign_typedData
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

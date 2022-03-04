package manager

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

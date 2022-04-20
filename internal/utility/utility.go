package utility

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"golang.org/x/crypto/sha3"
)

func WeiToEther(val *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(val), big.NewFloat(params.Ether))
}

func WeiToGwei(val *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(val), big.NewFloat(params.GWei))
}

func Gwei(wei int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(wei), big.NewInt(params.GWei))
}

func StringWithoutExponent(val *big.Float) string {
	f, _ := val.Float64()
	return strconv.FormatFloat(f, 'f', -1, 32)
}

// use for sign_personal
func SignHash(str string) []byte {
	data := []byte(str)
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// use for sign_typedData
func EIP712Hash(typedData apitypes.TypedData) (hash []byte, err error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash = crypto.Keccak256(rawData)
	return
}

func GetAddressFromPrivateKey(pk *ecdsa.PrivateKey) (*common.Address, error) {
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &fromAddress, nil
}

func GetInputDataForTokenTransfer(to common.Address, value *big.Int) []byte {
	transferFnSignature := []byte("transfer(address,uint256)") // do not include spaces in the string
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4] // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(value.Bytes(), 32)

	var inputData []byte
	inputData = append(inputData, methodID...)
	inputData = append(inputData, paddedAddress...)
	inputData = append(inputData, paddedAmount...)
	return inputData
}

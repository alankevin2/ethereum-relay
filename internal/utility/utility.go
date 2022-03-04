package utility

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

func WeiToEther(val *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(val), big.NewFloat(params.Ether))
}

func WeiToGwei(val *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(val), big.NewFloat(params.GWei))
}

func StringWithoutExponent(val *big.Float) string {
	f, _ := val.Float64()
	return strconv.FormatFloat(f, 'f', -1, 32)
}

func SignHash(str string) []byte {
	data := []byte(str)
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

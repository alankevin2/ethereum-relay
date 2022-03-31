package types

import (
	"math/big"
)

type TransactionState struct {
	Success   bool // success: status = 1, fail: status = 0
	Value     *big.Int
	From      string
	To        string
	GasPrice  *big.Int
	Gas       uint64
	Time      uint64 // in Second
	Chain     uint16 // current chain id number not more than 2000
	ChainName string
	UserNonce uint64
}

type EstimateGasInfo struct {
	Base *big.Int
	Tip  *big.Int
}

type TransactionRaw struct {
	Value                 *big.Int
	To                    string
	PreferredBaseGasPrice *big.Int
	PreferredTipGasPrice  *big.Int
	TokenSymbol           string
}

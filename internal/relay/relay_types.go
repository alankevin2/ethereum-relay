package relay

import (
	"ethereum-relay/config"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Relay struct {
	currentChainInfo config.ChainInfo
	client           *ethclient.Client
}

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

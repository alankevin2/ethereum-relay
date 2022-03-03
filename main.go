package main

import (
	"ethereum-relay/config"
	"ethereum-relay/internal/utility"
	"math/big"
)

func main() {
	config.InitializeConfiguration()
	utility.WeiToEther(big.NewInt(1000000))
}

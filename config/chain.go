package config

import (
	"fmt"
)

type ChainID uint16

const (
	mainnet           ChainID = 1
	modren            ChainID = 2
	ropsten           ChainID = 3
	rinkeby           ChainID = 4
	goerli            ChainID = 5
	kovan             ChainID = 42
	gethPrivateChains ChainID = 1337
)

type ChainInfo struct {
	Name         string
	ID           ChainID
	ProviderFile string
}

func RetrieveChainInfo(id ChainID) (ChainInfo, error) {
	var info ChainInfo
	switch id {
	case mainnet:
		info = ChainInfo{"mainnet", mainnet, "provider-mainnet.yml"}
	case rinkeby:
		info = ChainInfo{"rinkeby", rinkeby, "provider-testnet-rinkeby.yml"}
	default:
		return info, fmt.Errorf("no support yet for chain id : %d", id)
	}

	return info, nil
}

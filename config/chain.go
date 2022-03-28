package config

import (
	"fmt"
)

type ChainID uint16

const (
	Mainnet           ChainID = 1
	Modren            ChainID = 2
	Ropsten           ChainID = 3
	Rinkeby           ChainID = 4
	Goerli            ChainID = 5
	Kovan             ChainID = 42
	GethPrivateChains ChainID = 1337
	BscMainnet        ChainID = 56
	BscTestnet        ChainID = 97
)

type ChainInfo struct {
	Name         string
	ID           ChainID
	ProviderFile string
	Decimal      int8
}

func RetrieveChainInfo(id ChainID) (ChainInfo, error) {
	var info ChainInfo
	switch id {
	case Mainnet:
		info = ChainInfo{"mainnet", Mainnet, "provider-mainnet.yml", 18}
	case Rinkeby:
		info = ChainInfo{"rinkeby", Rinkeby, "provider-testnet-rinkeby.yml", 18}
	case BscMainnet:
		info = ChainInfo{"bscMainnet", BscMainnet, "provider-bsc-mainnet.yml", 18}
	case BscTestnet:
		info = ChainInfo{"bscTestnet", BscTestnet, "provider-bsc-testnet.yml", 18}
	default:
		return info, fmt.Errorf("no support yet for chain id : %d", id)
	}

	return info, nil
}

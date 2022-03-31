package relay

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"strings"

	"gitlab.inlive7.com/crypto/ethereum-relay/config"
	token "gitlab.inlive7.com/crypto/ethereum-relay/contracts/dist"
	extTypes "gitlab.inlive7.com/crypto/ethereum-relay/pkg/types"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Relay struct {
	currentChainInfo config.ChainInfo
	supportTokens    map[string]string
	client           *ethclient.Client
}

var instances = make(map[config.ChainID]*Relay)

/*
	Use Shared to get same Relay struct to do the logics,
	it stores ethclient.Client to prevent alloc & dealloc multiple times
*/
func Shared(chainID config.ChainID) *Relay {

	if instances[chainID] != nil {
		return instances[chainID]
	}

	info, err := config.RetrieveChainInfo(chainID)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	instance, err := createInstance(info)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	instances[chainID] = instance
	return instances[chainID]
}

/*
	This method is for hot-update usecase. If we manage to update the yml files,
	then destory instance to make it load the newer version of yml file.
*/
func Destory() {
	for _, v := range instances {
		v.destory()
	}
}

func (r Relay) QueryTransaction(txn string) (trans *extTypes.TransactionState, isPending bool, err error) {
	tx, isPending, err := r.client.TransactionByHash(context.Background(), common.HexToHash(txn))
	if err != nil {
		return trans, isPending, err
	}
	if isPending {
		return trans, isPending, err
	}

	receipt, err := r.client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return trans, isPending, err
	}
	// Txn Type: 2 (EIP-1559)
	// receipt.Type

	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), tx.GasPrice())
	if err != nil {
		return trans, isPending, err
	}

	h, err := r.client.HeaderByHash(context.Background(), receipt.BlockHash)
	if err != nil {
		return trans, isPending, err
	}

	return &extTypes.TransactionState{
		Success:   receipt.Status == 1,
		Value:     tx.Value(),
		From:      msg.From().Hex(),
		To:        tx.To().Hex(),
		GasPrice:  tx.GasPrice(),
		Gas:       tx.Gas(),
		Time:      h.Time,
		Chain:     uint16(tx.ChainId().Uint64()),
		ChainName: r.currentChainInfo.Name,
		UserNonce: msg.Nonce(),
	}, isPending, err
}

func (r Relay) TransferValue(privateKey string, data *extTypes.TransactionRaw) (string, error) {
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	fromAddress, err := r.getAddressFromPrivateKey(pk)
	if err != nil {
		return "", err
	}

	nonce, err := r.getNonceFromAddress(*fromAddress)
	if err != nil {
		return "", err
	}

	gasLimit := uint64(21000) // standard transfer limit, see https://ethereum.org/en/developers/docs/gas/, https://eips.ethereum.org/EIPS/eip-1559
	toAddress := common.HexToAddress(data.To)
	cID := big.NewInt(int64(r.currentChainInfo.ID))

	var tx *types.Transaction
	// BSC uses legacy transaction type
	if r.currentChainInfo.ID == config.BscMainnet || r.currentChainInfo.ID == config.BscTestnet {
		tx = types.NewTransaction(nonce, toAddress, data.Value, gasLimit, data.PreferredBaseGasPrice, nil)
	} else {
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   cID,
			Nonce:     nonce,
			GasFeeCap: data.PreferredBaseGasPrice,
			GasTipCap: data.PreferredTipGasPrice,
			Gas:       gasLimit,
			To:        &toAddress,
			Value:     data.Value,
			Data:      nil,
		})
	}

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(int64(r.currentChainInfo.ID))), pk)
	if err != nil {
		return "", err
	}
	result := r.client.SendTransaction(context.Background(), signedTx)
	return signedTx.Hash().String(), result
}

func (r Relay) TransferToken(privateKey string, data *extTypes.TransactionRaw) (string, error) {
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	fromAddress, err := r.getAddressFromPrivateKey(pk)
	if err != nil {
		return "", err
	}

	nonce, err := r.getNonceFromAddress(*fromAddress)
	if err != nil {
		return "", err
	}

	transferFnSignature := []byte("transfer(address,uint256)") // do not include spaces in the string
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4] // 0xa9059cbb

	toAddress := common.HexToAddress(data.To)
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	amount := data.Value
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var inputData []byte
	inputData = append(inputData, methodID...)
	inputData = append(inputData, paddedAddress...)
	inputData = append(inputData, paddedAmount...)

	token := strings.ToLower(data.TokenSymbol)
	tokenAddress := r.supportTokens[token]
	if tokenAddress == "" {
		return "", errors.New("token not match any of supported tokens")
	}
	if !common.IsHexAddress(tokenAddress) {
		return "", errors.New("token address is not valid")
	}
	tAddress := common.HexToAddress(tokenAddress)

	gasLimit, err := r.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: *fromAddress,
		To:   &tAddress,
		Data: inputData,
	})
	if err != nil {
		return "", err
	}

	cID := big.NewInt(int64(r.currentChainInfo.ID))

	var tx *types.Transaction
	// BSC uses legacy transaction type
	if r.currentChainInfo.ID == config.BscMainnet || r.currentChainInfo.ID == config.BscTestnet {
		tx = types.NewTransaction(nonce, tAddress, big.NewInt(0), gasLimit, data.PreferredBaseGasPrice, inputData)
	} else {
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   cID,
			Nonce:     nonce,
			GasFeeCap: data.PreferredBaseGasPrice,
			GasTipCap: data.PreferredTipGasPrice,
			Gas:       gasLimit,
			To:        &tAddress,
			Value:     big.NewInt(0),
			Data:      inputData,
		})
	}

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(cID), pk)
	if err != nil {
		return "", err
	}

	err = r.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

func (r Relay) GasPrice() (*extTypes.EstimateGasInfo, error) {
	cID := r.currentChainInfo.ID
	price, pErr := r.client.SuggestGasPrice(context.Background())
	tip, tErr := r.client.SuggestGasTipCap(context.Background())
	// BSC網路取得不到 GasTip
	if pErr != nil || (tErr != nil && cID != config.BscMainnet && cID != config.BscTestnet) {
		return nil, errors.New("GasPrice failed")
	}
	return &extTypes.EstimateGasInfo{
		Base: price,
		Tip:  tip,
	}, nil
}

func (r Relay) GetBalance(address string) (balance *big.Int, err error) {
	balance, err = r.client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func (r Relay) GetBalanceForToken(address string, symbol string) (*big.Int, uint8, error) {
	tokenAddress := r.supportTokens[strings.ToLower(symbol)]
	if tokenAddress == "" {
		return nil, 0, errors.New("can not find matched token")
	}
	tAddr := common.HexToAddress(tokenAddress)
	instance, err := token.NewToken(tAddr, r.client)
	if err != nil {
		log.Fatal(err)
	}
	balance, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		log.Fatal(err)

	}
	decimal, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)

	}
	return balance, decimal, nil
}

// func (r Relay) GasLimit(symbol string) (uint64, error) {
// 	s := strings.ToLower(symbol)
// 	cID := r.currentChainInfo.ID

// 	nativeToken := false
// 	nativeToken = (cID == config.BscMainnet || cID == config.BscTestnet && s == "bnb")
// 	nativeToken = (cID == config.Rinkeby || cID == config.Mainnet && s == "eth")

// 	// standard transfer limit, see https://ethereum.org/en/developers/docs/gas/, https://eips.ethereum.org/EIPS/eip-1559
// 	// also apply to BSC
// 	if nativeToken {
// 		return 21000, nil
// 	}

// 	fromAddress := common.HexToAddress("0xE34224f746F7Da45c870573850d4AbbfC8c3B1AC")
// 	toAddress := common.HexToAddress("0xef92aF139cDAdE4A3cB89bb72839c78a1f7406A7")

// 	gasLimit, err := r.client.EstimateGas(context.Background(), ethereum.CallMsg{
// 		From: *fromAddress,
// 		To:   &tAddress,
// 		Data: inputData,
// 	})
// }

// ******** PRIVATE ******** //

func createInstance(c config.ChainInfo) (*Relay, error) {
	p := config.GetProviderInfo(c.ProviderFile)
	client, err := ethclient.Dial(p.URL)
	if err != nil {
		return nil, errors.New("ethclient dial failed")
	}

	return &Relay{
		c,
		p.Tokens,
		client,
	}, nil
}

func (r Relay) destory() {
	r.client.Close()
}

func (r Relay) getAddressFromPrivateKey(pk *ecdsa.PrivateKey) (*common.Address, error) {
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &fromAddress, nil
}

func (r Relay) getNonceFromAddress(address common.Address) (uint64, error) {
	nonce, err := r.client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

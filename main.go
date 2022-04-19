package main

import (
	"fmt"

	"gitlab.inlive7.com/crypto/ethereum-relay/internal/relay"
)

func main() {
	// config.InitializeConfiguration()
	testing()
}

func testing() {
	// str := "{\"types\":{\"EIP712Domain\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"version\",\"type\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\"}],\"Login\":[{\"name\":\"nonce\",\"type\":\"string\"},{\"name\":\"contents\",\"type\":\"string\"}]},\"primaryType\":\"Login\",\"domain\":{\"name\":\"XBBCrypto\",\"version\":\"1.0\",\"chainId\":\"4\",\"verifyingContract\":\"0x0000000000000000000000000000000000000000\"},\"message\":{\"contents\":\"Login to XBBCrypto\",\"nonce\":\"123123123123\"}}"
	// fmt.Println(str)
	// reveal, _ := api.VerifySignature(str, "0xc16edd7cea6c29b5c38c74552c92780790b8e5c4370a3e59eca699cbe67f307e12dfe77a90b8ef4a027a9408b10173c288a69ef99abc554f28580d070d6fe2151c")
	// reveal, _ := api.VerifySignature("0xa6011a27752e2b49267977350b30fbf60fd63d434413382cd2249ea4e4c6f906", "0xfebc2d140d7fee300a9be8de26a0e06b68e8bc3dcfbecf73b215c0b68d78c61e7083b9e89f6a88cccd506cc1596ffa2307ec0bfff5e86dfd8e43bf3371e3ccf71c")

	// reveal, _ := api.VerifySignature("123456", "0x395d73df806b470e2211deecb8a8568c8cf164f7fd283cc37038ebb0c814cbeb24e2a9fe1726482bee6c45530851d240a30ce67a971205f895baa1cc17aa30241b")
	// fmt.Println(reveal)
	// r, p, e := relay.Shared(4).QueryTransaction("0xfe49a399dc9f6ea5a41b7eb415767a22d01054ab70ec93da0010a9d8b3ad6731")
	// fmt.Println(r)
	// fmt.Println(p)
	// fmt.Println(e)
	// a, b, c := api.CreateNewAccount()
	// fmt.Println(a, b, c)

	balance, err := relay.Shared(97).GetBalance("0xE34224f746F7Da45c870573850d4AbbfC8c3B1AC")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)
	// info, decimal, err := relay.Shared(1).GetBalanceForToken("0xE34224f746F7Da45c870573850d4AbbfC8c3B1AC", "usdt")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(info, decimal)

	// fmt.Println(utility.StringWithoutExponent(utility.WeiToGwei(info.Base)))
	// fmt.Println(utility.StringWithoutExponent(utility.WeiToGwei(info.Tip)))

	// a, pend, err := relay.Shared(4).QueryTransaction("0xfb5759bd0b7b4664470ca8857340cc8be591c2150eef0d77d53408d64e76d346")
	// fmt.Println(a, pend, err)

	// suggestGas, _ := relay.Shared(4).GasPrice()
	// fmt.Println(suggestGas)
	// hash, err := relay.Shared(4).TransferValue(
	// 	"2b6c64b688e50a652dd4cf66e478f2fcae8539f0096e18de0d5ea90c0dec2047",
	// 	&types.TransactionRaw{
	// 		To:                    "0xE34224f746F7Da45c870573850d4AbbfC8c3B1AC",
	// 		Value:                 utility.Gwei(1),
	// 		PreferredBaseGasPrice: new(big.Int).Quo(suggestGas.Tip, big.NewInt(2)),
	// 		PreferredTipGasPrice:  suggestGas.Tip,
	// 	})

	// fmt.Println(hash, err)

	// fmt.Println(api.GetBalance(4, "0xE34224f746F7Da45c870573850d4AbbfC8c3B1AC"))

	// limit, err := api.GetGasLimit(56, "CTZN", "0xE34224f746F7Da45c870573850d4AbbfC8c3B1AC", "0xc63013d45d51bec40f84b8d4aa515faf9f5d88cb", big.NewInt(7e14)) //77931
	// fmt.Println(limit, err)

	// info, err := api.GetGasPrice(1)
	// fmt.Println(info, err)
}

type B struct {
	value int
}

type A struct {
	bbb B
}

func (a *A) test(value B) {
	a.bbb = value
}

func (a *A) check() {
	a.bbb = B{value: 9999}
}

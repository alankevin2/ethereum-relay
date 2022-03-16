package main

import (
	"fmt"

	"gitlab.inlive7.com/crypto/ethereum-relay/pkg/api"
)

func main() {
	// config.InitializeConfiguration()
	testing()
}

func testing() {
	// str := "{\"types\":{\"EIP712Domain\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"version\",\"type\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\"}],\"Login\":[{\"name\":\"nonce\",\"type\":\"string\"},{\"name\":\"contents\",\"type\":\"string\"}]},\"primaryType\":\"Login\",\"domain\":{\"name\":\"XBBCrypto\",\"version\":\"1.0\",\"chainId\":\"4\",\"verifyingContract\":\"0x0000000000=000000000000000000000000000000\"},\"message\":{\"contents\":\"Login to XBBCrypto\",\"nonce\":\"123123123123\"}}"
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

	// info, err := relay.Shared(1).GasPrice()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(info)
	// fmt.Println(utility.StringWithoutExponent(utility.WeiToGwei(info.Base)))
	// fmt.Println(utility.StringWithoutExponent(utility.WeiToGwei(info.Tip)))

	// result := relay.Shared(4).SendTransaction(
	// 	"16e89bd5528ca6ee27f321e00ec76e6c00ecc6f61a90ac86be3791da4cb7702d",
	// 	&relay.TransactionRaw{
	// 		To:                "0xef92aF139cDAdE4A3cB89bb72839c78a1f7406A7",
	// 		Value:             utility.Gwei(2),
	// 		PreferredGasPrice: utility.Gwei(2),
	// 	})

	// fmt.Println(result)

	fmt.Println(api.GetBalance(4, "0xf92af139cdade4a3cb89bb72839c78a1f7406a7"))
}

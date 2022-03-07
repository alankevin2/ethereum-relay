package main

import (
	"fmt"
	"math/big"

	"gitlab.inlive7.com/crypto/ethereum-relay/config"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/relay"
	"gitlab.inlive7.com/crypto/ethereum-relay/internal/utility"
	"gitlab.inlive7.com/crypto/ethereum-relay/pkg/api"
)

func main() {
	config.InitializeConfiguration()
	testing()
}

func testing() {
	reveal, _ := api.VerifySignature("asdasd", "0x063822ca173c4c2ea7c0af6c23d9eb9b1dc398c97cc216a835eb2a2d1d081fdb274159b86aabec4c21c9134ab0e1d44b09b7a456802470a01618919591f034331b")
	// reveal, _ := api.VerifySignature("123456", "0x395d73df806b470e2211deecb8a8568c8cf164f7fd283cc37038ebb0c814cbeb24e2a9fe1726482bee6c45530851d240a30ce67a971205f895baa1cc17aa30241b")
	fmt.Println(reveal)
	// r, p, e := relay.Shared(4).QueryTransaction("0xfe49a399dc9f6ea5a41b7eb415767a22d01054ab70ec93da0010a9d8b3ad6731")
	// fmt.Println(r)
	// fmt.Println(p)
	// fmt.Println(e)
	// a, b, c := api.CreateNewAccount()
	// fmt.Println(a, b, c)

	info, err := relay.Shared(1).GasPrice()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
	fmt.Println(utility.StringWithoutExponent(utility.WeiToGwei(info.Base)))
	fmt.Println(utility.StringWithoutExponent(utility.WeiToGwei(info.Tip)))

	relay.Shared(1).SendTransaction("2b6c64b688e50a652dd4cf66e478f2fcae8539f0096e18de0d5ea90c0dec2047", "0xE34224f746F7Da45c870573850d4AbbfC8c3B1AC", big.NewInt(500))
}

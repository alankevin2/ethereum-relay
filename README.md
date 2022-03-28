# Solidity

### 1. 用solidity compiler輸出abi給go解析使用

    solc --overwrite --abi contracts/ERC20ByteCodeGenerator.sol -o contracts/dist/


### 2. 用solidity compiler輸出bin給EVM bytecode使用

    solc --overwrite --bin contracts/ERC20ByteCodeGenerator.sol -o contracts/dist/

### 3. 將智能合約輸出成go檔

    ~/go/bin/abigen --bin=contracts/dist/ERC20ByteCodeGenerator.bin --abi=contracts/dist/ERC20ByteCodeGenerator.abi --pkg=token --out=contracts/dist/ERC20Token.go 

#### P.S.
solc 透過這樣裝
```
brew update
brew tap ethereum/ethereum
brew install solidity
```

abigen則是
```
go get -u github.com/ethereum/go-ethereum
cd $(go env GOPATH)/pkg/github.com/ethereum/go-ethereum/
make
make devtools
```
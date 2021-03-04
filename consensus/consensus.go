package consensus

import (
	"2021/_03_公链/XianFengChain04/chain"
	"math/big"
)

/*
*定义区块链中使用的共识算法 接口
 */
type Consensus interface {
	FindNonce() int64
}

func NewPow(block chain.Block) Consensus {
	initTarget := big.NewInt(1)
	initTarget.Lsh(initTarget,255-DIFFICULTY)
	return PoW{
		Block:  block,
		Target: initTarget,
	}
}





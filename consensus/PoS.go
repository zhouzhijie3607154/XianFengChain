package consensus

import (
	"2021/_03_公链/XianFengChain04/chain"
	"fmt"
)

type PoS struct {
	Block chain.Block
}

func(pos PoS)FindNonce(block chain.Block)int64  {
	fmt.Println("使用pos算法进行共识机制")
}


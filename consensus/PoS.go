package consensus

import (
	"fmt"
)

type PoS struct {
	Block BlockInterface
}

func(pos PoS)FindNonce(block BlockInterface)int64  {
	fmt.Println("使用pos算法进行共识机制")
	return 0
}


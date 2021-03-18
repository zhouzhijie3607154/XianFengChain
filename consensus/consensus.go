package consensus

import (
	"2021/_03_公链/XianFengChain04/transaction"
	"math/big"
)

/*
*定义区块链中使用的共识算法 接口
 */
type Consensus interface {
	FindNonce() (int64 ,[32]byte)
}
/**
*定义区块结构体的接口标准
 */
type BlockInterface interface {
	GetHeight() int64
	GetVersion() int64
	GetTimeStamp() int64
	GetPreHash()[32]byte
	GetTransactions() []transaction.Transaction
}

func NewPow(block BlockInterface) Consensus {
	initTarget := big.NewInt(1)
	initTarget.Lsh(initTarget,255-DIFFICULTY)
	return PoW{
		Block:  block,
		Target: initTarget,
	}
}





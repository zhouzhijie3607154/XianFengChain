package consensus

import (
	"2021/_03_公链/XianFengChain04/chain"
	"2021/_03_公链/XianFengChain04/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

const DIFFICULTY = 10

type PoW struct {
	Block  chain.Block
	Target *big.Int
}

func (pow PoW) FindNonce() int64 {
	//给定一个nonce值，计算区块hash
	var nonce int64
	nonce = 0
	//无限循环
	for {

		hash := CalculateHash(pow.Block,nonce)

		//拿到系统目标值
		target := pow.Target
		//比较大小
		if bytes.Compare(hash[:], target.Bytes()) == -1 {
			return nonce
		}
		nonce++
	}

}
func CalculateHash(block chain.Block,nonce int64)[32]byte  {
	heightByte,_ := utils.Int2Byte(block.Height)
	versionByte,_ :=utils.Int2Byte(block.Version)
	timeByte,_ := utils.Int2Byte(block.TimeStamp)
	nonceByte,_:= utils.Int2Byte(nonce)
	blockByte :=bytes.Join([][]byte{heightByte,versionByte,timeByte,nonceByte,block.Data,block.PrevHash[:]},nil)
	block.Hash = sha256.Sum256(blockByte)
	//计算区块的hash
	return sha256.Sum256(blockByte)
}

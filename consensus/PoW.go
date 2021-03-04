package consensus

import (
	"2021/_03_公链/XianFengChain04/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

const DIFFICULTY = 20

type PoW struct {
	Block  BlockInterface
	Target *big.Int
}

func (pow PoW) FindNonce() int64 {
	//给定一个nonce值，计算区块hash
	var nonce int64
	nonce = 0
	//无限循环
	hashBig := new(big.Int)
	for {
		hash := CalculateHash(pow.Block,nonce)
		//拿到系统目标值
		target := pow.Target
		//比较大小
		hashBig = hashBig.SetBytes(hash[:])
		if hashBig.Cmp(target)== -1{
			return nonce
		}
		nonce++
	}

}
//计算区块并返回hash值
func CalculateHash(block BlockInterface,nonce int64)[32]byte  {
	heightByte,_ := utils.Int2Byte(block.GetHeight())
	versionByte,_ :=utils.Int2Byte(block.GetVersion())
	timeByte,_ := utils.Int2Byte(block.GetTimeStamp())
	nonceByte,_:= utils.Int2Byte(nonce)
	preHash := block.GetPreHash()
	blockByte :=bytes.Join([][]byte{heightByte,versionByte,timeByte,nonceByte,block.GetData(), preHash[:]},nil)
	return sha256.Sum256(blockByte)
}

package chain

import (
	"2021/_03_公链/XianFengChain04/consensus"
	"time"
)

/**
*区块的结构体定义 ：版上默时难随
 */
type Block struct {
	Height   int64    //高度
	Version  int64    //版本号
	PrevHash [32]byte //上一区块hash
	Hash     [32]byte //本区块hash
	//默克尔根
	TimeStamp int64 //时间戳
	//Difficulty int64
	Nonce int64 //随机数
	Data  []byte
}


/**
*创建创世区块
 */
func CreateGenesis(data []byte) Block {
	gensis := Block{
		Height:    0,
		Version:   0x00,
		PrevHash:  [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		TimeStamp: time.Now().Unix(),
		Data:      data,
	}
	//调用PoW共识算法，寻找随机数，计算哈希值
	proof := consensus.NewPow(gensis)
	gensis.Nonce,gensis.Hash = proof.FindNonce()

	return gensis
}

/**
*生成新区块的功能函数
 */
func NewBlock(height int64, prev [32]byte, data []byte) Block {

	block := Block{
		Height:    height + 1,
		Version:   0x00,
		PrevHash:  prev,
		TimeStamp: time.Now().Unix(),
		Data:      data,
	}
	//调用PoW共识算法，寻找随机数，计算哈希值
	proof := consensus.NewPow(block)
	block.Nonce,block.Hash = proof.FindNonce()
	return block
}

//实现接口 BlockInterface
func (block Block) GetHeight() int64 {
	return block.Height
}
func (block Block) GetVersion() int64 {
	return block.Version
}
func (block Block) GetTimeStamp() int64 {
	return block.TimeStamp
}
func (block Block) GetPreHash() [32]byte {
	return block.PrevHash
}
func (block Block) GetData() []byte {
	return block.Data
}

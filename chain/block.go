package chain

import "time"

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

func CreateGenesis(data []byte) Block {
	gensis := NewBlock(0, [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil)
	gensis.Height = 0
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

	//todo 设置哈希，寻找并设置随机数
	return block
}

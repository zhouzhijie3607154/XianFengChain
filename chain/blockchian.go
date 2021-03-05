package chain

//定义区块链结构体
type BlockChain struct {
	Blocks []Block
}

//创建一个区块链对象，初始化一个创世区块
func CreateChainWithGensis(data []byte) BlockChain {
	genesis := CreateGenesis(data)
	blocks := make([]Block, 0, )
	blocks = append(blocks, genesis)

	return BlockChain{Blocks: blocks}
}

//创建一个新区块，并添加到区块链中
func (chain *BlockChain) CreateNewBlock(data []byte) {
	//获取当前区块链上最新的区块
	lastBlock := chain.Blocks[len(chain.Blocks)-1]
	//新建一个区块
	newBlock := NewBlock(lastBlock.Height,lastBlock.Hash,data)
	//将区块添加到区块链上
	chain.Blocks = append(chain.Blocks,newBlock )
}

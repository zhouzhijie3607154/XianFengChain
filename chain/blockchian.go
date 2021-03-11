package chain

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"math/big"
)

const BLOCKS = "blocks"
const LASTHASH = "lasthash"

//定义区块链结构体
type BlockChain struct {
	//Blocks []Block
	DB        *bolt.DB
	LastBlock Block
	IteratorBlockHash [32]byte//迭代器当前迭代到的区块的Hash
}

func CreateChain(db *bolt.DB) BlockChain {
	var lastBlock Block
	db.Update(func(tx *bolt.Tx) error {
		bucket:= tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			bucket, _= tx.CreateBucket([]byte(BLOCKS))
		}
		lastHash:=bucket.Get([]byte(LASTHASH))
		if len(lastHash)<= 0 {
			return nil
		}
		//从文件当中读取出最新的区块，并赋值给chain.LastBlock
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, _= DeSerialize(lastBlockBytes)
		return nil

	})
	return BlockChain{
		DB:        db,
		LastBlock: lastBlock,
		IteratorBlockHash:lastBlock.Hash,
	}
}

//创建一个区块链对象，初始化一个创世区块
func (chain *BlockChain) CreateChainWithGensis(data []byte) error {
	flag := new(big.Int).SetBytes(chain.IteratorBlockHash[:]).Cmp(big.NewInt(0)) == 1
	if flag {
		return nil
	}
	var err error
	err = chain.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return  errors.New("操作区块链数据文件异常！")
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) == 0 {
			/*bucket已存在
			*key -> value
			*blockHash --> block序列化后的数据
			*lasthash -->最新区块的hash
			 */
			genesis := CreateGenesis(data)
			genesisBytes, err := genesis.Serialize()
			if err != nil {
				return err
			}
			bucket.Put(genesis.Hash[:], genesisBytes)
			bucket.Put([]byte(LASTHASH), genesis.Hash[:])
			chain.LastBlock = genesis
			chain.IteratorBlockHash = genesis.Hash
		} 
		//else {
			////从文件当中读取出最新的区块，并赋值给chain.LastBlock
			//lastHash := bucket.Get([]byte(LASTHASH))
			//lastBlockBytes := bucket.Get(lastHash)
			//chain.LastBlock, err = DeSerialize(lastBlockBytes)
			//chain.IteratorBlockHash = chain.LastBlock.Hash
			//if err != nil {
			//	return err
			//}
		//}

		return nil
	})
	return err
}

//创建一个新区块，并添加到区块链中
func (chain *BlockChain) CreateNewBlock(data []byte) error {
	/**
	*生成一个新区块，并存到bolt。BD文件中，由于涉及到存新区块，所以我们这里选择Updata
		1、从文件中查到当前存储的最新区块数据
		2、反序列化得到区块
		3、根据最新获取的区块生成一个新的区块
		4、将新区块序列化，存储到boltDB文件，同时更新最新区块hash
	*/
	err := chain.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块诗句库操作失败，请重试！")
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, err := DeSerialize(lastBlockBytes)
		if err != nil {
			return err
		}
		newBlock := NewBlock(lastBlock.Height, lastBlock.Hash, data)
		newBlockBytes, err := newBlock.Serialize()
		if err != nil {
			return err
		}
		err = bucket.Put(newBlock.Hash[:], newBlockBytes)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(LASTHASH), newBlock.Hash[:])
		if err != nil {
			return err
		}
		chain.LastBlock = newBlock
		chain.IteratorBlockHash = newBlock.Hash
		return nil
	})
	return err
}

/*
获取区块链上最新的区块
@lastblock ：最新区块数据
@err ：可能遇到的错误
*/
func (chain *BlockChain) GetLastBlock() Block {
	return chain.LastBlock
}

/**
*获取所有区块数据
1、找到最后一个区块
2、通过区块的PreHash依次找上一个区块，直至创世区块
3、每次找到的区块添加到blocks中，找到创世区块（读取完所有区块后返回blocks）
@return
	blocks: 从db文件中读取到的所有区块
	err ：读取过程中遇到的错误
*/
func (chain *BlockChain) GetAllBlocks() (blocks []Block, err error) {
	blocks = make([]Block, 0)
	err = chain.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			fmt.Println(1)
			return errors.New("blockchain.go的方法GetAllBlocks错误")
		}
		var currentHash []byte
		currentHash = bucket.Get([]byte(LASTHASH))
		if err != nil {
			return err
		}
		//for循环找所有区块
		for {
			currentBlockBytes := bucket.Get(currentHash)
			currentBlock, err := DeSerialize(currentBlockBytes)
			if err != nil {
				return err
			}
			blocks = append(blocks, currentBlock)
			currentHash = currentBlock.PrevHash[:]
			if currentBlock.Height == 0 {
				break
			}

		}
		return nil
	})
	return blocks, err
}

//实现接口：iterator的 HasNext方法，判断是否还有数据 true、false
func (chain *BlockChain) HasNext() (hasNext bool) {
	/*
		1、当前区块在哪？->preHash ->db
	*/

	db := chain.DB
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("HasNext：区块数据文件操作失败，请重试")
		}
		preBlockBytes := bucket.Get(chain.IteratorBlockHash[:])
		hasNext = len(preBlockBytes) != 0
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return hasNext
}

//该方法用于实现iterator的Next方法，取出一个block
func (chain *BlockChain) Next() (block Block) {
	/*
		1、知道当前在哪个区块-->找当前区块的上一个区块--->将找到的区块返回
	*/

	err := chain.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("BlockChain.Next：区块数据文件操作失败，请重试")
		}
		iteratorBlockBytes := bucket.Get(chain.IteratorBlockHash[:])
		iteratorBlock, err := DeSerialize(iteratorBlockBytes)
		if err != nil {
			return err
		}
		chain.IteratorBlockHash = iteratorBlock.PrevHash
		block = iteratorBlock
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	return block
}

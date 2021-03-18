package chain

import (
	"2021/_03_公链/XianFengChain04/transaction"
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
	DB                *bolt.DB
	LastBlock         Block
	IteratorBlockHash [32]byte //迭代器当前迭代到的区块的Hash
}

func CreateChain(db *bolt.DB) BlockChain {
	var lastBlock Block
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) <= 0 {
			return nil
		}
		//从文件当中读取出最新的区块，并赋值给chain.LastBlock
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, _ = DeSerialize(lastBlockBytes)
		return nil

	})
	return BlockChain{
		DB:                db,
		LastBlock:         lastBlock,
		IteratorBlockHash: lastBlock.Hash,
	}
}

//创建一个区块链对象，初始化一个创世区块
func (chain *BlockChain) CreateChainWithGenesis(txs []transaction.Transaction) error {
	flag := new(big.Int).SetBytes(chain.IteratorBlockHash[:]).Cmp(big.NewInt(0)) == 1
	if flag {
		return nil
	}
	var err error
	err = chain.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("操作区块链数据文件异常！")
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) == 0 {
			/*bucket已存在
			*key -> value
			*blockHash --> block序列化后的数据
			*lasthash -->最新区块的hash
			 */
			genesis := CreateGenesis(txs)
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
func (chain *BlockChain) CreateNewBlock(txs []transaction.Transaction) error {
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
			return errors.New("区块数据库操作失败，请重试！")
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, err := DeSerialize(lastBlockBytes)
		if err != nil {
			return err
		}
		newBlock := NewBlock(lastBlock.Height, lastBlock.Hash, txs)
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

/* 2021年3月17日14:39:05
改方法用于查询出指定地址的 UTXO 集合并返回
*/
func (chain *BlockChain) SearchUTXO(addr string) []transaction.UTXO {
	//准备两个容器
	spend := make([]transaction.TxInput, 0)  //花费记录的容器
	inCome := make([]transaction.UTXO, 0) //收入记录的容器

	//遍历拿到每一个区块
	for chain.HasNext() {
		block := chain.Next()
		//遍历区块中的每一个交易
		for _, tx := range block.Transactions {
			//遍历交易中的每一个交易输入
			for _, input := range tx.Inputs {
				if string(input.ScriptSig) != addr {
					continue
				}
				spend = append(spend, input)

			}
			//遍历交易中的每一个交易输出
			for index, output := range tx.OutPuts {
				if string(output.ScriptPub) != addr { //与当前遍历地址无关直接跳过
					continue
				}
				//与当前地址相等则创建一个记录（input）记录该交易ID 下标
				utxo := transaction.UTXO{
					TxId: tx.TxHash,
					Vout: index,
					TxOutPut: output,
				}
				inCome = append(inCome,utxo )

			}
		}
	}
	utxos := make([]transaction.UTXO, 0)
	//将收入与花费两个容器中的相同input进行比较，找出未花费的收入
	var isComeSpent bool
	for _, come := range inCome {
		for _, spen := range spend {
			if come.TxId == spen.TxId && come.Vout == spen.Vout {
				isComeSpent = true
				break
			}
		}
		//如果该笔收入找遍所有的花费记录都没有匹配到，即
		if !isComeSpent {
			utxos = append(utxos, come)
		}
	}
	return utxos
}

//定义区块链的发送交易的功能
func (chain *BlockChain) SendTransaction(from, to string, amount float64) error {
	//1、先把from中的可花费的钱（utxo）找出来
	utxos := chain.SearchUTXO(from)
	var totalBalance float64
	for _,utxo := range utxos{
		totalBalance += utxo.Value
	}
	if totalBalance  < amount{
		return errors.New("余额不足，交易失败！！！")
	}
	//2、可花费的钱总额大于要花费的amount，才构建交易

	newTx, err := transaction.CreateNewTransaction(from, to, amount)
	if err != nil {
		return err
	}
	err = chain.CreateNewBlock([]transaction.Transaction{*newTx})
	if err != nil {
		return err
	}
	return nil

}

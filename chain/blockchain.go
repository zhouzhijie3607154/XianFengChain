package chain

import (
	"2021/_03_公链/XianFengChain04/transaction"
	"2021/_03_公链/XianFengChain04/wallet"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"math/big"
)

const BLOCKS = "blocks"
const LASTHASH = "lasthash"

/**
 * 定义区块链结构体，该结构体用于管理区块
 */
type BlockChain struct {
	//切片
	//[block1,block2,block3]
	//Blocks []Block
	DB        *bolt.DB
	LastBlock Block
	//cursor游标
	//current
	IteratorBlockHash [32]byte //表示当前迭代到了那个区块，该变量用于记录迭代到的区块hash

	Wallet *wallet.Wallet //引入钱包管理功能
}

func CreateChain(db *bolt.DB) (*BlockChain, error) {
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
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, _ = Deserialize(lastBlockBytes)
		return nil
	})
	wallet, err := wallet.LoadAddressFromDB(db)
	if err != nil {
		fmt.Println("初始化钱包失败,", err.Error())
		return nil, err
	}
	return &BlockChain{
		DB:                db,
		LastBlock:         lastBlock,
		IteratorBlockHash: lastBlock.Hash,
		Wallet:            wallet,
	}, nil
}

/**
 * 创建coinbase交易的方法
 */
func (chain *BlockChain) CreateCoinBase(addr string) error {
	//参数检查 : 对用户传递的地址进行有效性地址检查
	isAddrValid := chain.Wallet.CheckAddress(addr)
	if !isAddrValid {
		return errors.New("输入的地址不合法,请检查后重试!")
	}

	//1、创建一笔coinbase交易
	coinbase, err := transaction.CreateCoinBase(addr)
	if err != nil {
		return err
	}

	//2、把coinbase交易存到区块中
	err = chain.CreateGensis([]transaction.Transaction{*coinbase})

	return err
}

/**
 * 创建一个区块链对象，包含一个创世区块
 */
func (chain *BlockChain) CreateGensis(txs []transaction.Transaction) error {
	hashBig := new(big.Int)
	hashBig.SetBytes(chain.LastBlock.Hash[:])
	//最新区块hash有值，则说明区块文件中创世区块已经存在了
	if hashBig.Cmp(big.NewInt(0)) == 1 {
		return nil
	}

	var err error
	//gensis持久化到db中去
	engine := chain.DB
	engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil { //没有桶
			bucket, err = tx.CreateBucket([]byte(BLOCKS))
			if err != nil {
				return err
				//panic("操作区块存储文件失败,请重试！")
			}
		}
		//先查看
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) == 0 { //第一次
			gensis := CreateGenesis(txs)
			genSerBytes, _ := gensis.Serialize()
			//bucket已经存在
			// key -> value
			// blockHash -> 区块序列化以后的数据
			bucket.Put(gensis.Hash[:], genSerBytes) //把创世区块保存到boltdb中去
			//使用一个标志，用来记录最新区块的hash，以标明当前文件中存储到了最新的哪个区块
			bucket.Put([]byte(LASTHASH), gensis.Hash[:])
			//把geneis赋值给chain的lastblock
			chain.LastBlock = gensis
			chain.IteratorBlockHash = gensis.Hash
		}
		return nil
	})
	return err
}

/**
 * 生成一个新区块
 */
func (chain *BlockChain) CreateNewBlock(txs []transaction.Transaction) error {
	//目的：生成一个新区块，并存到bolt.DB文件中去(持久化）
	//手段（步骤）：
	//1、从文件中查到当前存储的最新区块数据
	lastBlock := chain.LastBlock
	//3、根据获取的最新区块生成一个新区块
	newBlock := NewBlock(lastBlock.Height, lastBlock.Hash, txs)
	//4、将最新区块序列化，得到序列化数据
	var err error
	newBlockSerBytes, err := newBlock.Serialize()
	if err != nil {
		return err
	}
	//5、将序列化数据存储到文件、同时更新最新区块的标记lasthash，更新为最新区块的hash
	db := chain.DB
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块数据库操作失败，请重试!")
		}
		//将新生成的区块保存到文件中
		bucket.Put(newBlock.Hash[:], newBlockSerBytes)
		//更新最新区块的标记lasthash，更新为最新区块的hash
		bucket.Put([]byte(LASTHASH), newBlock.Hash[:])
		//更新内存中的blockchain的LastBlock
		chain.LastBlock = newBlock
		chain.IteratorBlockHash = newBlock.Hash
		return nil
	})
	return err
}

//获取最新的区块数据
func (chain *BlockChain) GetLastBlock() Block {
	return chain.LastBlock
}

//获取所有的区块数据
func (chain *BlockChain) GetAllBlocks() ([]Block, error) {
	//目的：获取所有的区块
	//手段（步骤）：
	//1、找到最后一个区块
	db := chain.DB
	var err error
	blocks := make([]Block, 0)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块数据库操作失败,请重试！")
			return err
		}
		var currentHash []byte
		currentHash = bucket.Get([]byte(LASTHASH))
		//2、根据最后一个区块依次往前找
		for {
			currentBlockBytes := bucket.Get(currentHash)
			currentBlock, err := Deserialize(currentBlockBytes)
			if err != nil {
				break
			}
			//3、每次找到的区块放入到一个[]Block容器中
			blocks = append(blocks, currentBlock)
			//4、找到最开始的创世区块时，就结束了，不再找了
			if currentBlock.Height == 0 {
				break
			}
			currentHash = currentBlock.PrevHash[:]
		}
		return nil
	})
	return blocks, err
}

/**
 * 该方法用于实现迭代器Iterator的HasNext方法,用于判断是否还有数据
 * 如果有数据，返回true，否则返回false
 */
func (chain *BlockChain) HasNext() bool {
	//区块0 -> 区块1 -> 区块2 -> 区块3
	//最新区块: 区块3
	//步骤: 当前区块在哪 ->  preHash -> db
	db := chain.DB
	var hasNext bool
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败，请重试！")
		}
		//获取前一个区块的区块数据
		blockBytes := bucket.Get(chain.IteratorBlockHash[:])
		//如果获取不到前一个区块的数据，说明前面没有区块了
		hasNext = len(blockBytes) != 0
		return nil
	})
	return hasNext
}

/**
 * 该方法用于实现迭代器Iterator的Next方法，用于取出一个数据
 * 此处，因为是区块链数据集合，因此返回的数据类型是Block
 */
func (chain *BlockChain) Next() Block {
	//1、知道当前在哪个区块
	//2、找当前区块的前一个区块
	//3、找到的区块返回
	db := chain.DB
	var iteratorBlock Block
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败,请重试！")
		}
		//前一个区块的数据
		blockBytes := bucket.Get(chain.IteratorBlockHash[:])
		iteratorBlock, _ = Deserialize(blockBytes)
		//迭代到当前区块后，更新游标的区块内容
		chain.IteratorBlockHash = iteratorBlock.PrevHash
		return nil
	})
	return iteratorBlock
}

/**
 * 该方法用于查询出指定地址的UTXOs集合并返回
 */
func (chain *BlockChain) SearchUTXOsFromDB(addr string) []transaction.UTXO {

	//花费记录的容器
	spend := make([]transaction.TxInput, 0)

	//收入记录的容器
	inCome := make([]transaction.UTXO, 0)

	//迭代遍历每一个区块
	for chain.HasNext() {
		block := chain.Next()
		//遍历区块中的交易
		for _, tx := range block.Transactions {
			//a、遍历每个交易的交易输入
			for _, input := range tx.Inputs {
				//找到了花费记录
				if input.VerifyInputWithAddress(addr) {
					spend = append(spend, input)
				}
			}
			//b、遍历每个交易的交易输出:收入的钱记录下来
			for index, output := range tx.Outputs {
				if output.CheckPubKHashWithAddr(addr) {
					utxo := transaction.NewUTXO(tx.TxHash, index, output)
					inCome = append(inCome, utxo)
				}
			}
		}
	}

	//遍历收入集合和花费集合，把已花的剔除，找出未花费的记录
	utxos := make([]transaction.UTXO, 0)
	var isComeSpent bool
	for _, come := range inCome {
		isComeSpent = false
		for _, spen := range spend {
			if come.IsUTXOSpent(spen) {//该笔收入已被花费
				isComeSpent = true
				break
			}
		}
		if !isComeSpent { //当前遍历到的come没有被消费
			utxos = append(utxos, come)
		}
	}
	return utxos
}

/**
 * 该方法用于实现地址余额的统计
 */
func (chain *BlockChain) GetBalance(addr string) (float64, error) {
	//检查地址的合法性
	isAddrValid := chain.Wallet.CheckAddress(addr)
	if !isAddrValid {
		return 0.0000, errors.New("非法地址")
	}

	_, totalBalance := chain.GetUTXOsWithBalance(addr, []transaction.Transaction{})
	return totalBalance, nil
}

/**
 * 该方法用于实现地址余额统计和地址所可以花费的utxo集合
 */
func (chain BlockChain) GetUTXOsWithBalance(addr string, txs []transaction.Transaction) ([]transaction.UTXO, float64) {
	//1、遍历bolt.DB文件，找区块中的可用的utxo的集合
	dbUtxos := chain.SearchUTXOsFromDB(addr)

	//2、找一遍内存中已经存在但还未存到文件中的交易
	memSpends := make([]transaction.TxInput, 0)
	memInComes := make([]transaction.UTXO, 0)
	for _, tx := range txs {
		//a、遍历交易输入，把花出去的钱记录下来
		for _, input := range tx.Inputs {
			if input.VerifyInputWithAddress(addr) {// input 属于 addr
				memSpends = append(memSpends, input)
			}
		}
		//b、遍历交易输出，把得到的钱记录下来
		for outIndex, output := range tx.Outputs {
			if output.CheckPubKHashWithAddr(addr) { // output 属于 addr
				utxo := transaction.NewUTXO(tx.TxHash, outIndex, output)
				memInComes = append(memInComes, utxo)
			}
		}
	}
	//3、经过内存中的交易的遍历以后，剩下的才是最终可用的utxo集合
	utxos := make([]transaction.UTXO, 0)
	var isUTXOSpend bool

	for _, utxo := range dbUtxos {//遍历dbUTXO集合,看看内存中是否有已经花费了的
		isUTXOSpend = false

		for _, spend := range memSpends {
			if utxo.IsUTXOSpent(spend) {
				isUTXOSpend = true
			}
		}

		if !isUTXOSpend {
			utxos = append(utxos, utxo)
		}
	}
	//把内存中的交易输出也加入到可用的utxo集合中
	utxos = append(utxos, memInComes...)

	var totalBalance float64
	for _, utxo := range utxos {
		totalBalance += utxo.Value
	}
	return utxos, totalBalance
}

/**
 * 定义区块链的发送交易的功能
 * @param froms 交易发送者的地址切片
 * @param tos 交易收款者的地址切片
 * @param amount 交易的金额切片
 */
func (chain *BlockChain) SendTransaction(froms []string, tos []string, amounts []float64) error {
	//1.地址合法性检查,非法地址直接返回一个错误
	for i := 0; i < len(froms); i++ {
		isFormValid := chain.Wallet.CheckAddress(froms[i])
		isToValid := chain.Wallet.CheckAddress(tos[i])
		if !isFormValid || !isToValid {
			return errors.New("地址不合法,请检查后重试!")
		}
	}

	//2.创建一个交易切片,用于存储生成的交易
	newTxs := make([]transaction.Transaction, 0) //内存

	//3.遍历交易发起者列表,拿到下标 from_index 每个切片的长度都相同,所以通用
	for from_index, from := range froms {
		//3.1 先把from的可花费的utxos给找出来
		utxos, totalBalance := chain.GetUTXOsWithBalance(from, newTxs)

		//3.2 如果from的钱不够,终止交易
		if totalBalance < amounts[from_index] {
			return errors.New(from + "余额不足，赶紧去搬砖挣钱")
		}

		//3.3 如果from的钱够用,判断需要哪些钱? 切片中多少钱就够本次交易花费
		totalBalance = 0
		var utxoNum int
		for index, utxo := range utxos {
			totalBalance += utxo.Value
			if totalBalance > amounts[from_index] {
				utxoNum = index
				break
			}
		}

		//3.4取出交易发起者的私钥,用于创建一笔新的交易 newTx
		fromKeyPair := chain.Wallet.GetKeyPairByAddr(from)

		newTx, err := transaction.CreateNewTransaction(utxos[0:utxoNum+1],
			from,
			tos[from_index],
			fromKeyPair.Pub,
			amounts[from_index])
		if err != nil {
			return err
		}

		//3.4.1 对签名验证,只有签名通过,才能将
		err = newTx.SignTx(fromKeyPair.Priv,utxos[0:utxoNum+1])
		if err != nil {
			return err
		}
		//3.5新的交易 newTx 添加到 之前创建的 交易切片
		newTxs = append(newTxs, *newTx)
	}

	// 4. 创建区块,把交易切片扔进去
	err := chain.CreateNewBlock(newTxs)
	if err != nil {
		return err
	}

	return nil
}

//创建新地址
func (chain *BlockChain) GetNewAddress() (string, error) {
	return chain.Wallet.NewAddress()
}

//导出指定地址的密钥对
func (chain *BlockChain) DumpPrivateKey(addr string) (*wallet.KeyPair, error) {
	return chain.Wallet.DumpKeyPair(addr)
}

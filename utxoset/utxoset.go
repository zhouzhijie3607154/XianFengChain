package utxoset

import (
	"2021/_03_公链/XianFengChain04/transaction"
	"2021/_03_公链/XianFengChain04/utils"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

const UTXOSET = "utxoset" //存放utxoset的桶名

// UTXOset集合,用于存放区块链的所有UTXO 实现快速查询
// map 中 以 地址为 key  以utxo切片为值
type UTXOSet struct {
	//UTXOs  map[string][]transaction.UTXO
	Engine *bolt.DB
}

func LoadUTXOSet(db *bolt.DB) *UTXOSet  {
	return &UTXOSet{Engine: db}
}

//查询某个地址的所有可用UTXO集合,统计的utxo金额,和一个可能遇到的err
func (utxoset UTXOSet) QueryUTXOsByAddr(address string) (utxos []transaction.UTXO,totalBalance float64,err error) {
	fmt.Println("queryUTXO 收到的参数为: ",address)
	utxos = make([]transaction.UTXO, 0)
	err = utxoset.Engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UTXOSET))
		if bucket == nil {
			return errors.New("抱歉,utxoset 桶未创建")
		}

		utxosBytes := bucket.Get([]byte(address))
		if len(utxosBytes) <= 0 {
			return nil
		} //未查到该地址的UTXO数据

		decoder := gob.NewDecoder(bytes.NewReader(utxosBytes))
		err = decoder.Decode(&utxos)
		fmt.Println("解码后的utxo有多少张: ",len(utxos))
		for i,v := range utxos{
			fmt.Printf("第 %d 张utxo 的 金额为 : %f\n",i,v.Value)
		}

		return err
	})
	for _, utxo := range utxos {
		totalBalance += utxo.Value
	}
	return utxos, totalBalance, err
}

//当某个地址有新的 UTXO时,把新产生的UTXO存入到UTXOSet中
func (utxoset UTXOSet) AddUTXOsWithAddr(address string, newUTXOs []transaction.UTXO) (bool, error) {
	var err error
	var utxosAll = make([]transaction.UTXO, 0) //用户所有的 UTXO
	utxoset.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UTXOSET))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(UTXOSET))
			if err != nil {
				return err
			}
		}

		utxosBytes := bucket.Get([]byte(address))
		// 该address 之前已存有 UTXO在set中,不能直接覆盖
		if len(utxosBytes) != 0 {
			decoder := gob.NewDecoder(bytes.NewReader(utxosBytes))
			err = decoder.Decode(&utxosAll)
		}
		//把新产生的utxo追加到该地址的utxo切片中,再序列化存到db文件中
		utxosAll := append(utxosAll, newUTXOs...)
		utxoBytes, err := utils.Encode(utxosAll)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(address), utxoBytes)

		return err
	})

	return err == nil, err
}

//当某个地址使用了UTXO时,把使用了的UTXO从  UTXOSet中删除
func (utxoset UTXOSet) DeleteUTXOsWithAddr(address string, records []SpendRecord) (bool, error) {
	var err error
	//账户已有的utxo
	var utxoExsited = make([]transaction.UTXO, 0)
	err = utxoset.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UTXOSET))
		if bucket == nil {
			return errors.New("抱歉,utxoset 桶未创建")
		}

		utxosBytes := bucket.Get([]byte(address))
		if len(utxosBytes) == 0 {
			return errors.New("该用户db里没钱!  账户地址为:" + address)
		}
		decoder := gob.NewDecoder(bytes.NewReader(utxosBytes))
		err = decoder.Decode(&utxoExsited)
		if err != nil {
			return err
		}

		//判断records与找到的utxos的关系,records必须是utxos的子集,消费的必须是已有的钱
		isSub := IsSubUTXOs(utxoExsited, records)
		if !isSub {
			return errors.New("本地未找到所消费的utxo")
		}

		//删除之后账户的 utxo
		remainUTXOs := make([]transaction.UTXO, 0)
		for _, record := range records {
			isSpent := false
			for _, utxoE := range utxoExsited {
				if record.EqualUTXO(utxoE) { //如果已有的utxo出现在消费记录中,他就被消费了
					isSpent = true
				}

				if !isSpent { //保留没删掉的
					remainUTXOs = append(remainUTXOs, utxoE)
				}
			}

		}

		remainUTXOsBytes, err := utils.Encode(remainUTXOs)

		bucket.Put([]byte(address), remainUTXOsBytes)

		return err
	})

	return err == nil, err
}

func IsSubUTXOs(utxos []transaction.UTXO, records []SpendRecord) bool {
	for _, record := range records {

		isContains := false

		for _, utxo := range utxos {
			if record.EqualUTXO(utxo) {
				isContains = true
			}
			if !isContains {
				return false
			}
		}
	}
	return true
}

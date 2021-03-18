/**
2021年3月12日10:35:59
author：admin
filename ：transaction.go
*/
package transaction

import (
	"2021/_03_公链/XianFengChain04/utils"
	"crypto/sha256"
)

/**
*该结构体定义交易 分交易输入和交易输出
*
 */
type Transaction struct {
	TxHash  [32]byte   //交易的整体数据hash 用作唯一标识
	Inputs  []TxInput  //交易输入
	OutPuts []TxOutPut //交易输出

}

/*
创建coinbase交易
*/
func CraeteCoinbase(address string) (*Transaction, error) {
	output0 := TxOutPut{
		Value:     50,
		ScriptPub: []byte(address),
	}
	coinbase := Transaction{
		TxHash:  [32]byte{},
		Inputs:  nil,
		OutPuts: []TxOutPut{output0},
	}
	bytes, err := utils.GobEncode(coinbase)
	if err != nil {
		return nil, err
	}
	coinbase.TxHash = sha256.Sum256(bytes)
	return &coinbase, nil
}

/*
该函数用于构建一笔普通的交易，返回构件好的交易实例
*/
func CreateNewTransaction(from, to string, amout float64) (*Transaction, error) {
	/*
		1.构建inputs
		2、构建outputs
		3、构建transaction
		4、计算trans哈希，赋值给Txhash
	*/
	inputs := make([]TxInput, 0)
	outputs := make([]TxOutPut, 0)
	//定义一个变量，记录form一共付了多少钱
	var inputAmmout float64
	outPut0 := TxOutPut{
		Value:     amout,
		ScriptPub: []byte(to),
	}
	outputs = append(outputs, outPut0)
	//判断是否需要找零？
	if (inputAmmout - amout) > 0 {
		outPut1 := TxOutPut{
			Value:     inputAmmout- amout,
			ScriptPub: []byte(from),
		}
		outputs =  append(outputs, outPut1)
	}
//交易输入本质上就是一个交易输出的引用？
/*
如何找到那笔交易、引用那一笔交易输出、引用哪一笔交易输出的索引？
	遍历区块，遍历交易，遍历交易输出，筛选交易输出-->得到有关from的交易输出的集合
如何找到没有花费出去的交易输出？
	遍历区块，遍历交易，遍历交易输入，筛选交易输入-->得到有关 from的交易输入的集合
		form的UTXO = 输出的集合 - 输入的集合
 */

	newTransaction := Transaction{
		TxHash:  nil,
		Inputs:  inputs,
		OutPuts: outputs,
	}
	txbytes, err := utils.GobEncode(newTransaction)
	if err != nil {
		return nil, err
	}
	newTransaction.TxHash = sha256.Sum256(txbytes)
	return &newTransaction, nil
}

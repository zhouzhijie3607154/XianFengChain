package transaction

import (
	"2021/_03_公链/XianFengChain04/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

const REWARDSIZE = 50

/**
 * 定义交易的结构体
 */
type Transaction struct {
	//交易哈希
	TxHash [32]byte
	//交易输入
	Inputs []TxInput
	//交易输出
	Outputs []TxOutput
}

/**
 * 该函数用于定义一个coinbase交易，并返回该交易结构体
 */
func CreateCoinBase(addr string) (*Transaction, error) {
	//1.由系统一笔交易输出
	output0 := NewTxOutput(REWARDSIZE, addr)
	//2.系统生成 一笔 coinbase 交易
	coinbase := Transaction{
		Outputs: []TxOutput{*output0},
	}

	//3.对coinbase交易进行 gob 序列化,生成交易 Id 并赋值
	coinbaseBytes, err := utils.Encode(coinbase)
	if err != nil {
		return nil, err
	}
	coinbase.TxHash = sha256.Sum256(coinbaseBytes)

	return &coinbase, nil
}

/**
 * 该函数用于构建一笔普通的交易，返回构建好的交易实例

 */
func CreateNewTransaction(utxos []UTXO, from, to string, pubk []byte, amount float64) (*Transaction, error) {

	//1、构建inputs
	inputs := make([]TxInput, 0) //用于存放交易输入的容器
	var inputAmount float64      //该变量用于记录转账发起者一共付了多少钱
	//input -> 交易输入本质上对某个交易的交易输出UTXO的引用
	for _, utxo := range utxos {
		input := NewTxInput(utxo.TxId, utxo.Vout, pubk)

		inputAmount += utxo.Value

		//把构建好的input存入到交易输入容器中
		inputs = append(inputs, input)
	}

	//2、构建outputs
	outputs := make([]TxOutput, 0) //用于存放交易输出的容器
	//构建转账接收者的交易输出
	output0 := NewTxOutput(amount, to)

	outputs = append(outputs, *output0) //把第一个交易输出放入到专门存交易输出的容器中

	//判断是否需要找零,如果需要找零，则需要构建一个新的找零输出
	if inputAmount-amount > 0 {
		output1 := NewTxOutput(inputAmount-amount, from)
		outputs = append(outputs, *output1)
	}

	//3、构建transaction
	newTransaction := Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}

	//4、计算transaction的哈希,并赋值
	transactionBytes, err := utils.Encode(newTransaction)
	newTransaction.TxHash = sha256.Sum256(transactionBytes)
	if err != nil {
		return nil, err
	}
	//5、将构建的transaction实例进行返回
	return &newTransaction, nil
}

/*
对交易进行签名
*/
func (tx *Transaction) SignTx(priv *ecdsa.PrivateKey, utxos []UTXO) (err error) {
	if len(tx.Inputs) != len(utxos) {
		err = errors.New("签名失败")
		return err
	}

	//复制一份 交易 txCopy 在副本中进行签名,避免签名过程中原始数据被修改
	txCopy := tx.CopyTx()

	//在复制品中签名,签名的结果复制给原始 tx
	for i := 0; i < len(txCopy.Inputs); i++ {
		//1.遍历得到每一笔交易输入
		input := txCopy.Inputs[i]

		//2.遍历得到每一个 UTXO
		utxo := utxos[i]

		//3.将交易输入中的原始公钥赋值为 公钥hash
		input.PubK = utxo.PubKHash

		//4.计算交易的哈希值 txHash
		txHash, err := txCopy.CalculateTxHash()
		if err != nil {
			return err
		}
		//5.使用私钥对得到的txHash进行签名
		r, s, err := ecdsa.Sign(rand.Reader, priv, txHash)
		if err !=nil {
			return err
		}
		//6.对原始交易中的  input中的Sig字段进行赋值
		tx.Inputs[i].Sig = append(r.Bytes(),s.Bytes()...)

		//7.最后清空 交易输入中的 原始公钥
		input.PubK = nil
	}
	return nil
}

//拷贝交易对象实例
func (tx Transaction) CopyTx() Transaction {
	copyTransaction := tx
	return copyTransaction
}


//计算交易的哈希值
func (tx *Transaction)CalculateTxHash()([]byte,error)  {
	txBytes, err := utils.Encode(tx)
	return utils.Hash256(txBytes),err
}

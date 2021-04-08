package transaction

import (
	"2021/_03_公链/XianFengChain04/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
)

const REWARDSIZE = 50

//定义交易的结构体
type Transaction struct {
	//交易哈希
	TxHash [32]byte
	//交易输入
	Inputs []TxInput
	//交易输出
	Outputs    []TxOutput
	LockedTime int64 //时间戳,确保生成的交易的hash值唯一
}

//该函数用于定义一个coinbase交易，并返回该交易结构体
func CreateCoinBase(addr string) (*Transaction, error) {
	//1.由系统一笔交易输出
	output0 := NewTxOutput(REWARDSIZE, addr)
	//2.系统生成 一笔 coinbase 交易
	coinbase := Transaction{
		Outputs:    []TxOutput{*output0},
		LockedTime: time.Now().Unix(),
	}

	//3.对coinbase交易进行 gob 序列化,生成交易 Id 并赋值
	coinbaseBytes, err := utils.Encode(coinbase)
	if err != nil {
		return nil, err
	}
	coinbase.TxHash = sha256.Sum256(coinbaseBytes)

	return &coinbase, nil
}

//该函数用于构建一笔普通的交易，返回构建好的交易实例
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
		Inputs:     inputs,
		Outputs:    outputs,
		LockedTime: time.Now().Unix(),
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

//对交易进行签名
func (tx *Transaction) SignTx(priv *ecdsa.PrivateKey, utxos []UTXO) (err error) {
	//如果是coinbase交易,不需要签名
	if tx.IsCoinBase() {
		return nil
	}

	if len(tx.Inputs) != len(utxos) {
		err = errors.New("签名失败")
		return err
	}

	//复制一份 交易 txCopy 在副本中进行签名,避免签名过程中原始数据被修改
	txCopy := tx.CopyTx()

	//在复制品中签名,签名的结果复制给原始 tx
	for index, _ := range txCopy.Inputs {
		//1.将交易输入中的原始公钥赋值为 公钥hash
		txCopy.Inputs[index].PubK = utxos[index].PubKHash

		//2.计算交易的哈希值 txHash
		txHash, err := txCopy.CalculateTxHash()
		if err != nil {
			return err
		}
		fmt.Printf("签名时的交易hash数据: %x\n ", txHash)

		//3.使用私钥对得到的txHash进行签名
		r, s, err := ecdsa.Sign(rand.Reader, priv, txHash)
		if err != nil {
			return err
		}

		//4.对原始交易中的  input中的Sig字段进行赋值
		tx.Inputs[index].Sig = append(r.Bytes(), s.Bytes()...)

		//5.最后清空 交易输入中的 原始公钥
		txCopy.Inputs[index].PubK = nil
		fmt.Printf("签名后的PubK信息: %x\n ", tx.Inputs[index].PubK)

	}
	return nil
}

//交易的签名验证方法,该方法返回一个布尔值,ture为验证通过,否则签名不通过
func (tx *Transaction) VerifyTx(utxos []UTXO) (bool, error) {
	//如果是coinbase交易,不需要验证
	if tx.IsCoinBase() {
		return true, nil
	}

	if len(tx.Inputs) != len(utxos) {
		fmt.Println(len(tx.Inputs), len(utxos))
		return false, errors.New("txInputs length should equal utxos length")
	}
	// 验签 : 需要 公钥 签名结果 原始数据 -->  ecdsa.Verify 函数调用
	//已有的 公钥 : tx.Input[i].PubK
	//   签名数据 : tx.Input[i].Sig
	//签名时是对 整个交易数据进行签名 ,sign字段也在交易中
	txCopy := tx.CopyTx()
	for index, _ := range txCopy.Inputs {

		//a. 清空 签名结果
		txCopy.Inputs[index].Sig = nil

		//b. 副本数据改回签名前的状态
		txCopy.Inputs[index].PubK = utxos[index].PubKHash

		//对副本数据 计算 改造后的交易 hash
		txHash, err := txCopy.CalculateTxHash()
		if err != nil {
			return false, err
		}

		//调用 api进行签名验证
		//公钥格式转换 : []byte --> PublicKey
		pubKey := utils.PubBytesToEcdsaPubKey(elliptic.P256(), tx.Inputs[index].PubK)

		//签名结果 : []byte -- > r,s *big.Int
		r, s := utils.SignBytesToSignature(tx.Inputs[index].Sig)

		//签名验证
		isVerify := ecdsa.Verify(pubKey, txHash, r, s)

		if !isVerify {
			return false, errors.New("签名验证失败..")
		}
	}
	return true, nil
}

//拷贝交易对象实例
func (tx Transaction) CopyTx() Transaction {
	/*	inputs := make([]TxInput, 0)
		for _, input := range tx.Inputs {
			txIn := TxInput{
				TxId: input.TxId,
				Vout: input.Vout,
				PubK: input.PubK,
				Sig:  input.Sig,
			}
			inputs = append(inputs, txIn)
		}

		outputs := make([]TxOutput, 0)
		for _, output := range tx.Outputs {
			txOut := TxOutput{
				Value:    output.Value,
				PubKHash: output.PubKHash,
			}
			outputs = append(outputs, txOut)
		}
		hash := tx.TxHash

		return Transaction{
			TxHash:  hash,
			Inputs:  inputs,
			Outputs: outputs,
		}

	*/
	outputs := make([]TxOutput, 0)
	inputs := make([]TxInput, 0)
	copy(tx.Outputs, outputs)
	copy(tx.Inputs, inputs)
	return Transaction{
		TxHash:  tx.TxHash,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

//计算交易的哈希值
func (tx *Transaction) CalculateTxHash() ([]byte, error) {
	txBytes, err := utils.Encode(tx)
	return utils.Hash256(txBytes), err
}

//判断是否为coinbase交易
func (tx Transaction) IsCoinBase() bool {
	if len(tx.Inputs) == 0 && len(tx.Outputs) == 1 {
		return true
	}
	return false
}

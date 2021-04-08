package transaction

import (
	"2021/_03_公链/XianFengChain04/utils"
	"2021/_03_公链/XianFengChain04/wallet"
	"bytes"
)

//定义结构体UTXO，表示未花费的交易输出
type UTXO struct {
	TxId     [32]byte //该笔收入来自哪个交易
	Vout     int      //该笔收入来自交易的哪个输出
	TxOutput          //该表输入的面额和收入者
}

//创建一个UTXO 的结构实例
func NewUTXO(txId [32]byte, vout int, out TxOutput) UTXO {

	return UTXO{
		TxId:     txId,
		Vout:     vout,
		TxOutput: out,
	}
}

//验证该 UTXO 是否已经某笔输入 引用被消费
func (utxo *UTXO) IsUTXOSpent(spend TxInput) bool {
	/* utxo 的交易哈希(TxId) 交易索引(Vout)  公钥哈希(PubKHash) 与交易输入都一致时
	则说明该笔 utxo 已经被消费,返回 true  ,一个不同则返回false(未被消费)
	*/
	//1.判断 交易hash是否一致
	equalTxId := bytes.Compare(utxo.TxId[:], spend.TxId[:]) == 0

	//2.判断索引下标
	equalVout := utxo.Vout == spend.Vout

	//3.判断 utxo.PubKHash(公钥哈希) 和 spend.PubK(原始公钥)
	pubk := spend.PubK
	hash256 := utils.Hash256(pubk)
	ripemd160 := utils.HashRipemd160(hash256)
	pubkHash := append(wallet.VERSION, ripemd160...)

	equalPubKHash := bytes.Compare(utxo.PubKHash, pubkHash) == 0

	return equalTxId && equalVout && equalPubKHash
}

////判断该utxo 与给定 的 UTXO是否为同一个
//func (utxo *UTXO) EqualUTXO(record utxoset.SpendRecord) bool {
//	equalTxID := bytes.Compare(utxo.TxId[:], record.TxId[:]) == 0
//	equalVout := utxo.Vout == record.Vout
//
//	return  equalTxID  && equalVout
//}

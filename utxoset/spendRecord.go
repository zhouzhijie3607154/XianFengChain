package utxoset

import "2021/_03_公链/XianFengChain04/transaction"

/**
* 该结构体用于定义一笔交易的消费记录
 */
type SpendRecord struct {
	TxId [32]byte
	Vout int
}

//判断该消费记录是否引用了指定的 utxo
func (record *SpendRecord) EqualUTXO(utxo transaction.UTXO) bool {
	isEqualVout := record.Vout == utxo.Vout
	isEqualTxId := record.TxId == utxo.TxId
	return isEqualTxId && isEqualVout
}

//创建一个新的SpendRecord
func NewSpendRecord(txid [32]byte,vout int)SpendRecord  {
	return SpendRecord{
		TxId: txid,
		Vout: vout,
	}
}

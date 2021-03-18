package transaction

/***
* 	value   金额
*	scrpitPub 锁定脚本
*
 */
type TxOutPut struct {
	Value     float64
	ScriptPub []byte
}

//定义结构体UTXO，表示未花费的交易输出
type UTXO struct {
	TxId     [32]byte //该笔收入来自哪个交易
	Vout     int      //该笔收入来自交易的哪个输出
	TxOutPut          //该笔收入的面额和收入者

}

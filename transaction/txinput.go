package transaction

/*
*交易id 	TxId
*交易输出id  Vout
*解锁脚本	ScriptSig
*
*
*
 */
type TxInput struct {
	TxId      [32]byte
	Vout      int
	ScriptSig []byte
}

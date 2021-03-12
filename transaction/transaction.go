/**
2021年3月12日10:35:59
author：admin
filename ：transaction.go
*/
package transaction

/**
*该结构体定义交易 分交易输入和交易输出
*
 */
type Transaction struct {
	TxHash  [32]byte	//交易的整体数据hash 用作唯一标识
	Inputs  []TxInput  //交易输入
	OutPuts []TxOutPut //交易输出

}

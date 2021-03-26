package transaction

import (
	"bytes"
	"github.com/mr-tron/base58"
)

/**
 * 定义交易输出的结构体
 */
type TxOutput struct {
	Value     float64 //转账的数量
	//ScriptPub []byte  //锁定脚本 2021年3月25日10:36:23 拆为 PubKHash
	PubKHash []byte	 //公钥Hash
}
/*
	 ScriptSig:解锁脚本: <sig> <pubK>
	 ScriptPubK:锁定脚本: DUP HASH160 <PubKHash> EQUALVERIFY CHECKSIG
	操作指令 :  DUP HASH160  EQUALVERIFY CHECKSIG
	 数据 :  <sig> <pubK> <PubKHash>
	 指令栈 : 压栈 和 弹栈
	 指令 : 3 2 ADD 5 EQUAL

*/
//生成一笔交易输出
func NewTxOutput(value float64,addr string)(*TxOutput)  {
	//1.对 addr 进行 base58 反编码
	reAddr,_ := base58.Decode(addr)
	//2.去除校验位,得到公钥hash
	pubHash := reAddr[:len(reAddr) -4]
	out := TxOutput{
		Value:    value,
		PubKHash: pubHash,
	}
	return &out
}
/*
验证某个交易输出是否属于某个地址 属于返回true 不属于返回false
 */
func(output *TxOutput) CheckPubKHashWithAddr(addr string)bool  {
	//1.对 addr 进行 base58 反编码
	reAddr,_ := base58.Decode(addr)
	//2.去除校验位,得到公钥hash
	pubHash := reAddr[:len(reAddr) -4]

	//3.比较给定addr 的公钥哈希是否与 交易输出的公钥哈希 是否相等
	return bytes.Compare(output.PubKHash,pubHash) == 0
}

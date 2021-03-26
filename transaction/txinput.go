package transaction

import (
	"2021/_03_公链/XianFengChain04/utils"
	"2021/_03_公链/XianFengChain04/wallet"
	"bytes"
	"github.com/mr-tron/base58"
)

/**
 * 定义交易输入的结构体
 */
type TxInput struct {
	TxId [32]byte //该字段确定引用自哪笔交易
	Vout int      //该字段确定引用自该交易的哪个输出
	//ScriptSig []byte   //该字段表示使用交易输出的证明，解锁脚本
	//于2021年3月25日10:41:25 拆分为 Sig 和 PubK两个字段

	Sig  []byte //
	PubK []byte
}

/*
	 ScriptSig:解锁脚本: <sig> <pubK>
	 ScriptPubK:锁定脚本: DUP HASH160 <PubKHash> EQUALVERIFY CHECKSIG
	操作指令 :  DUP HASH160  EQUALVERIFY CHECKSIG
	 数据 :  <sig> <pubK> <PubKHash>
	 指令栈 : 压栈 和 弹栈
	 指令 : 3 2 ADD 5 EQUAL

*/
/* 生成一笔新的交易输入 */
func NewTxInput(txid [32]byte, vout int, pubk []byte) TxInput {
	return TxInput{
		TxId: txid,
		Vout: vout,
		//Sig:  nil,
		PubK: pubk,
	}
}

/*验证某个 TxInput 是否是某个特定 地址的消费  已消费返回true 否则false*/
func (input *TxInput) VerifyInputWithAddress(address string) bool {
	// 方法一 : input.PubK 变换计算得到 addr 与 address 进行比较

	// 方法二 : input.PubK 变换计算得到 pubKHash  与 address 变换的PubHash2进行比较

	// 1. input.PubK 变换计算得到 pubKHash
	pubk := input.PubK
	hash256 := utils.Hash256(pubk)
	ripemd160 := utils.HashRipemd160(hash256)
	pubHash := append(wallet.VERSION, ripemd160...)

	//2. address 变换的PubHash1
	reAddress, _ := base58.Decode(address)
	rePubKHash := reAddress[:len(reAddress)-4]

	//input.PubK 变换计算得到 pubKHash1  与 address 变换的 PubHash2 进行比较
	return bytes.Compare(pubHash,rePubKHash)==0
}

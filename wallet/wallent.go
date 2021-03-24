package wallet

import (
	"2021/_03_公链/XianFengChain04/utils"
	"bytes"
	"crypto/sha256"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

/**用于
*定义 wallet结构体,用于管理地址和对应的密钥对信息

(Address)公钥生成的地址 --> 一对密钥(公钥,私钥)
*/
type Wallet struct {
	Address map[string]*KeyPair
}
/*
初始化 Wallet
 */



/**
  新生成一个比特币的地址,
*/
func (wallet *Wallet)NewAddress() (string, error) {
	// 1~2.生成一对密钥对
	keyPair, err := NewKeyPair()
	if err != nil {
		return "", err
	}

	//3、对公钥进行sha256哈希
	pubHash := sha256.Sum256(keyPair.Pub)
	//4、reipemd160 计算
	ripe := ripemd160.New()
	ripe.Write(pubHash[:])
	ripemdPub := ripe.Sum(nil)
	//5、追加版本号
	versionPub := append([]byte{0x00}, ripemdPub...)

	//6、两次哈希
	firstHash := utils.Hash256(versionPub)
	secondHash := utils.Hash256(firstHash)

	//7.截取前四个字节作为地址校验位
	check := secondHash[:4]

	//8.拼接到versionPub后面
	originAddress := append(versionPub, check...)

	//9.base58编码
	address := base58.Encode(originAddress)

	//10.生成的新地址添加到钱包中
	wallet.Address[address] = keyPair

	return address, nil
}

//该函数用于检查地址是否合法,返回一个bool类型的值,合法返回true 否则false
func(wallet *Wallet) CheckAddress(addr string) bool {
	if len(addr) <= 8 {
		return false
	}

	//1.使用base58对传入的地址进行解码
	reAddrBytes, _ := base58.Decode(addr)

	//2.取出校验位
	reCheck := reAddrBytes[len(reAddrBytes)-4:]

	//3.取出 versioinPubHash
	reVersionPubHash := reAddrBytes[:len(reAddrBytes)-4]

	//4.把reversionPubHash双哈希后取前四个字节与reCheck进行校验
	reFirstHash := utils.Hash256(reVersionPubHash)
	reSecondHash := utils.Hash256(reFirstHash)

	if bytes.Compare(reSecondHash[:4], reCheck) == 0 { //如果校验位相同,返回true
		return true
	}
	return false
}

package wallet

import (
	"2021/_03_公链/XianFengChain04/utils"
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

/**用于
*定义 wallet结构体,用于管理地址和对应的密钥对信息

(Address)公钥生成的地址 --> 一对密钥(公钥,私钥)
(Engine)持久化钱包内的地址
*/
type Wallet struct {
	Address map[string]*KeyPair
	Engine  *bolt.DB
}

//存地址的桶名
var KEYSTORE = "KeyStore"

//存地址的 Key
var ADDRESSS_KEY = "addressKey"

/**
  新生成一个比特币的地址,
*/
func (wallet *Wallet) NewAddress() (string, error) {
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
	err = wallet.SaveAddressToDB()
	if err != nil {
		fmt.Println("保存地址出现错误",err.Error())
	}

	return address, err
}

//该函数用于检查地址是否合法,返回一个bool类型的值,合法返回true 否则false
func (wallet *Wallet) CheckAddress(addr string) bool {
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

/*
钱包中的地址持久化到DB文件中
*/
func (wallet *Wallet) SaveAddressToDB()(error) {
	var err error
	err = wallet.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		//如果桶为空,创建桶
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(KEYSTORE))
			if err != nil {
				return err
			}

		}
		//gob 注册接口 编码方式
		gob.Register(elliptic.P256())
		buff := new(bytes.Buffer)
		encoder := gob.NewEncoder(buff)
		err = encoder.Encode(wallet.Address)
		if err != nil {
			fmt.Println("gob加密出错了",err.Error())
			return err
		}
		err= bucket.Put([]byte(ADDRESSS_KEY), buff.Bytes())
		//todo
		fmt.Println("存了啥进去呀",wallet.Address)
		return err
	})
	return err
}

/*
加载DB文件中的所有钱包地址
*/
func LoadAddressFromDB(db *bolt.DB)(*Wallet,error){
	var err error
	var address = make(map[string]*KeyPair)
		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(KEYSTORE))
			//如果KeyStore桶不存在,创建KeyStore桶
			if bucket == nil {
				return errors.New("都没桶呢,你在捣鼓啥?")
			}

			//如果有KeyStore桶,读取数据
			addressBytes := bucket.Get([]byte(ADDRESSS_KEY))
			//如果桶中有数据...取出address...反序列化
			if len(addressBytes) >= 0 {
				//gob注册接口 编码
				gob.Register(elliptic.P256())
				decoder := gob.NewDecoder(bytes.NewReader(addressBytes))
				err = decoder.Decode(&address)
			}

		return err
	})
	if err !=nil {
		return nil,err
	}
	//创建Wallet
	return &Wallet{
		Address: address,
		Engine:  db,
	},nil
}
/**
查询钱包中的所有地址
 */
func(wallet *Wallet)GetAddressList()([]string){
	addressList := make([]string, 0)

	for address,_ := range wallet.Address{
		addressList = append(addressList, address)
	}
	return addressList
}


// 1EpRnkdAYinLhmFhJUvonZYnHwdVQ359fp
// 1JQrgBshyApwBReEoCM4XSVbiPbTFVSWy1
// 12K2rdAgZpzC9TFvLbfwRgDfjrgHrTVusS

/** 存储所有地址到DB文件时  map中的address值太离散了, 不放便取值 (Address) ,改为 map 直接序列化 : gob注册接口 Curve
*for address, keyPair := range wallet.Address {
				err = encoder.Encode(keyPair)

				keypairBytes := bucket.Get([]byte(address))
				if len(keypairBytes) == 0 {
					bucket.Put([]byte(address), buff.Bytes())
				}
				buff.Reset() //清空buffer  复用
			}
*/

/**		查询db文件中的所有地址,不需要了
var err error
var address = make(map[string]*KeyPair)
err = wallet.Engine.View(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte(KEYSTORE))
	//如果KeyStore桶不存在,创建KeyStore桶
	if bucket == nil {
		return err
	}

	//如果有KeyStore桶,读取数据
	addressBytes := bucket.Get([]byte(ADDRESSS_KEY))
	if len(addressBytes) != 0 {//如果桶中有数据...取出address...反序列化
		//gob注册接口 编码
		gob.Register(elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(addressBytes))
		err = decoder.Decode(&address)
	}
	return err
})*/

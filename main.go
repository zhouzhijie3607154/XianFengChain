package main

import (
	"2021/_03_公链/XianFengChain04/chain"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/prometheus/common/log"
)

func main() {
	blockChain := chain.CreateChainWithGensis([]byte("hello world！"))
	blockChain.CreateNewBlock([]byte("hello world ，too！"))
	fmt.Println(len(blockChain.Blocks))
	fmt.Printf("%+v\n",blockChain)
	bytes, err := blockChain.Blocks[0].Serialize()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	block, err := chain.DeSerialize(bytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("输出序列化与反序列化后的blockData：\n%s\n",block.Data)
	db, err := bolt.Open()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err:= tx.CreateBucket([]byte("hello world"))
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(""), []byte(""))
		data := bucket.Get([]byte(""))

		return nil
	})
}

package main

import (
	"2021/_03_公链/XianFengChain04/chain"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const BLOCKS = "xianfengchain04.db"

func main() {
	db, err := bolt.Open(BLOCKS, 0600, nil)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	blockchain := chain.CreateChain(db)
	err = blockchain.CreateChainWithGensis([]byte("hello world"))
	if err != nil {
		log.Println("main.go第20行：",err.Error())
	}
	err = blockchain.CreateNewBlock([]byte("hello"))
	if err != nil {
		log.Println("main.go第24行",err.Error())
	}
	lastBlock := blockchain.GetLastBlock()

	fmt.Printf("最新区块：%+v\n",lastBlock)
	blocks,err := blockchain.GetAllBlocks()
	if err  != nil {
		log.Println("main.go 第34行",err.Error())
	}
	for _,v :=range blocks{
		fmt.Printf("第%d个区块，%s\n",v.Height,v.Data)
	}

	//blockChain := chain.CreateChainWithGensis([]byte("hello world！"))
	//blockChain.CreateNewBlock([]byte("hello world ，too！"))
	//fmt.Println(len(blockChain.Blocks))
	//fmt.Printf("%+v\n",blockChain)
	//bytes, err := blockChain.Blocks[0].Serialize()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//block, err := chain.DeSerialize(bytes)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//fmt.Printf("输出序列化与反序列化后的blockData：\n%s\n",block.Data)
	//err = db.Update(func(tx *bolt.Tx) error {
	//	bucket, err:= tx.CreateBucket([]byte("hello world"))
	//	if err != nil {
	//		return err
	//	}
	//	err = bucket.Put(blockChain.Blocks[0].Hash[:],bytes)
	//	data := bucket.Get(blockChain.Blocks[0].Hash[:])
	//
	//	return nil
	//})
}

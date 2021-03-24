package main

import (
	"2021/_03_公链/XianFengChain04/chain"
	"2021/_03_公链/XianFengChain04/client"
	"fmt"
	"github.com/boltdb/bolt"
)

const BLOCKS = "xianfengchain04.db"

func main() {
	//打开数据库文件
	db, err := bolt.Open(BLOCKS, 0600, nil)
	if err !=nil {
		fmt.Println("打开数据文件失败",err.Error())
		return
	}

	defer db.Close() //xxx.db.lock
	//创建一条区块链
	blockChain ,err:= chain.CreateChain(db)
	if err !=nil {
		fmt.Println("创建区块链失败",err.Error())
		return
	}

	//创建一个 client
	cmdClient := client.CmdClient{*blockChain}

	cmdClient.Run()
}



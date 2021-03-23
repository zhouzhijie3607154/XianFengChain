package main

import (
	"2021/_03_公链/XianFengChain04/chain"
	"2021/_03_公链/XianFengChain04/client"
	"github.com/boltdb/bolt"
)

const BLOCKS = "xianfengchain04.db"

func main() {

	//打开数据库文件
	db, err := bolt.Open(BLOCKS, 0600, nil)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close() //xxx.db.lock
	blockChain := chain.CreateChain(db)
	cmdClient := client.CmdClient{blockChain}

	cmdClient.Run()
}

package main

import (
	"2021/_03_公链/XianFengChain04/chain"
	"fmt"
	"github.com/boltdb/bolt"
)

const BLOCKS = "xianfengchain04.db"

func main() {
	db, err := bolt.Open(BLOCKS, 0600, nil)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	blockchain := chain.CreateChain(db)
	err = blockchain.CreateChainWithGenesis([]byte("data"))
	if err != nil {
		fmt.Println(err.Error())
	}
	blockchain.CreateNewBlock([]byte("你好"))
	blocks, err := blockchain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _,v := range blocks{
		fmt.Println(v.Height)
	}
	//cmdline:=client.CmdClient{Chain:blockchain}
	//cmdline.Run()


}

package main

import (
	"2021/_03_公链/XianFengChain04/chain"
	"2021/_03_公链/XianFengChain04/client"
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
	//coinbase, err := transaction.CraeteCoinbase("123456")
	if err !=nil {
		fmt.Println(err)
		return
	}
	//err = blockchain.CreateChainWithGenesis([]transaction.Transaction{*coinbase})
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//blocks, err := blockchain.GetAllBlocks()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//for _,v := range blocks{
	//	fmt.Println(v)
	//}
	cmdline:=client.CmdClient{Chain:blockchain}
	cmdline.Run()


}

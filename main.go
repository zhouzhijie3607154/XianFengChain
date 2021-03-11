package main

import (
	"2021/_03_公链/XianFengChain04/chain"
	"2021/_03_公链/XianFengChain04/client"
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
	cmdline:=client.CmdClient{Chain:blockchain}
	cmdline.Run()
}

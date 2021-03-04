package main

import (
	"2021/_03_公链/XianFengChain04/chain"
	"fmt"
)

func main() {
	fmt.Println("Hello World !")
	block0 := chain.CreateGenesis([]byte("Hello World !"))
	block1 := chain.NewBlock(block0.Height, block0.Hash, []byte("OK！"))

	fmt.Printf("%+v",block0)
	fmt.Printf("%+v",block1)
}

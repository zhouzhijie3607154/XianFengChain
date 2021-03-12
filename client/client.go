/**
*CreateDate:2021年3月11日09:44:16
*FileName：client.go
*Author：Zhou
*
*该结构体另一了用于实行命令行参数解析的结构体
 */

package client

import (
	"2021/_03_公链/XianFengChain04/chain"
	"flag"
	"fmt"
	"math/big"
	"os"
)

type CmdClient struct {
	Chain chain.BlockChain
}
//命令行交互功能
func (cmd *CmdClient) Run() {
	if len(os.Args)==1 {//无参数、直接调帮助文档
		cmd.Help()
	}else {
	switch os.Args[1] {//有参数，调用参数对应的方法
	case GENERATEGENESIS:
		if len(os.Args[2:])>2 {
			fmt.Println("only  -genesis")
			return
		}
		var genesis string
		flagSet := flag.NewFlagSet(GENERATEGENESIS, flag.ExitOnError)
		flagSet.StringVar(&genesis,"genesis","","用户输入的创世区块从存储数据")
		flagSet.Parse(os.Args[2:])

		flag := new(big.Int).SetBytes(cmd.Chain.LastBlock.Hash[:]).Cmp(big.NewInt(0)) == 1
		if flag {
			fmt.Println("创世区块已存在！")
		}else {
			err := cmd.Chain.CreateChainWithGensis([]byte(genesis))
			if err != nil {
				fmt.Println(err.Error())
			}else {
				fmt.Println("generate genesis block success ! genesis block's data is",genesis)
			}
		}

	case CREATEBLOCK:
		if len(os.Args[2:])>2 {
			fmt.Println("参数无法识别：",os.Args[2:])
			return
		}
		var data string
		flagSet := flag.NewFlagSet(CREATEBLOCK, flag.ExitOnError)
		flagSet.StringVar(&data,"data","","用户输入的新区块数据")
		flagSet.Parse(os.Args[2:])

	//flag 为true时，说明 已经有创世区块了。
		flag := new(big.Int).SetBytes(cmd.Chain.LastBlock.Hash[:]).Cmp(big.NewInt(0)) == 1
		if !flag {
			fmt.Println("Error: ,your client don't have a blockchain now !")
			fmt.Println("Use generategenesis [-genesis 'data'] to create a new blockchain ")
		}
		err := cmd.Chain.CreateNewBlock([]byte(data))
		if err != nil {
			fmt.Println("main.go: ",err.Error())

		}else {
			fmt.Println("main.go:" ,"generate a new block OK and data is",data)
		}
	case GETLASTBLOCK:
		if len(os.Args[1:])>1 {
			fmt.Println("参数无法识别：",os.Args[1:])
			return
		}
		IsChain := new(big.Int).SetBytes(cmd.Chain.LastBlock.Hash[:]).Cmp(big.NewInt(0)) == 1
		if !IsChain {
			fmt.Println("Error: ,your client don't have a blockchain now !")
			fmt.Println("Use generategenesis [-genesis 'data'] to create a new blockchain ")
			return
		}
		fmt.Println("调用获取最新区块功能！")
		flag.String(GETLASTBLOCK,"","用于获取最新区块的数据")
		fmt.Printf("%+v\n",cmd.Chain.LastBlock)

	case GETALLBLOCKS:
		fmt.Println("调所有区块的功能！")
	case HELP:
		fmt.Println("调用帮助说明功能！")
	case VERSION:
		fmt.Println("调用获取版本信息！")
	default:
		fmt.Println("go run main.go: Unknown subcommand.")
		fmt.Println("Run 'go run main.go -help' for usage")
	}}
}

func (cmd *CmdClient) Help() {
	fmt.Println("---------------Welcome to XianFengChain04 Application------------")
	fmt.Println("XianfengChain04 is a custom BlockChain Project,the project plan to ")
	fmt.Println("")
	fmt.Println("USAGE")
	fmt.Println("\t go run main.go command [arguments]")
	fmt.Println()
	fmt.Println("\t generategensis \tcreate a genesis block and save to boltdb file. \n\t\t  -genesis \tcan save  your data with custom")
	fmt.Println("\t createblock \tcreate a new block and save to boltdb file. \n\t\t  -data \tcan save your custom data")
	fmt.Println("\t getlastblock \treturn a last block on boltdb file. \n")
	fmt.Println("Use go run main.go help [topic] for more information about that topic.")

}

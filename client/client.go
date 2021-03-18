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
	"2021/_03_公链/XianFengChain04/transaction"
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
	if len(os.Args) == 1 { //无参数、直接调帮助文档
		cmd.Help()
	} else {
		switch os.Args[1] { //有参数，调用参数对应的方法
		case GENERATEGENESIS:
			if len(os.Args[2:]) > 2 {
				fmt.Println("only  -genesis")
				return
			}
			var address string
			flagSet := flag.NewFlagSet(GENERATEGENESIS, flag.ExitOnError)
			flagSet.StringVar(&address, "address", "", "用户指定的地址")
			flagSet.Parse(os.Args[2:])

			flag := new(big.Int).SetBytes(cmd.Chain.LastBlock.Hash[:]).Cmp(big.NewInt(0)) == 1
			if flag {
				fmt.Println("创世区块已存在！")
			} else {
				//1、生成一个铸币交易
				coinbase, err := transaction.CraeteCoinbase(address)
				if err != nil {
					fmt.Println("抱歉，创建coinbase交易遇到错误，请重试！")
					return
				}
				err = cmd.Chain.CreateChainWithGenesis([]transaction.Transaction{*coinbase})
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("generate genesis block success ! genesis block's data is", address)
				}
			}

		case SENDTRANSACTION:
			if len(os.Args[2:]) > 6 {
				fmt.Println("参数无法识别：", os.Args[2:])
				return
			}
			var from string
			var to string
			var amount float64
			flagSet := flag.NewFlagSet(SENDTRANSACTION, flag.ExitOnError)
			flagSet.StringVar(&from, "from", "", "发起者地址")
			flagSet.StringVar(&to, "to", "", "接受者地址")
			flagSet.Float64Var(&amount, "amount", 0, "转账的金额")
			flagSet.Parse(os.Args[2:])

			//flag 为true时，说明 已经有创世区块了。
			flag := new(big.Int).SetBytes(cmd.Chain.LastBlock.Hash[:]).Cmp(big.NewInt(0)) == 1
			if !flag {
				fmt.Println("Error: ,your client don't have a blockchain now !")
				fmt.Println("Use generategenesis [-genesis 'data'] to create a new blockchain ")
			}
			err := cmd.Chain.SendTransaction(from, to, amount)
			if err !=nil {
				fmt.Println("抱歉，发送交易时出错",err.Error())
				return
			}
			fmt.Println("恭喜，交易发送成功！")
		case GETLASTBLOCK:
			if len(os.Args[1:]) > 1 {
				fmt.Println("参数无法识别：", os.Args[1:])
				return
			}
			IsChain := new(big.Int).SetBytes(cmd.Chain.LastBlock.Hash[:]).Cmp(big.NewInt(0)) == 1
			if !IsChain {
				fmt.Println("Error: ,your client don't have a blockchain now !")
				fmt.Println("Use generategenesis [-genesis 'data'] to create a new blockchain ")
				return
			}
			fmt.Println("调用获取最新区块功能！")
			flag.String(GETLASTBLOCK, "", "用于获取最新区块的数据")
			fmt.Printf("%+v\n", cmd.Chain.LastBlock)

		case GETALLBLOCKS:
			fmt.Println("调所有区块的功能！")
		case HELP:
			fmt.Println("调用帮助说明功能！")
		case VERSION:
			fmt.Println("调用获取版本信息！")
		default:
			fmt.Println("go run main.go: Unknown subcommand.")
			fmt.Println("Run 'go run main.go -help' for usage")
		}
	}
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

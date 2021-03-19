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
	"2021/_03_公链/XianFengChain04/utils"
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
			cmd.GenerateGenesis()
		case SENDTRANSACTION:
			cmd.SendTransaction()
		case GETBALANCE:
			cmd.GetBalance()
		case GETLASTBLOCK:
			cmd.GetLastBlock()
		case GETALLBLOCKS: //查询所有区块
			blocks, err := cmd.Chain.GetAllBlocks()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			for _, block := range blocks {
				for _, tx := range block.Transactions {
					for index, input := range tx.Inputs {
						fmt.Printf("第%d 个区块中的 第 %v 笔交易输入姓名：%s\n", block.Height+1, index+1, input.ScriptSig)
					}
					for index, output := range tx.OutPuts {
						fmt.Printf("第%d 个区块中的 第 %v 笔交易输出金额：%+v , 交易输出对象为：%s\n", block.Height+1, index+1, output.Value, output.ScriptPub)
					}

				}
			}
		case HELP:
			cmd.Help()
		case VERSION:
			fmt.Println("调用获取版本信息！")
		default:
			fmt.Println("go run main.go: Unknown subcommand.")
			fmt.Println("Run 'go run main.go -help' for usage")
		}
	}
}
func (cmd *CmdClient) GetLastBlock() {
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

}

//发送交易功能
func (cmd *CmdClient) SendTransaction() {
	if len(os.Args[2:]) > 6 {
		fmt.Println("参数无法识别：", os.Args[2:])
		return
	}
	var from string
	var to string
	var amount string
	flagSet := flag.NewFlagSet(SENDTRANSACTION, flag.ExitOnError)
	flagSet.StringVar(&from, "from", "", "发起者地址")
	flagSet.StringVar(&to, "to", "", "接受者地址")
	flagSet.StringVar(&amount, "amount", "", "转账的金额")
	flagSet.Parse(os.Args[2:])
	///解析多个参数数组
	fromSlice, err := utils.JSONArrayToString(from)
	if err != nil {
		fmt.Println(from)
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}
	toSlice, err := utils.JSONArrayToString(to)
	if err != nil {
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}
	amountSlice, err := utils.JSONArrayToFloat(amount)
	if err != nil {
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}
	//判断各参数列表长度是否一致
	if len(fromSlice) != len(toSlice) || len(toSlice) != len(amountSlice) {
		fmt.Println("各参数长度不一致")
	}
	//flag 为true时，说明 已经有创世区块了。
	flag := new(big.Int).SetBytes(cmd.Chain.LastBlock.Hash[:]).Cmp(big.NewInt(0)) == 1
	if !flag {
		fmt.Println("Error: ,your client don't have a blockchain now !")
		fmt.Println("Use generategenesis [-genesis 'data'] to create a new blockchain ")
	}

	err = cmd.Chain.SendTransaction(fromSlice,toSlice ,amountSlice )
	if err != nil {
		fmt.Println("抱歉，发送交易时出错", err.Error())
		return
	}
	fmt.Println("恭喜，交易发送成功！")
}

//创建创世区块功能
func (cmd *CmdClient) GenerateGenesis() {
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
		err := cmd.Chain.CreateCoinBase(address)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("generate genesis block success ! genesis block's reward 50 BTC to ", address)
		}
	}

}

//Help 命令方法
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
func (cmd *CmdClient) GetBalance() {
	getbalance := flag.NewFlagSet(GETBALANCE, flag.ExitOnError)
	var addr string
	getbalance.StringVar(&addr, "address", "", "用户的地址")
	getbalance.Parse(os.Args[2:])
	//先判断是否有创世区块
	hasGenesis := new(big.Int).SetBytes(cmd.Chain.LastBlock.Hash[:]).Cmp(big.NewInt(0)) == 0
	if hasGenesis {
		fmt.Println("该网络不存在，无法查询")
		return
	}
	//有这条链，才调用查询功能
	balance := cmd.Chain.GetBalance(addr)
	fmt.Printf("地址 [%s] 的余额是：%f\n", addr, balance)
}
/*
go run main.go sendtransaction -from [\"yugu\",\"xiaobing\"] -to [\"xiaobing\",\"shipeng\"] -amount [10,10]
go run main.go sendtransaction -from "[\"yugu\",\"xiaobing\"]" -to "[\"xiaobing\",\"shipeng\"]" -amount "[10,10]"

 */
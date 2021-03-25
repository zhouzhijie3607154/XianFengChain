package client

const (
	GENERATEGENSIS  = "generategensis"  //ccoinbase -addr
	SENDTRANSACTION = "sendtransaction" //sendTransaction from to amount
	GETBALANCE      = "getbalance"      //获取地址的余额功能
	GETLASTBLOCK    = "getlastblock"    //获取最新区块
	GETALLBLOCKS    = "getallblocks"    //获取所有区块
	GETNEWADDRESS   = "getnewaddress"   //生成新的比特币地址
	GETALLADDRESS   = "getalladdress"    //获取所有的已生成的地址
	DUMPPRIKEY		= "dumpprivkey"		//导出指定地址的私钥
	HELP            = "help"
)

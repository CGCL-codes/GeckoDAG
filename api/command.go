package api

import (
	"Jight/config"
	jrpc "Jight/rpc"
	"Jight/utils"
	"Jight/wallet"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

type CMD struct {
	rpcPort int
}


func (cmd *CMD) Run() {
	// TODO nodeID should be mac or ip address, nodeID specialize a node
	nodeID := "self"

	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name: "rpcport, r",
			Usage: "Daemon RPC port to connect to",
			Destination: &cmd.rpcPort,
		},
	}

	app.Commands = []cli.Command{
		// createdagchain CMD
		{
			Name:        "createdagchain",
			Aliases:     []string{"cdc"},
			Usage:       "Create a dagchain and send genesis tx reward to ADDRESS",
			UsageText:   "createdagchain -address ADDRESS",
			Description: "Create a dagchain",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "address, a"},
			},
			Action: func(c *cli.Context) error {
				address := c.String("address")
				cmd.createDagChain(address)
				return nil
			},
		},


		// getbalance CMD
		{
			Name:        "getbalance",
			Aliases:     []string{"gb"},
			Usage:       "Get balance of ADDRESS",
			UsageText:   "getbalance -address ADDRESS",
			Description: "Get balance of ADDRESS",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "address, a"},
			},
			Action: func(c *cli.Context) error {
				address := c.String("address")
				cmd.getBalance(address)
				return nil
			},
		},

		// list all address of wallet CMD
		{
			Name:        "listaddress",
			Aliases:     []string{"la"},
			Usage:       "Lists all addresses from the wallet file",
			UsageText:   "listaddress",
			Description: "Lists all addresses from the wallet file",
			Action: func(c *cli.Context) error {
				cmd.listAddress(nodeID)
				return nil
			},
		},

		// send money from f to t CMD
		{
			Name:        "send",
			Aliases:     []string{"sd"},
			Usage:       "Send AMOUNT of coins from FROM address to TO",
			UsageText:   "send -from FROM -to TO -amount AMOUNT",
			Description: "Send AMOUNT of coins from FROM address to TO",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "from, f"},
				cli.StringFlag{Name: "to, t"},
				cli.UintFlag{Name: "amount, a"},
				cli.StringFlag{Name: "income, i", Value: ""},
			},
			Action: func(c *cli.Context) error {
				from := c.Int("from")
				to := c.Int("to")
				amount := c.Int("amount")
				income := c.String("income")
				cmd.send(from, amount, to, income)
				return nil
			},
		},

		{
			Name: "initializewallets",
			Aliases: []string{"iw"},
			Usage: "Initialize the wallets for test according to the 1000 genesis addresses",
			Action: func(c *cli.Context) error {
				if err := cmd.initializeWallets(nodeID); err != nil {
					return  err
				}
				return nil
			},
		},

		{
			Name: "createbatchtxs",
			Aliases: []string{"cbt"},
			Usage: "Create batch transactions to test confirmation latency",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "times, t", Usage: "how many times to create batch transactions"},
				cli.IntFlag{Name: "speed, s", Usage: "speed to issue transactions by all the system, s transactions/min"},
			},
			Action: func(c *cli.Context) error {
				times := c.Int("times")
				speed := c.Int("speed")
				if err:= cmd.createBatchTxs(times, speed); err!=nil {
					return err
				}
				return nil
			},

		},
		{
			Name: "refreshtips",
			Aliases: []string{"rt"},
			Usage: "Clean the old tips",
			Action: func(c *cli.Context) error {
				if err:= cmd.refreshTips(); err!=nil {
					return err
				}
				return nil
			},

		},
		{
			Name: "fetchtxssentbyaccount",
			Aliases: []string{"ftsa"},
			Usage: "Fetch transactions sent by an account",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "account, t", Usage: "specify the account"},
			},
			Action: func(c *cli.Context) error {
				accountId :=  c.Int("account")
				if err:= cmd.fetchTxsSentByAccount(accountId); err!=nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "fetchtussentbyaccount",
			Aliases: []string{"ftusa"},
			Usage: "Fetch transaction unions sent by an account",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "account, t", Usage: "specify the account"},
			},
			Action: func(c *cli.Context) error {
				accountId :=  c.Int("account")
				if err:= cmd.fetchTUSentbyaccount(accountId); err!=nil {
					return err
				}
				return nil
			},

		},
		{
			Name: "fetchtx",
			Aliases: []string{"ftx"},
			Usage: "Fetch a transaction from DB",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "txId, t", Usage: "specify the transaction"},
			},
			Action: func(c *cli.Context) error {
				txId :=  c.Int("txId")
				if err:= cmd.fetchTx(txId); err!=nil {
					return err
				}
				return nil
			},

		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}

}

// TODO this func should implemented in daemon
func (cmd *CMD) createDagChain(address string) {
	fmt.Println("create DAGCHAIN: ", address)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	cdcc := jrpc.CreateDagChainCMD{
		Address: "tJight123344566765765756",
	}
	var cdcr jrpc.CreateDagChainReply

	err = client.Call("Jightd.CreateDagChain", cdcc, &cdcr)
	if err != nil {
		log.Fatal("CreateDagChain error: ", err)
	} else {
		fmt.Println("CreateDagChain CMD: ", cdcr.CdcrField)
	}

}

func (cmd *CMD) initializeWallets(nodeID string) error {
	if utils.CheckFileExisted(config.WALLET_FILE) {
		return errors.New("wallet file is already existed")
	}
	for _, addr := range(config.GenesisAddresses) {
		cmd.createWallet(nodeID, addr)
	}
	return nil
}


func (cmd *CMD) createWallet(nodeID string, address string) {
	fmt.Println("create wallet")

	// write the new address to local wallet_self.dat
	wallets, _ := wallet.LoadWallets(nodeID)
	wallets.CreateWallet(address)
	wallets.SaveToFile(nodeID)

	// write the new address to DB via rpc
	log.Println("cmd.rpcPort: ", cmd.rpcPort)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	cwc := jrpc.CreateWalletCMD{
		Address: address,
	}

	var cwr jrpc.CreateWalletReply

	err = client.Call("Jightd.CreateWallet", cwc, &cwr)
	if err != nil {
		log.Fatal("CreateWallet error: ", err)
	}

	fmt.Printf("Generated address: %s\n", address)

}

func (cmd *CMD) getBalance(address string) {
	fmt.Println("get balance of ", address)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	gbc := jrpc.GetBalanceCMD{
		Address: address,
	}

	var gbr jrpc.GetBalanceReply

	err = client.Call("Jightd.GetBalance", gbc, &gbr)
	if err != nil {
		log.Fatal("GetBalance error: ", err)
	} else {
		fmt.Printf("Balance of %s: %s\n",address, gbr.Balance)
	}
}

/*
// Merge conflict
func (cmd *CMD) listAddress() {
	fmt.Println("list address")
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	lac := jrpc.ListAddressCMD{
	}

	var lar jrpc.ListAddressReply

	err = client.Call("Jightd.ListAddress", lac, &lar)
	if err != nil {
		log.Fatal("ListAddress error: ", err)
	} else {
		fmt.Println("ListAddress CMD: ", lar.Addresses)
*/

func (cmd *CMD) listAddress(nodeID string) {
	wallets, err := wallet.LoadWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cmd *CMD) send(from int, amount int, to int, income string) {
	//fmt.Println("send from ", from, "to ", to, "amount ", amount)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	sc := jrpc.SendCmd{
		From: from,
		To: to,
		IncomeTx: income,
		Amount: amount,
	}

	var sr jrpc.SendReply

	err = client.Call("Jightd.Send", sc, &sr)
	if err != nil {
		log.Fatal("Send error: ", err)
	} else {
		//fmt.Printf("Transaction hash: %s\n", sr.Tx)
	}
}

/**
@param times: the times to repeat
*/
func (cmd * CMD) createBatchTxs(times int, speed int) error {
	fmt.Println("create batch transactins ", times, "times, with speed: ", speed)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	cbtc := jrpc.CreateBatchTxsCMD{
		Times:times,
		Speed:speed,
	}

	var cbtr jrpc.CreateBatchTxsReply

	err = client.Call("Jightd.CreateBatchTxs", cbtc, &cbtr)
	if err != nil {
		log.Fatal("Send error: ", err)
		return err
	} else {
		fmt.Println("Successful createBatchTxs")
		return nil
	}
}

func (cmd *CMD) refreshTips() error {
	//fmt.Println("Refresh tips")
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	rtc := jrpc.RefreshTipsCMD{
	}

	var rtr jrpc.RefreshTipsReply

	err = client.Call("Jightd.RefreshTips", rtc, &rtr)
	if err != nil {
		log.Fatal("Send error: ", err)
		return err
	} else {
		//fmt.Printf("Count of tips before refresh: %d\n", rtr.TipsCntBefore)
		//fmt.Printf("Count of tips after refresh: %d\n", rtr.TipsCntAfter)
		return nil
	}
}

func (cmd *CMD) fetchTxsSentByAccount(accountId int) error {
	fmt.Printf("Fetch transactions sent by account: %d\n", accountId)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	ftac := jrpc.FetchTxsByAccountCMD{
		AccountId: accountId,
	}

	var ftar jrpc.FetchTxsByAccountReply

	err = client.Call("Jightd.FetchTxsSentByAccount", ftac, &ftar)
	if err != nil {
		log.Fatal("Send error: ", err)
		return err
	} else {
		log.Printf("Txs: %v\n", ftar.TxsNumber)
		log.Printf("Txs count: %d\n", len(ftar.TxsNumber))
		return nil
	}
}


func (cmd *CMD) fetchTUSentbyaccount(accountId int) error {
	fmt.Printf("Fetch transaction unions sent by account: %d\n", accountId)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	ftac := jrpc.FetchTUsByAccountCMD{
		AccountId: accountId,
	}

	var ftar jrpc.FetchTUsByAccountReply

	err = client.Call("Jightd.FetchTUsSentByAccount", ftac, &ftar)
	if err != nil {
		log.Fatal("Send error: ", err)
		return err
	} else {
		log.Printf("TUs: %v\n", ftar.TUsNumber)
		log.Printf("TUs count: %d\n", len(ftar.TUsNumber))
		return nil
	}
}


func (cmd *CMD) fetchTx(txId int) error {
	fmt.Printf("Fetch transaction from DB: %d\n", txId)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	ftc := jrpc.FetchTxCMD{
		TxId: txId,
	}

	var ftr jrpc.FetchTxReply

	err = client.Call("Jightd.FetchTx", ftc, &ftr)
	if err != nil {
		log.Fatal("Send error: ", err)
		return err
	} else {
		if ftr.Existed {
			log.Printf("Tx: %v exists!\n", ftc.TxId)
		} else {
			log.Printf("Tx: %v does not exist!\n", ftc.TxId)
		}
		return nil
	}
}

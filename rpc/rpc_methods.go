package rpc

import (
	"Jight/config"
	"Jight/dagchain"
	"Jight/p2p"
	"Jight/wallet"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

func (jd *Jightd) GetBalance (cmd GetBalanceCMD, reply *GetBalanceReply) error {
	account, err := dagchain.FetchAccount(jd.DC.DBOthers, cmd.Address)
	if err != nil {
		return errors.New("no match address: " + cmd.Address)
	}

	balance := account.Balance

	reply.Balance = strconv.FormatUint(uint64(balance), 10)

	return nil
}

/**
for test confirmation latency
the transaction will decrease the balance of sender directly
and the transaction will increase the balance of receiver directly
*/
func (jd *Jightd) Send (cmd SendCmd, reply *SendReply) error {
	// get the income value if existed
	var income int
	var incomeTxByteSlice [32]byte
	if cmd.IncomeTx != "" {
		fmt.Println("cmd.IncomeTx is not nil")
		incomeTxByteSlice, err := hex.DecodeString(cmd.IncomeTx)
		if err !=nil {
			return err
		}
		tx, err := dagchain.FetchTx(jd.DC.DB, incomeTxByteSlice)
		if err!=nil {
			return err
		}
		income = tx.Value
	}

	// check if balance is enough
	fmt.Println("cmd.From: "+dagchain.AccountInt2Addr[cmd.From])
	fromAccount, err := dagchain.FetchAccount(jd.DC.DBOthers, dagchain.AccountInt2Addr[cmd.From])
	if err !=nil {
		return err
	}
	if fromAccount.Balance + income - cmd.Amount < 0 {
		return errors.New("balance is not enough for this transaction")
	} else {
		// substract the balance of fromAccount
		fromAccount.Balance = fromAccount.Balance + income - cmd.Amount
		dagchain.StoreAccount(jd.DC.DBOthers, *fromAccount)
	}

	// add the balance of toAccount
	toAccount, err := dagchain.FetchAccount(jd.DC.DBOthers, dagchain.AccountInt2Addr[cmd.To])
	if err!=nil {
		return err
	}
	toAccount.Balance = toAccount.Balance + cmd.Amount
	dagchain.StoreAccount(jd.DC.DBOthers, *toAccount)

	// TODO nodeID should be mac or ip address, nodeID specialize a node
	nodeID := "self"
	ws, err := wallet.LoadWallets(nodeID)
	if err != nil {
		return err
	}

	wallet := ws.GetWallet(dagchain.AccountInt2Addr[cmd.From])


	gt, lastTC, isTU := jd.DC.CreateGT(wallet.PrivateKey, incomeTxByteSlice, dagchain.AccountInt2Addr[cmd.From], dagchain.AccountInt2Addr[cmd.To],
		cmd.From, cmd.To, cmd.Amount)

	for i, gx := range dagchain.GXs {
		fmt.Printf("The %dth gx is cited by %d times\n", i, gx.FetchCitedCount())
	}

	reply.Tx = hex.EncodeToString([]byte("justreply"))

	jd.DC.AddGT(gt, lastTC, isTU, true)
	p2p.SyncGT(gt)
	dagchain.GXs[gt.FetchNumber()] = gt

	// if it is a transaction union, prune the old transactions
	if isTU {
		senderAccount := dagchain.AccountMap[cmd.From]
		dagchain.PruneOldTxs(jd.DC, senderAccount)
	}
	//fmt.Println()
	//fmt.Println(dagchain.PrintGXs())
	//fmt.Println(dagchain.PrintAccounts())

	return nil
}

/**
the transaction will decrease the balance of sender directly
however, the balance of receiver will not increase directly
the balance of receiver will increase until it is cited by a later transaction via 'income reference'
*/
/*func (jd *Jightd) Send (cmd SendCmd, reply *SendReply) error {
	// get the income value if existed
	var income uint32
	var incomeTxByteSlice []byte
	if cmd.IncomeTx != "" {
		fmt.Println("cmd.IncomeTx is not nil")
		incomeTxByteSlice, err := hex.DecodeString(cmd.IncomeTx)
		if err !=nil {
			return err
		}
		tx, err := dagchain.FetchTx(jd.DC.DB, incomeTxByteSlice)
		if err!=nil {
			return err
		}
		income = tx.Value
	}

	// check if balance is enough
	fromAccount, err := dagchain.FetchAccount(jd.DC.DB, []byte(cmd.From))
	if err !=nil {
		return err
	}
	if fromAccount.Balance + income - cmd.Amount < 0 {
		return errors.New("balance is not enough for this transaction")
	} else {
		// substract the balance of fromAccount
		fromAccount.Balance = fromAccount.Balance + income - cmd.Amount
		dagchain.StoreAccount(jd.DC.DB, *fromAccount)
	}

	// TODO nodeID should be mac or ip address, nodeID specialize a node
	nodeID := "self"
	ws, err := wallet.LoadWallets(nodeID)
	if err != nil {
		return err
	}

	wallet := ws.GetWallet(cmd.From)

	tx := jd.DC.CreateTx(wallet.PrivateKey, incomeTxByteSlice, []byte(cmd.From), []byte(cmd.To), cmd.Amount)

	reply.Tx = hex.EncodeToString(tx.Hash)
	log.Println("reply.Tx: ", reply.Tx)

	jd.DC.AddTx(tx, true)
	p2p.SyncTx(tx)

	return nil
}*/

func (jd *Jightd) CreateWallet (cmd CreateWalletCMD, reply *CreateWalletReply) error {
	address := cmd.Address
	// lastid of new account is nil
	acc := dagchain.CreateNewAccount(address, 0)
	log.Println("func CreateWallet is called")
	return dagchain.StoreAccount(jd.DC.DBOthers, *acc)
}


func (jd *Jightd) CreateBatchTxs (cmd CreateBatchTxsCMD, reply *CreateBatchTxsReply) error {
	for t := 0; t< cmd.Times; t++ {
		for i:=0; i<cmd.Speed; i++ {
			r1 := rand.Intn(int(config.GENESIS_ADDR_COUNT))
			r2 := rand.Intn(int(config.GENESIS_ADDR_COUNT))
			fromAccount := config.GenesisAddresses[r1]
			toAccount := config.GenesisAddresses[r2]

			// TODO nodeID should be mac or ip address, nodeID specialize a node
			nodeID := "self"
			ws, err := wallet.LoadWallets(nodeID)
			if err != nil {
				return err
			}

			wallet := ws.GetWallet(fromAccount)

			//dagchain.AccountMap[r1].TxCount++
			/*var gt dagchain.GeneralTx
			var tc *dagchain.TxContent
			if dagchain.AccountMap[r1].TxCount % config.MERGE_PERIOD != 0 {
				gt = jd.DC.CreateTx(wallet.PrivateKey, [32]byte{}, fromAccount, toAccount, r1, r2, 0)
			} else {
				gt, tc = jd.DC.CreateTU(wallet.PrivateKey, [32]byte{}, fromAccount, toAccount, r1, r2, 0)
			}*/

			gt, lastTC, isTU := jd.DC.CreateGT(wallet.PrivateKey, [32]byte{}, fromAccount, toAccount,
				r1, r2, 0)



			jd.DC.AddGT(gt, lastTC, isTU,true)
			p2p.SyncGT(gt)
			dagchain.GXs[gt.FetchNumber()] = gt

			// if it is a transaction union, prune the old transactions
			if isTU {
				senderAccount := dagchain.AccountMap[r1]
				dagchain.PruneOldTxs(jd.DC, senderAccount)
			}
		}
		//fmt.Println("Txs: ")
		//for _, tx := range transaction.Txs {
		//	fmt.Printf("tx: %d, ver1: %d, ver2: %d\n", tx.Number, tx.ValidateNum[0], tx.ValidateNum[1])
		//}

		// cited tips will be clean every NETWORK_LATENCY seconds
		if t % config.NETWORK_LATENCY == 0 {
			log.Printf("Create batch transactions %d th times\n", t)

			var uncitedTips = make(map[int]*dagchain.Tip)
			//fmt.Println("Original tips: ")
			for _, t := range jd.DC.Tips {
				//fmt.Printf("tip: %d, ver1: %d, ver2: %d, cited: %t, verified: %t\n", t.TxNum, t.VerifyNum[0],
				//	t.VerifyNum[1], t.Cited, t.Verified)
				if !t.Cited {
					uncitedTips[t.TxNum] = t
				}
			}
			jd.DC.Tips = uncitedTips

			verifiedTipsCount := 0
			//fmt.Println("Tips after clean: ")
			for _, t := range jd.DC.Tips {
				//fmt.Printf("tip: %d, ver1: %d, ver2: %d, cited: %t, verified: %t\n", t.TxNum, t.VerifyNum[0],
				//	t.VerifyNum[1], t.Cited, t.Verified)
				if t.Verified {
					verifiedTipsCount ++
				}
			}

			fmt.Printf("Total tips count: %d, verified count: %d\n", len(jd.DC.Tips), verifiedTipsCount)
		}

	}
	return nil
}

func (jd *Jightd) RefreshTips (cmd RefreshTipsCMD, reply *RefreshTipsReply) error {
	reply.TipsCntBefore = len(jd.DC.Tips)
	var uncitedTips = make(map[int]*dagchain.Tip)
	//fmt.Println("Original tips: ")
	for _, t := range jd.DC.Tips {
		//fmt.Printf("tip: %d, ver1: %d, ver2: %d, cited: %t, verified: %t\n", t.TxNum, t.VerifyNum[0],
		//	t.VerifyNum[1], t.Cited, t.Verified)
		if !t.Cited {
			//uncitedTips = append(uncitedTips, t)
			uncitedTips[t.TxNum] = t
		}
	}
	jd.DC.Tips = uncitedTips
	reply.TipsCntAfter = len(jd.DC.Tips)
	return nil
}

func (jd *Jightd) FetchTxsSentByAccount (cmd FetchTxsByAccountCMD, reply *FetchTxsByAccountReply) error {
	accountID := cmd.AccountId
	reply.TxsNumber = dagchain.AccountMap[accountID].WithoutPruneIds
	return nil
}


func (jd *Jightd) FetchTUsSentByAccount (cmd FetchTUsByAccountCMD, reply *FetchTUsByAccountReply) error {
	log.Println("In function FetchTUsSentByAccount")
	accountID := cmd.AccountId
	txsNumber := dagchain.AccountMap[accountID].WithoutPruneIds
	var tusNumber []int
	for _, tx := range txsNumber {
		if tu, ok := dagchain.GXs[tx].(*dagchain.TU); ok {
			tusNumber = append(tusNumber, tx)
			log.Printf("Lenght of transaction content: %d\n", len(tu.TCList))
		}
	}
	reply.TUsNumber = tusNumber
	return nil
}

func (jd *Jightd) FetchTx (cmd FetchTxCMD, reply *FetchTxReply) error {
	log.Println("In function FetchTx")
	existed := dagchain.FetchTx2(jd.DC.DBMerging, cmd.TxId)
	if existed {
		reply.Existed = true
	} else {
		reply.Existed = false
	}
	return nil
}
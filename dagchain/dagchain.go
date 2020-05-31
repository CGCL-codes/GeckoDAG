package dagchain

import (
	"Jight/config"
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"math/rand"
	"time"
)

/*
 Different nodes may see different views of tips due to the network latency
 However, since we cannot deploy the nodes all across the world currently, we have to simulate the different views of tips
 Specifically,
	1) if a node issues a tx citing two tips by itself, it will delete the tips directly;
	2) if a node receives a tx from other nodes which cites two tips, only the 'cited' flags in these two tips are set
	3) tips with 'cited' flags being set will be deleted later, either periodically or triggered by the cli
*/
type Tip struct {
	TxNum	int
	TipID [32]byte	// tip txid
	Cited bool	// denote if it has been cited by the txs from other nodes
	Verified bool // indicate if the tip verify the sample transaction
	//Account []byte	// Account address of tip
}

type DagChain struct {
	Tips map[int]*Tip // tips of dagchain, persistence
	//Account map[string][]byte		// no persistence, reindex at CreateDagChain or RecoverDagChain, to reduce storage costs and speed up FindParent
	DB *leveldb.DB	// no merging and pruning
	DBMerging *leveldb.DB  // with merging and pruning
	DBOthers *leveldb.DB // DB used to store wallet and initialize transactions
}

/*
	--------------
	  tx related
	--------------
*/
// receive tx from p2p network, then store it to db
// local; denote if the tx is created locally or received from other nodes
func (dc *DagChain) AddGT(gt GeneralTx, tc *TxContent, isTU bool, local bool) error {
	// verify tx from network
	if gt.Verify() == false {
		return errors.New("Tx verify failed, invalid Tx")
	}

	oldTips:= gt.FetchLatestValidateNum()
	//var citedChangeTips []*Tip
	for _, t := range dc.Tips {
		if t.TxNum == oldTips[0] || t.TxNum == oldTips[1] {
			t.Cited = true
			//citedChangeTips = append(citedChangeTips, t)
		}
	}

	newTip := &Tip{gt.FetchNumber(), gt.FetchHash(), false, gt.CheckVerification()}

	//dc.Tips = append(dc.Tips, newTip)
	dc.Tips[newTip.TxNum] = newTip

	// write gt to DB
	// 1.1 write gt to db with no merging and pruning
	if !isTU {
		if err := StoreTx2(dc.DB, gt); err != nil {
			log.Panic("store tx to db failed: ", err)
			return err
		}
	} else {
		// if it is a TU, only the single transaction is needed to be stored
		// for test: a mock transaction is used
		var mockTx GeneralTx =  &Transaction{gt.FetchNumber(), [2]int{1, 2}, config.MOCK_TX, config.MOCK_TX,
			[2][32]byte{config.MOCK_TX, config.MOCK_TX}, config.MOCK_TX,[34]byte{}, 0, 1000,
			[34]byte{},0, time.Now().Unix(), [64]byte{}, false, 0}

		if err := StoreTx2(dc.DB, mockTx); err != nil {
			log.Panic("store tx to db failed: ", err)
			return err
		}
	}
	// 1.2 write gt to db with merging and pruning
	if !isTU {
		if err := StoreTx2(dc.DBMerging, gt); err != nil {
			log.Panic("store tx to dbmerging failed: ", err)
			return err
		}
	} else {
		if err := StoreTx2(dc.DBMerging, gt); err != nil {
			log.Panic("store tx to dbmerging failed: ", err)
			return err
		}
		if err := RefreshTCL2(dc.DBMerging, tc); err!=nil {
			log.Panic("refresh tcl to dbmerging failed: ", err)
			return err
		}
	}

	/*// write to DB
	err := StoreTx(dc.DB, *tx)
	if err != nil {
		log.Panic("store tx failed: ", err)
	}
	for _, v := range citedChangeTips {
		PutTip(dc.DB, *v)
	}
	PutTip(dc.DB, *newTip)*/

	/*if local {
		fmt.Println("Receive tx locally: ", tx)

		// remove old tips and add new tip
		var index []int
		oldTips := tx.Validate
		for i, t := range dc.Tips {
			if bytes.Equal(t.TipID, oldTips[0]) || bytes.Equal(t.TipID, oldTips[1]) {
				index = append(index, i)
			}
		}

		newTip := &Tip{tx.Number, tx.ValidateNum, tx.Hash, false, tx.Verification}
		log.Println("length of oldTips: ", len(index))
		if len(index) == 2 {
			dc.Tips[index[0]] = newTip
			dc.Tips = append(dc.Tips[:index[1]], dc.Tips[index[1]:]...)
		} else if len(index) == 1 {
			dc.Tips[index[0]] = newTip
		} else {
			dc.Tips = append(dc.Tips, newTip)
		}
		// write to DB
		err := StoreTx(dc.DB, *tx)
		if err != nil {
			log.Panic("store tx failed: ", err)
		}
		for _, v := range oldTips {
			RemoveTip(dc.DB, v)
		}
		PutTip(dc.DB, *newTip)
	} else {
		fmt.Println("Receive tx remotely: ", tx)
		// set the cited flag and add new tip
		//oldTips := tx.Validate
		var citedChangeTips []*Tip
		//for _, t := range dc.Tips {
		//	if bytes.Equal(t.TipID, oldTips[0]) || bytes.Equal(t.TipID, oldTips[1]) {
		//		t.cited = true
		//		citedChangeTips = append(citedChangeTips, t)
		//	}
		//}
		newTip := &Tip{tx.Hash, false, tx.Verification}
		dc.Tips = append(dc.Tips, newTip)

		// write to DB
		err := StoreTx(dc.DB, *tx)
		if err != nil {
			log.Panic("store tx failed: ", err)
		}
		for _, v := range citedChangeTips {
			PutTip(dc.DB, *v)
		}
		PutTip(dc.DB, *newTip)
	}*/

	//fmt.Printf("Tips count: %d\n", len(dc.Tips))

	return nil
}

// create a general transaction
// the return bool value denotes if the general transaction is a transaction union
func (dc *DagChain) CreateGT(privKey ecdsa.PrivateKey, income [32]byte, sender, receiver string, senderNo, receiverNo int, value int) (GeneralTx, *TxContent, bool) {
	senderAccount := AccountMap[senderNo]
	senderAccount.TxCount++
	if senderAccount.TxCount % config.MERGE_PERIOD == 0 {
		gt, tc := dc.CreateTU(privKey, income, sender, receiver, senderNo, receiverNo, value)
		return gt, tc, true
	} else {
		gt := dc.CreateTx(privKey, income, sender, receiver, senderNo, receiverNo, value)
		return gt, nil, false
	}
}

// create (mining) a tx
func (dc *DagChain) CreateTx(privKey ecdsa.PrivateKey, income [32]byte, sender, receiver string, senderNo, receiverNo int, value int) *Transaction {
	var validateRef [2][32]byte
	verification := false

	//index, vTips := dc.SelectTips()
	_, vTips := dc.SelectTips(senderNo)
	for dc.VerifyTips(vTips) == false {
		_, vTips = dc.SelectTips(senderNo)
	}
	//log.Println("length of selectedTips: ", len(vTips))
	lastTx := dc.FindParent(sender)

	for i, v := range vTips {
		validateRef[i] = v.TipID
		v.Cited = true
		if v.Verified {
			verification = true
		}
	}

	tx := NewTx([2]int{vTips[0].TxNum, vTips[1].TxNum}, lastTx, validateRef, income, sender, senderNo, value, receiver, verification)

	senderAccount := AccountMap[senderNo]
	senderAccount.LastIdNo = tx.Number
	senderAccount.WithoutMergeIds[(senderAccount.TxCount-1)%config.MERGE_PERIOD]=tx.Number
	senderAccount.WithoutPruneIds = append(senderAccount.WithoutPruneIds, tx.Number)

	// sign a tx
	tx.Sign(privKey)
	tx.Nonce = tx.Pow()
	tx.Hash = tx.HashTx()

	// add the citedcount
	GXs[vTips[0].TxNum].AddCitedCount()
	GXs[vTips[1].TxNum].AddCitedCount()

	return tx
}

// create a transaction union
func (dc *DagChain) CreateTU(privKey ecdsa.PrivateKey, income [32]byte, sender, receiver string, senderNo, receiverNo int, value int) (*TU, *TxContent) {
	senderAccount := AccountMap[senderNo]
	tipsNo, vTips := dc.SelectTips(senderNo)
	for dc.VerifyTips(vTips) == false {
		tipsNo, vTips = dc.SelectTips(senderNo)
	}
	//log.Println("length of selectedTips: ", len(vTips))
	//lastTx := dc.FindParent(sender)

	log.Println("tipsNo before a TU:", tipsNo)

	if senderAccount.LatestTU!=nil {
		log.Println("Before a new TU, tu.ValidateNum: ", senderAccount.LatestTU.ValidateNum)
	}
	tu, lastTC := NewTU(tipsNo, senderAccount, income, value, receiver)
	log.Println("After a new TU, tu.ValidateNum: ", tu.ValidateNum)
	senderAccount.LatestTU = tu
	senderAccount.WithoutPruneIds = append(senderAccount.WithoutPruneIds, tu.Number)

	tu.Signature = [64]byte{}

	// add the citedcount
	/*GXs[tipsNo[0]].AddCitedCount()
	GXs[tipsNo[1]].AddCitedCount()*/

	return tu, lastTC
}


// Signature a tx with privKey
func (dc *DagChain) SignTx(tx *Transaction, privKey ecdsa.PrivateKey)  {
	if tx.IsGenesisTx() {
		return
	}

	tx.Sign(privKey)
}

// Verify if a tx is valid
func (dc *DagChain) VerifyTx(tx *Transaction) bool {
	if tx.IsGenesisTx() {
		return true
	}

	return tx.Verify()
}

// Init a dagchain used in daemon
func Init() (*DagChain, error) {
	// db will be closed in the main() function
	db, dbPruning, dbOthers, err := LoadDB()

	if err != nil {
		return nil, err
	}
	/*tips, err := FetchTips(db)
	if err != nil {
		return nil, err
	}*/
	var tips map[int]*Tip
	//address := []byte(config.GENESIS_TO_ADDRESS)

	// if tips is nil, initialize the new dagchain
	if tips == nil {
		tips, err = InitializeGenesisTxs(dbOthers)
		if err!=nil {
			return nil, err
		}
	}

	return &DagChain{tips, db, dbPruning, dbOthers}, nil
}

/*
	--------------------
	  dagchain related
	--------------------
*/
// Initialize the 1000 genesis transactions and 1000 addresses, each address has 10000 coins
// @parameter address: receiver of genesis tx
func InitializeGenesisTxs(db *leveldb.DB) (map[int]*Tip, error) {
	var tips = make(map[int]*Tip)
	for i, addr:= range(config.GenesisAddresses) {
		genesisTx := NewGenesisTx(config.GENESIS_VALUE, []byte(addr))
		genesisTx.Nonce = genesisTx.Pow()
		genesisTx.Hash = genesisTx.HashTx()
		genesisID := genesisTx.Hash
		genesisTx.SenderNum = 0-i-1

		acc := CreateNewAccount(addr, genesisTx.Value)
		acc.LastIdNo = -1
		AccountMap[i] = acc
		AccountAddr2Int[addr] = i
		AccountInt2Addr[i] = addr

		GXs[genesisTx.Number] = genesisTx

		if i < config.TIPS_COUNT {
			//tips = append(tips, &Tip{genesisTx.Number, genesisID, false, false})
			tips[i] = &Tip{genesisTx.Number, genesisID, false, false}
		}
		// write to db
		if err := StoreTx(db, *genesisTx); err != nil {
			return nil, err
		}
		if err := StoreAccount(db, *acc); err != nil {
			return nil, err
		}
		/*if err := StoreTips(db, tips); err != nil {
			return nil, err
		}*/
	}

	for i:= 0; i< config.GENESIS_ADDR_COUNT; i++ {
		log.Println(AccountInt2Addr[i])
		log.Println(AccountAddr2Int[AccountInt2Addr[i]])
	}

	tips[5].Verified = true
	GXs[5].SetVerification(true)
	return tips, nil
}

/*
	--------------
     tips related
	--------------
*/

//serialize tips
func (tip Tip) Serialize() []byte {
	var encode bytes.Buffer

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(tip)

	if err != nil {
		log.Panic("tip encode fail:", err)
	}

	return encode.Bytes()
}

//deserialize tips
func DeserializeTip(data []byte) Tip {
	var tip Tip

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&tip)
	if err != nil {
		log.Panic("tip decode fail:", err)
	}

	return tip
}

// @return tips of DagChain TODO
func (dc *DagChain) FindTips() map[int]*Tip {

	return dc.Tips
}

// @return last tx's hash of special Account address
func (dc *DagChain) FindParent(account string) [32]byte {
	// Todo read db file of Account to find an Account

	acc, err := FetchAccount(dc.DBOthers, account)
	if err != nil {
		log.Panic(err)
	}

	return acc.LastId

	/*lastTx, errTip := database.FetchTx(dc.DB, acc.LastId)
	if errTip != nil {
		log.Panic(errTip)
	}

	return lastTx*/

	/*for _, tip := range tips {
		tx := transaction.FindTxByHash(tip.TipID)
		if bytes.Equal(tx.Sender, Account) {
			return tip, tx
		}
	}*/
}

// select two tips to verify, randomly
func (dc *DagChain) SelectTips(senderNo int) ([2]int, [2]*Tip) {
	allTips := dc.FindTips()
	var tipsSlice []int = make([]int, len(allTips))
	i := 0
	for n:= range allTips {
		tipsSlice[i] = n
		i++
	}
	//log.Printf("alltips: %v\n", allTips)

	var twoTips [2]*Tip
	var index [2]int

	rand.Seed(time.Now().UnixNano())
	first := rand.Intn(len(allTips))
	//log.Printf("first: %d\n", first)

	firstTipNum := tipsSlice[first]

	//log.Printf("firstTipNum: %d\n", firstTipNum)

	//log.Printf("allTips[firstTipNum].TxNum: %d\n", allTips[firstTipNum].TxNum)

	for GXs[allTips[firstTipNum].TxNum].FetchSenderNum() == senderNo {
		first = rand.Intn(len(allTips))
		firstTipNum = tipsSlice[first]
	}

	second := rand.Intn(len(allTips))
	secondTipNum := tipsSlice[second]
	for GXs[allTips[secondTipNum].TxNum].FetchSenderNum() == senderNo ||
		GXs[allTips[secondTipNum].TxNum].FetchSenderNum() == GXs[allTips[firstTipNum].TxNum].FetchSenderNum() {
		second = rand.Intn(len(allTips))
		secondTipNum = tipsSlice[second]
	}
	/*for first == second {
		log.Println("first == second")
		second = rand.Intn(len(allTips))
	}*/

	index[0] = allTips[firstTipNum].TxNum
	index[1] = allTips[secondTipNum].TxNum

	twoTips[0] = allTips[firstTipNum]
	twoTips[1] = allTips[secondTipNum]


	/*twoTips := [2]*Tip {
		&Tip{[]byte("tseafooler111111111"), false},
		&Tip{[]byte("tseafooler222222222"), false},
	}
	index := [2]int {2, 3}*/

	return index, twoTips
}

// verify selected tips tx
func (dc *DagChain) VerifyTips(tips [2]*Tip) bool {
	// verify selected tx
	/*for _, t := range tips {
		tipTx, err := FetchTx(dc.DB, t.TipID)
		if err != nil {
			log.Panic(err)
		}
		if tipTx.Verify() == false {
			return false
		}
	}*/

	return true
}

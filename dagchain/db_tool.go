package dagchain

import (
	"bytes"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"strconv"
	"errors"
)

// nq: store tx to db
func StoreTx(db *leveldb.DB, tx Transaction) error {
	blockdata := tx.Serialize()
	_ = db.Put(bytes.Join([][]byte{[]byte("t"),tx.Hash[:]},[]byte{}),[]byte(blockdata),nil)
	return nil
}

// for test: store tx to db
func StoreTx2(db *leveldb.DB, gt GeneralTx) error {
	blockdata := gt.Serialize()
	numberString := strconv.Itoa(gt.FetchNumber())
	numberBytes := []byte(numberString)
	key := bytes.Join([][]byte{[]byte("t"), numberBytes}, []byte{})
	err := db.Put(key,[]byte(blockdata), nil)
	if err!=nil {
		return errors.New("store tx error")
	}
	//_ = db.Put(bytes.Join([][]byte{[]byte("t"),tx.Hash[:]},[]byte{}),[]byte(blockdata),nil)

	if gt.FetchNumber() == 880 {
		fmt.Println("key 880 in StoreTxs:", key)
	}
	return nil
}

func StoreTxEth(db *leveldb.DB, gt GeneralTxEth) error {
	blockdata := gt.Serialize()
	numberString := strconv.FormatInt(gt.FetchNumber(), 10)
	numberBytes := []byte(numberString)
	err := db.Put(numberBytes,[]byte(blockdata), nil)
	if err!=nil {
		return errors.New("store tx error")
	}
	//_ = db.Put(bytes.Join([][]byte{[]byte("t"),tx.Hash[:]},[]byte{}),[]byte(blockdata),nil)

	/*if gt.FetchNumber() == 880 {
		fmt.Println("key 880 in StoreTxs:", numberBytes)
	}*/
	return nil
}

// for test: remove tx from db
func RemoveTx2(db *leveldb.DB, gtID int) error {
	numberString := strconv.Itoa(gtID)
	numberBytes := []byte(numberString)
	key := bytes.Join([][]byte{[]byte("t"), numberBytes}, []byte{})
	err := db.Delete(key,nil)
	if err != nil {
		return errors.New("remove failed")
	}
	return nil
}

func RemoveTxEth(db *leveldb.DB, gtID int64) error {
	numberString := strconv.FormatInt(gtID, 10)
	numberBytes := []byte(numberString)
	err := db.Delete(numberBytes,nil)
	if err != nil {
		return errors.New("remove failed")
	}
	return nil
}

// for test: get tx via txid from db
func FetchTx2(db *leveldb.DB, txid int) bool {
	numberString := strconv.Itoa(txid)
	numberBytes := []byte(numberString)
	key := bytes.Join([][]byte{[]byte("t"), numberBytes}, []byte{})
	log.Println("key:", key)
	_, err :=db.Get(key, nil)
	if err != nil {
		log.Println("FetchTx2 error: ", err.Error())
		return false
	}
	return true
}

// for test: store tcl to db
// the new tcl will overwrite the old tcl
func StoreTCL2(db *leveldb.DB, tcl *TxContentList) error {
	key := []byte("l")
	log.Println("key:", key)
	err :=db.Put(key, tcl.Serialize(), nil)
	if err != nil {
		return errors.New("storeTCL error")
	}
	return nil
}

func StoreTCLEth(db *leveldb.DB, tcl *TxContentEthList) error {
	key := []byte("l")
	log.Println("key:", key)
	err :=db.Put(key, tcl.Serialize(), nil)
	if err != nil {
		return errors.New("storeTCL error")
	}
	return nil
}

// for test: fetch tcl to db
func FetchTCL2(db *leveldb.DB) *TxContentList {
	key := []byte("l")
	log.Println("key:", key)
	data, err :=db.Get(key, nil)
	if err != nil {
		log.Println("FetchTCL error: ", err.Error())
		return nil
	}
	tcl := DeserializeTCL(data)
	return &tcl
}

func FetchTCLEth(db *leveldb.DB) *TxContentEthList {
	key := []byte("l")
	log.Println("key:", key)
	data, err :=db.Get(key, nil)
	if err != nil {
		log.Println("FetchTCL error: ", err.Error())
		return nil
	}
	tcl := DeserializeTCLEth(data)
	return &tcl
}

// for test: add a new tc to db
func RefreshTCL2(db *leveldb.DB, tc *TxContent) error {
	key := []byte("l")
	log.Println("key:", key)
	tcl := FetchTCL2(db)
	var newTcs []TxContent
	if tcl != nil {
		newTcs = *tcl
	}
	newTcs = append(newTcs, *tc)
	newTcl := TxContentList(newTcs)
	return StoreTCL2(db, &newTcl)
}

func RefreshTCLEth(db *leveldb.DB, tc *TxContentEth) error {
	key := []byte("l")
	log.Println("key:", key)
	tcl := FetchTCLEth(db)
	var newTcs []TxContentEth
	if tcl != nil {
		newTcs = *tcl
	}
	newTcs = append(newTcs, *tc)
	newTcl := TxContentEthList(newTcs)
	return StoreTCLEth(db, &newTcl)
}

// for test: compact the db
func CompactDB(db *leveldb.DB) error {
	dbRange := util.BytesPrefix([]byte("t"))
	return db.CompactRange(*dbRange)
}

// nq: get tx via hash from db
func FetchTx(db *leveldb.DB, hash []byte) (*Transaction, error) {
	data,err :=db.Get(bytes.Join([][]byte{[]byte("t"),hash},[]byte{}),nil)
	if err != nil {
		return nil, err
	}
	tx := DeserializeTx(data)
	return &tx, nil
}

// nq: store account to db
func StoreAccount(db *leveldb.DB, acc Account) error {
	Accountdata := acc.Serialize()
	key := bytes.Join([][]byte{[]byte("a"), acc.Account[:]},[]byte{})
	value := []byte(Accountdata)
	log.Println("Store key: ", key)
	if err := db.Put(key ,value,nil); err!=nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// nq: store accountEth to db
func StoreAccountEth(db *leveldb.DB, acc AccountEth) error {
	accountData := acc.Serialize()
	key := []byte(acc.Account)
	log.Println("Store key: ", key)
	if err := db.Put(key, accountData,nil); err!=nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// nq: get account via address from db
func FetchAccount(db *leveldb.DB, addr string) (*Account, error) {
	key := bytes.Join([][]byte{[]byte("a"), []byte(addr)},[]byte{})
	//fmt.Println("Fetch key: ", key)
	data, err := db.Get(key,nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	acc := DeserializeAcc(data)
	return &acc, nil
}

// nq: get account via address from db
func FetchAccountEth(db *leveldb.DB, addr string) (*AccountEth, error) {
	key := []byte(addr)
	//fmt.Println("Fetch key: ", key)
	data, err := db.Get(key,nil)
	if err != nil {
		return nil, err
	}
	acc := DeserializeAccEth(data)
	return &acc, nil
}

// nq: store all tips to db
func StoreTips(db *leveldb.DB, tips []*Tip) error {
	for i := 0; i < len(tips); i++ {
		data := tips[i].Serialize()
		_ = db.Put(bytes.Join([][]byte{[]byte("tip"),tips[i].TipID[:]},[]byte{}),[]byte(data),nil)
	}
	return nil
}

// Add a new tip
func PutTip(db *leveldb.DB, tip Tip) error {
	Tipdata := tip.Serialize()
	err := db.Put(bytes.Join([][]byte{[]byte("tip"),tip.TipID[:]},[]byte{}),[]byte(Tipdata),nil)
	if err != nil {
		return errors.New("Add Tip failed.")
	}
	return nil
}

// Delete a specified tip with tip tx hash
func RemoveTip(db *leveldb.DB, tipID []byte) error {
	err := db.Delete(bytes.Join([][]byte{[]byte("tip"),tipID},[]byte{}),nil)
	if err != nil {
		return errors.New("Remove failed.")
	}
	return nil
}

// nq: get all tips from db
func FetchTips(db *leveldb.DB) ([]*Tip, error) {
	var temptips []*Tip
	iter := db.NewIterator(util.BytesPrefix([]byte("tip")),nil)
	for iter.Next() {
		data := DeserializeTip(iter.Value())
		temptips = append(temptips, &data)
	}
	iter.Release()
	tips := temptips[:]
	return tips,nil
}


// open db instance or create a db if not exist
func LoadDB() (*leveldb.DB, *leveldb.DB, *leveldb.DB, *leveldb.DB, *leveldb.DB, *leveldb.DB, error) {
	log.Println("LoadDB() function is called")
	db, err := leveldb.OpenFile("Jightdb",nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	dbPruning, err := leveldb.OpenFile("JightdbMerging",nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dbOthers, err := leveldb.OpenFile("JightdbOthers",nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dbEth, err := leveldb.OpenFile("JightdbEth",nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	dbEthPruning, err := leveldb.OpenFile("JightdbEthMerging",nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dbEthOthers, err := leveldb.OpenFile("JightdbEthOthers",nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	return db, dbPruning, dbOthers, dbEth, dbEthPruning, dbEthOthers, nil
}


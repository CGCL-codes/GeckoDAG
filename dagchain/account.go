package dagchain

import (
	"Jight/config"
	"bytes"
	"encoding/gob"
	"log"
	"strconv"
)

var AccountMap = make(map[int]*Account, config.GENESIS_ADDR_COUNT)
var AccountAddr2Int = make(map[string]int, config.GENESIS_ADDR_COUNT)
var AccountInt2Addr = make(map[int]string, config.GENESIS_ADDR_COUNT)

var AccountEthMap = make(map[int64]*AccountEth)

var accountNo int = 0
var accountNoEth int64 = 0

type Account struct {
	Account [34]byte // address of account
	AccountNo int	// No. of the account
	LastId  [32]byte // hash id of account's last tx
	LastIdNo int 	// No. of the account's last tx
	Balance int // Balance of this Account
	TxCount int // Count of transactions sent by the account
	LatestTU *TU // Latest transaction union of this account
	WithoutMergeIds [config.MERGE_PERIOD-1]int // Ids of transactions without merge
	WithoutPruneIds []int // Ids of transactions send by the account while without pruning
}

type AccountEth struct {
	Account string // address of account
	AccountNo int64	// No. of the account
	LastId  [64]byte // hash id of account's last tx
	LastIdNo int64	// No. of the account's last tx
	Balance int // Balance of this Account
	TxCount int // Count of transactions sent by the account
	LatestTU *TUEth // Latest transaction union of this account
	WithoutMergeIds [config.MERGE_PERIOD-1]int64 // Ids of transactions without merge
	WithoutPruneIds []int64 // Ids of transactions send by the account while without pruning
}

func PrintAccounts() string {
	var returnString string
	returnString = returnString + "\n"
	for i:=0; i< len(AccountMap); i++ {
		returnString = returnString + "Account number: " +  strconv.Itoa(i) + " "
		returnString = returnString + "LastIdNo: " + strconv.Itoa(AccountMap[i].LastIdNo) + " "
		returnString = returnString + "TxCount: " +  strconv.Itoa(AccountMap[i].TxCount)+ " "
		if AccountMap[i].LatestTU != nil {
			returnString = returnString + "LatestTU: " +  strconv.Itoa(AccountMap[i].LatestTU.Number)+ " "
		}
		returnString = returnString + "\n"
	}
	return returnString
}

func CreateNewAccount(addr string, balance int) *Account {
	addrBytes := []byte(addr)
	var accountBytes [34]byte
	copy(accountBytes[:], addrBytes)

	acc := &Account{accountBytes, accountNo, [32]byte{}, 0, balance,
		0, nil, [config.MERGE_PERIOD-1]int{},[]int{}}
	accountNo++
	return acc
}

func CreateNewAccountEth(addr string, balance int, isSender bool) *AccountEth {
	var txCount int
	if isSender {
		txCount = 1
	} else {
		txCount = 0
	}
	acc := &AccountEth{addr, accountNoEth, [64]byte{}, 0, balance,
		txCount, nil, [config.MERGE_PERIOD-1]int64{},[]int64{}}
	accountNoEth++
	return acc
}

// serialize Account
func (acc Account) Serialize() []byte {
	var encode bytes.Buffer

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(acc)

	if err != nil {
		log.Panic("Account encode fail:", err)
	}

	return encode.Bytes()
}

// serialize AccountEth
func (acc AccountEth) Serialize() []byte {
	var encode bytes.Buffer

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(acc)

	if err != nil {
		log.Panic("AccountEth encode fail:", err)
	}

	return encode.Bytes()
}

// Deserialize Account
func DeserializeAcc(data []byte) Account {
	var acc Account

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&acc)
	if err != nil {
		log.Panic("Account decode fail:", err)
	}

	return acc
}

// Deserialize AccountEth
func DeserializeAccEth(data []byte) AccountEth {
	var acc AccountEth

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&acc)
	if err != nil {
		log.Panic("AccountEth decode fail:", err)
	}

	return acc
}
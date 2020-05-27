package dagchain

import (
	"Jight/config"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"testing"
	"time"
)

func TestTransaction_Serialize(t *testing.T) {
	//address := "158ewJ1itTAAE8gy1Hk45JdAmqvbpd0001"
	tx := NewTx([2]int{1,2}, [32]byte{}, [2][32]byte{}, [32]byte{}, "", 0, 1000, "", false)
	txBytes := tx.Serialize()
	fmt.Printf("Length of serialized tx: %d\n", len(txBytes))
	hash := tx.HashTx()
	fmt.Printf("Length of tx hash: %d\n", len(hash))

	tx1 := NewTx([2]int{1,2}, config.MOCK_TX, [2][32]byte{}, config.MOCK_TX,
		"158ewJ1itTAAE8gy1Hk45JdAmqvbpd0006", 0, 1000, "158ewJ1itTAAE8gy1Hk45JdAmqvbpd0007", false)
	tx1Bytes := tx1.Serialize()
	fmt.Printf("Length of serialized tx: %d\n", len(tx1Bytes))
	hash1 := tx1.HashTx()
	fmt.Printf("Length of tx hash: %d\n", len(hash1))
}

func TestTransaction_Sign(t *testing.T) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Error(err.Error())
	}
	tx := NewTx([2]int{1,2}, [32]byte{}, [2][32]byte{}, [32]byte{}, "", 0, 1000, "", false)
	tx.Sign(*private)
	log.Printf("Length of tx signature: %d\n", len(tx.Signature))
	log.Println(tx.Signature)
}

func TestDeserializeTCL(t *testing.T) {
	var tcl TxContentList = []TxContent {
		TxContent{
			Receiver:  config.MOCK_ACCOUNT,
			Value:     1,
			Timestamp: time.Now().Unix(),
			Income:    config.MOCK_TX,
			Nonce:     0,
		},
		TxContent{
			Receiver: config.MOCK_ACCOUNT,
			Value: 2,
			Timestamp: time.Now().Unix(),
			Income: config.MOCK_TX,
			Nonce: 0,
		},
		TxContent{
			Receiver: config.MOCK_ACCOUNT,
			Value: 3,
			Timestamp: time.Now().Unix(),
			Income: config.MOCK_TX,
			Nonce: 0,
		},
		TxContent{
			Receiver: config.MOCK_ACCOUNT,
			Value: 4,
			Timestamp: time.Now().Unix(),
			Income: config.MOCK_TX,
			Nonce: 0,
		},
	}

	var newTC = TxContent {
		Receiver: config.MOCK_ACCOUNT,
		Value: 5,
		Timestamp: time.Now().Unix(),
		Income: config.MOCK_TX,
		Nonce: 0,
	}

	db, err := leveldb.OpenFile("Testdb",nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if RefreshTCL2(db, &newTC)!=nil {
		t.Error("refresh tcl error at the first time")
	}
	tcl2 := FetchTCL2(db)
	log.Println(tcl2)

	if StoreTCL2(db, &tcl)!=nil {
		t.Error("store tcl error")
	}
	tcl2 = FetchTCL2(db)
	if tcl2 == nil {
		t.Error("fetch tcl error")
	}
	log.Println(tcl2)

	if RefreshTCL2(db, &newTC)!=nil {
		t.Error("refresh tcl error")
	}
	tcl2 = FetchTCL2(db)
	if tcl2 == nil {
		t.Error("fetch tcl error")
	}
	log.Println(tcl2)
}

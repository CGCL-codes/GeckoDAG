package dagchain

import (
	"Jight/config"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"testing"
)

func TestStoreTx(t *testing.T) {
	db, err := leveldb.OpenFile("Testdb",nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	for i:= 0; i<100; i++ {
		var gt GeneralTx = NewTx([2]int{1, 2}, config.MOCK_TX, [2][32]byte{config.MOCK_TX, config.MOCK_TX}, config.MOCK_TX, "", 0, 1000, "", false)
		log.Println(gt.FetchNumber())
		StoreTx2(db, gt)
	}

	dbRange := util.BytesPrefix([]byte("t"))

	iter := db.NewIterator(dbRange, nil)

	for iter.Next() {
		log.Printf("k: %v, v: %v\n", iter.Key(), iter.Value())
	}
}

func TestDeleteTx(t *testing.T) {
	db, err := leveldb.OpenFile("Testdb",nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	for i:= 0; i<999000; i++ {
		err = RemoveTx2(db, i)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

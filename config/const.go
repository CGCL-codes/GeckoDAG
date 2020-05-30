package config

import (
	"fmt"
	"math"
)

// TODO how to set a pow target which need few computational
const POW_TARGET_BITS int = 24
const MAX_NONCE uint32 = math.MaxUint32

const NETWORK_LATENCY int = 10

// TODO genesis tx value
const GENESIS_VALUE int = 1000
const GENESIS_ADDR_COUNT int = 1000
const TIPS_COUNT int = 240

const MERGE_PERIOD int = 100


// Default data storage location
const LOC string = "E:\\Jight\\"

// TODO db location storing dagchain
const DB_LOCATION string = LOC + "db"


// Wallet const
const (
	VERSION = byte(0x00)
	ADDRESS_CHECK_SUM_LEN = 4
	WALLET_FILE = "wallet_%s.dat"
)

var GenesisAddresses [GENESIS_ADDR_COUNT]string

// convert a integer (< 1000) to a 4 byte string
func convertIntTo4Byte(i int) string {
	if i < 10 {
		return fmt.Sprintf("000%d", i)
	} else if i < 100 {
		return fmt.Sprintf("00%d", i)
	} else if i < 1000 {
		return fmt.Sprintf("0%d", i)
	} else {
		return fmt.Sprintf("%d", i)
	}
}

func init() {
	for i:= 0; i< int(GENESIS_ADDR_COUNT); i++ {
		GenesisAddresses[i] = "158ewJ1itTAAE8gy1Hk45JdAmqvbpd" + convertIntTo4Byte(i)
	}
}

var MOCK_ACCOUNT [34]byte = [34]byte{'1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1',
	'1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}

var MOCK_TX [32]byte = [32]byte{'1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1',
	'1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}

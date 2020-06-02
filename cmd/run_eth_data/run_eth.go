package main

import (
	jrpc "Jight/rpc"
	"encoding/csv"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

const (
	FileDirPath = "/home/seafooler/ethTx_data"
	IntermediateSuffixrName = "_NormalTransaction"
	FileSuffix = "_NormalTransaction.csv"
	Factor = 1000000
)

func ReadCSVFile(fileName string) (*csv.Reader, error) {
	// 1. Open the file
	recordFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return nil, err
	}

	// 2. Initialize the reader
	reader := csv.NewReader(recordFile)
	return reader, nil
}

func MakeFullFileName(num int) string {
	startNum := num*Factor
	stopNum := (num+1)*Factor - 1
	prefix := strconv.Itoa(startNum) + "to" + strconv.Itoa(stopNum)
	fullFileName := FileDirPath + "/" + prefix + IntermediateSuffixrName + "/" + prefix + FileSuffix
	fmt.Println("FileName: ", fullFileName)
	return fullFileName
}

func sendTx(rpcPort, from, to, income string, amount int) {
	client, err := rpc.DialHTTP("tcp", "localhost:"+ rpcPort)
	defer client.Close()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	sc := jrpc.SendEthCmd{
		From: from,
		To: to,
		IncomeTx: income,
		Amount: amount,
	}

	var sr jrpc.SendEthReply

	err = client.Call("Jightd.SendEth", sc, &sr)
	if err != nil {
		log.Fatal("Send error: ", err)
	} else {
		//fmt.Printf("Transaction hash: %s\n", sr.Tx)
	}
}

func refreshTips(rpcPort string) {
	client, err := rpc.DialHTTP("tcp", "localhost:"+ rpcPort)
	defer client.Close()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	rtc := jrpc.RefreshTipsCMD{}

	var rtr jrpc.RefreshTipsReply

	err = client.Call("Jightd.RefreshTipsEth", rtc, &rtr)

	if err != nil {
		log.Fatal("Send error: ", err)
	} else {
		fmt.Printf("Before refresh, tips cnt: %d; after refresh, tips cnt: %d\n", rtr.TipsCntBefore,
		rtr.TipsCntAfter)
	}
}

func main() {
	argsCnt := os.Args
	if len(argsCnt) != 3 {
		log.Fatal("Exactly two arguments are needed")
	}

	fileNum := os.Args[1]
	num, err := strconv.Atoi(fileNum)
	if err != nil {
		log.Fatal(err)
	}

	fullFileName := MakeFullFileName(num)

	reader, err := ReadCSVFile(fullFileName)
	if err != nil {
		log.Fatal(err)
	}

	rpcPort := os.Args[2]

	var blockNum = ""

	for i := 0; ; i++{
		line, err := reader.Read()
		if err!=nil {
			break
		}

		if line[4] != "None" {
			sendTx(rpcPort, line[3], line[4], "", 100)
			//fmt.Printf("%d th line, from: %s\n", i, line[3])
			if i % 100 == 0 {
				refreshTips(rpcPort)
			}
			if blockNum != line[0] {
				fmt.Printf("Process block: %s\n", blockNum)
				blockNum = line[0]
			}
		}
	}

}


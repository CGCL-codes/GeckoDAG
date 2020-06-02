package main

import (
	"Jight/api"
	"Jight/dagchain"
	"Jight/p2p"
	jrpc "Jight/rpc"
	"Jight/utils"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

var cmd api.CMDDeamon

var dc *dagchain.DagChain

var dcEth *dagchain.DagChainEth

// initialize the variables from the arguments
func init() {
	cmd.Run()
	fmt.Println("Cmd variables: %d, %d, %s, %s, %s, %s\n",
		cmd.RPCPort, cmd.P2PPort, cmd.Pid, cmd.FullAddrsPath, cmd.PrivKeyPath, cmd.TargetPath)
	if cmd.HelpFlag {
		os.Exit(0)
	}
	// initiazlie the dagchain
	var err error
	dc, dcEth, err = dagchain.Init()
	if err != nil {
		log.Fatal(err)
	}
}

// start the rpc server
func init() {
	jd := new(jrpc.Jightd)
	jd.DC = dc
	jd.DCEth = dcEth
	err := rpc.Register(jd)
	if err != nil {
		log.Fatal("Format of service Jightd isn't correct.", err)
	}

	rpc.HandleHTTP()

	listener, e := net.Listen("tcp", ":" + strconv.Itoa(cmd.RPCPort))
	if e != nil{
		log.Fatal("Listen error: ", e)
	}

	log.Printf("Serving RPC server on port %d", cmd.RPCPort)
	go http.Serve(listener, nil)

}

// start the p2p listener
func init() {
	pk, err := utils.GetPrivKey(cmd.CreatePK, cmd.PrivKeyPath, 0)
	if err != nil {
		log.Fatal(err)
	}

	host, err := p2p.MakeBasicHost(cmd.P2PPort, pk, cmd.FullAddrsPath)
	if err!=nil {
		log.Fatal(err)
	}

	txDealerHandler := dc.AddGT

	p2p.OpenP2PListen(cmd.Pid, host, txDealerHandler)
	log.Printf("Open a port %d for p2p connection, fullAddr is stored in %s\n", cmd.P2PPort, cmd.FullAddrsPath)

	targetAddrs, err := utils.ReadStringsFromFile(cmd.TargetPath)
	if err!=nil {
		log.Fatal(err)
	}

	for _, a := range targetAddrs {
		log.Println("Target address: ", a)
		err := p2p.ConnectP2PNode(&a, cmd.Pid, host, txDealerHandler)
		if err!=nil {
			log.Fatal(err)
		}
		log.Println(fmt.Sprintf("Successfully connect to the server: %s", a))
	}
}

/*// write the transaction to the local DB
func writeTx2DB(tx *transaction.Transaction) error {
	log.Println("Transaction: ", *tx)
	return nil
}*/

func main() {
	defer dc.DB.Close()
	defer dc.DBMerging.Close()
	defer dc.DBOthers.Close()
	defer dcEth.DB.Close()
	defer dcEth.DBMerging.Close()
	defer dcEth.DBOthers.Close()
	/*tx := transaction.Transaction{
		Hash: []byte("h1111111111111"),
		Parent: []byte("p22222222222222"),
		Validate: [][]byte{[]byte("v333333333333"), []byte("v44444444444")},
		Sender: []byte("s44444444444"),
		Value: 12,
		Receiver: []byte("r555555555555"),
		Nonce: 199883,
		Timestamp: 3232342424242,
		Signature: []byte("si666666666666"),
	}
	p2p.SyncTx(&tx)*/
	select{}
}
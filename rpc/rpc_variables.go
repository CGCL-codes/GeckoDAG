package rpc

import "Jight/dagchain"

type Jightd struct{
	DC *dagchain.DagChain
	DCEth *dagchain.DagChainEth
}

type GetBalanceCMD struct {
	Address string
}

type GetBalanceReply struct {
	Balance string
}

type ListAddressCMD struct {
}

type ListAddressReply struct {
	Addresses []string
}

type SendCmd struct {
	From int
	To int
	IncomeTx string
	Amount int
}

type SendEthCmd struct {
	From string
	To string
	IncomeTx string
	Amount int
}

type SendReply struct {
	Tx string
}

type SendEthReply struct {
	Tx string
}

type CreateDagChainCMD struct {
	Address string
}

type CreateDagChainReply struct {
	CdcrField string
}

type CreateWalletCMD struct {
	Address string
}

type CreateWalletReply struct {
	CwrField string
}

type CreateBatchTxsCMD struct {
	Times int
	Speed int
}

type CreateBatchTxsReply struct {
}

type RefreshTipsCMD struct {
}

type RefreshTipsReply struct {
	TipsCntBefore int
	TipsCntAfter int
}

type FetchTxsByAccountCMD struct {
	AccountId int
}

type FetchTxsByAccountReply struct {
	TxsNumber []int
}

type FetchTUsByAccountCMD struct {
	AccountId int
}

type FetchTUsByAccountReply struct {
	TUsNumber []int
}

type FetchTxCMD struct {
	TxId int
}

type FetchTxReply struct {
	Existed bool
}

// The following variables are just for test
var balancesList = map[string]string{
	"fsdjklfjskljfkldsjklfdjsklf": "123243",
	"fjiosdjfioqrewrewrewtrtrggg": "896969",
}

var addressList = []string {
	"fsdjklfjskljfkldsjklfdjsklf",
	"fjiosdjfioqrewrewrewtrtrggg",
	"r8i9w0jfsfjsdkfjdskfjh8r9ww",
}




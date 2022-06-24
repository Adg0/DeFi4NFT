package d4t

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

type Conf struct {
	Network     string `json:"network,omitempty"`
	RootDir     string `json:"dir,omitempty"`
	AbiDir      string `json:"abi_dir,omitempty"`
	TealDir     string `json:"teal_dir,omitempty"`
	CodecDir    string `json:"codec_dir,omitempty"`
	DryrunDir   string `json:"dryrun_dir,omitempty"`
	ResponseDir string `json:"response_dir,omitempty"`
	Address     string `json:"address"`
	Token       string `json:"token"`
	Header      string `json:"header"`
}

var (
	rootDir = "./conf.json"
	net     = getConf("network")
)

func InitAlgodClient() (*algod.Client, error) {
	// Initialize an algodClient
	commonClient, err := common.MakeClient(getConf("address"), getConf("header"), getConf("token"))
	if err != nil {
		log.Fatalf("Failed to make common client: %+v", err)
	}
	return (*algod.Client)(commonClient), nil
}

func debugAppCall(algodClient *algod.Client, atc future.AtomicTransactionComposer, dryrunDump, response string) []future.ABIMethodResult {
	// gather signatures
	stxns, _ := atc.GatherSignatures()
	stx := make([]types.SignedTxn, len(stxns))
	for i, sigTxns := range stxns {
		stxn := types.SignedTxn{}
		msgpack.Decode(sigTxns, &stxn)
		stx[i] = stxn
	}

	// Create the dryrun request object
	dryrunRequest, err := future.CreateDryrun(algodClient, stx, nil, context.Background())
	if err != nil {
		log.Fatalf("Failed creating dryrun: %+v", err)
	}

	// Pass dryrun request to algod server
	dryrunResponse, _ := algodClient.TealDryrun(dryrunRequest).Do(context.Background())

	// Inspect the response to check result
	os.WriteFile(dryrunDump, msgpack.Encode(dryrunRequest), 0666)
	drr, err := json.MarshalIndent(dryrunResponse, "", "")
	if err != nil {
		log.Fatalf("Failed JSON marshal indent: %+v", err)
	}
	os.WriteFile(response, drr, 0666)

	ret, err := atc.Execute(algodClient, context.Background(), 2)
	if err != nil {
		log.Fatalf("Failed to execute call: %+v", err)
	}
	for _, r := range ret.MethodResults {
		log.Printf("%s returned %+v", r.TxID, r.ReturnValue)
	}
	return ret.MethodResults
}

func ConfigASA(algodClient *algod.Client, acct crypto.Account, mngID, d4tID, lqtID, assetID uint64) (err error) {
	// Get information about the asset
	asset, err := algodClient.GetAssetByID(assetID).Do(context.Background())
	if err != nil {
		log.Fatalf("Found error getting Asset:%#v\n", err)
	}

	// Stop attempt if account isn't Manager of asset
	if asset.Params.Manager != acct.Address.String() {
		log.Fatalf("You are not the manager of this asset!\nOnly address %s can configure this asset.", asset.Params.Manager)
	}

	// Change admin addresses to manager, d4t and liquidator contract
	new_manager := crypto.GetApplicationAddress(mngID).String()
	new_freeze := crypto.GetApplicationAddress(d4tID).String()
	new_clawback := crypto.GetApplicationAddress(lqtID).String() // liquidatorAddr
	strictEmptyAddressChecking := true
	note := []byte(nil)
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}
	txn, err := future.MakeAssetConfigTxn(acct.Address.String(), note, txParams, assetID, new_manager, asset.Params.Reserve, new_freeze, new_clawback, strictEmptyAddressChecking)
	if err != nil {
		log.Fatalf("Failed to send transaction MakeAssetConfig Txn: %s\n", err)
	}

	// sign the transaction
	err = signSendWait(algodClient, acct.PrivateKey, txn)
	return
}

func OptinASA(algodClient *algod.Client, acct crypto.Account, assetID uint64) (err error) {
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}
	txn, err := future.MakeAssetAcceptanceTxn(acct.Address.String(), []byte(nil), txParams, assetID)
	if err != nil {
		log.Fatalf("Failed to send transaction MakeAssetAcceptance Txn: %s\n", err)
	}
	err = signSendWait(algodClient, acct.PrivateKey, txn)
	return
}

func signSendWait(algodClient *algod.Client, sk ed25519.PrivateKey, txn types.Transaction) (err error) {
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(sk, txn)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %s\n", err)
	}
	log.Printf("Transaction ID: %s\n", txid)

	// Broadcast the transaction to the network
	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s\n", err)
	}

	// Wait for transaction to be confirmed
	waitForConfirmation(txid, algodClient)
	return
}

// fetch data from file
func getConf(selector string) (ret string) {
	r, err := os.Open(rootDir)
	if err != nil {
		log.Fatalf("Failed to open app root file: %+v", err)
	}

	d, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("Failed to read file: %+v", err)
	}

	// contract location
	dir := &Conf{}
	if err := json.Unmarshal(d, dir); err != nil {
		log.Fatalf("Failed to unmarshal contract: %+v", err)
	}

	switch selector {
	case "abi":
		ret = dir.AbiDir
	case "dryrun":
		ret = dir.DryrunDir
	case "response":
		ret = dir.ResponseDir
	case "codec":
		ret = dir.CodecDir
	case "teal":
		ret = dir.TealDir
	case "root":
		ret = dir.RootDir
	case "network":
		ret = dir.Network
	case "address":
		ret = dir.Address
	case "token":
		ret = dir.Token
	case "header":
		ret = dir.Header
	default:
		ret = dir.RootDir
	}
	return
}

func getAbi(app string) (abi string) {
	abi = getConf("abi")
	switch app {
	case "d4t":
		abi += "d4t.json"
	case "lqt":
		abi += "lqt.json"
	case "mng":
		abi += "mng.json"
	default:
		abi += "d4t.json"
	}
	return
}

func getTeal(app string) (teal string) {
	teal = getConf("teal")
	switch app {
	case "d4t":
		teal += "d4tApp.teal"
	case "d4tClear":
		teal += "d4tClear.teal"
	case "lqt":
		teal += "liquidatorApp.teal"
	case "mng":
		teal += "managerApp.teal"
	case "lsig":
		teal += "logicSigDelegated.teal"
	default:
		teal += "clearState.teal"
	}
	return
}

// Make 44t application call to earn USDCa at 3%
func Earn(algodClient *algod.Client, acct crypto.Account, xids []uint64, aamt, lvr uint64, lsa []byte) (err error) {
	contract, err := getContract("app")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}

	signer := future.BasicAccountTransactionSigner{Account: acct}

	mcp := future.AddMethodCallParams{
		AppID:           contract.Networks[net].AppID,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "earn"), []interface{}{xids, aamt, lvr, lsa}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}

	debugAppCall(algodClient, atc, "./dryrun/earn.msgp", "./dryrun/response/earn.json")
	return
}

func Optin(algodClient *algod.Client, acct crypto.Account, app uint64, setContract string) (err error) {
	contract, err := getContract(setContract)
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}

	signer := future.BasicAccountTransactionSigner{Account: acct}

	mcp := future.AddMethodCallParams{
		AppID:           contract.Networks[net].AppID,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.OptInOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "optin"), []interface{}{app}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}

	debugAppCall(algodClient, atc, "./dryrun/optin.msgp", "./dryrun/response/optin.json")
	return
}

// Make 44t application call to borrow against provided collateral
func Borrow(algodClient *algod.Client, acct crypto.Account, lender types.Address, usdc, jusd, mng, lqt uint64, xids, camt, lamt []uint64, lsigFile string) (err error) {
	contract, err := getContract("d4t")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}
	txParams.FlatFee = true
	txParams.Fee = types.MicroAlgos(4 * txParams.MinFee)

	signer := future.BasicAccountTransactionSigner{Account: acct}

	mcp := future.AddMethodCallParams{
		AppID:           contract.Networks[net].AppID,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	txParams.Fee = 0
	txn, _ := future.MakeAssetTransferTxn(lender.String(), acct.Address.String(), lamt[0], nil, txParams, "", usdc)
	lsa, err := FetchLsigFromFile(lsigFile)
	if err != nil {
		log.Fatalf("Failed to get lsa from file: %+v", err)
	}
	signerLsa := future.LogicSigAccountTransactionSigner{LogicSigAccount: lsa}
	stxn := future.TransactionWithSigner{Txn: txn, Signer: signerLsa}
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "borrow"), []interface{}{stxn, xids, camt, lamt, lender, xids[0], jusd, mng, lqt}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}

	debugAppCall(algodClient, atc, "./dryrun/borrow.msgp", "./dryrun/response/borrow.json")
	return

}

// Make 44t application call to repay loan and unfreeze asset
func Repay(algodClient *algod.Client, acct crypto.Account, mng, lqt, usdc uint64, xids, ramt []uint64) (err error) {
	contract, err := getContract("d4t")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}
	txParams.FlatFee = true
	txParams.Fee = types.MicroAlgos(3 * txParams.MinFee)

	signer := future.BasicAccountTransactionSigner{Account: acct}

	d4t := contract.Networks[net].AppID
	mcp := future.AddMethodCallParams{
		AppID:           d4t,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	txParams.Fee = 0
	txn, _ := future.MakeAssetTransferTxn(acct.Address.String(), crypto.GetApplicationAddress(d4t).String(), ramt[0], nil, txParams, "", usdc)
	stxn := future.TransactionWithSigner{Txn: txn, Signer: signer}
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "repay"), []interface{}{stxn, xids, ramt, xids[0], mng, lqt}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}

	debugAppCall(algodClient, atc, "./dryrun/repay.msgp", "./dryrun/response/repay.json")
	return
}

// Make 44t application call to claim USDCa for JUSD
func Claim(algodClient *algod.Client, acct crypto.Account, mng, amt, usdc, jusd uint64) (err error) {
	contract, err := getContract("d4t")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}
	txParams.FlatFee = true
	txParams.Fee = types.MicroAlgos(3 * txParams.MinFee)

	signer := future.BasicAccountTransactionSigner{Account: acct}

	d4t := contract.Networks[net].AppID
	mcp := future.AddMethodCallParams{
		AppID:           d4t,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	txParams.Fee = 0
	txn, _ := future.MakeAssetTransferTxn(acct.Address.String(), crypto.GetApplicationAddress(d4t).String(), amt, nil, txParams, "", jusd)
	stxn := future.TransactionWithSigner{Txn: txn, Signer: signer}
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "claim"), []interface{}{stxn, usdc, mng}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}

	debugAppCall(algodClient, atc, "./dryrun/claim.msgp", "./dryrun/response/claim.json")
	return
}

func ConfigureApps(algodClient *algod.Client, acct crypto.Account, lqt, d4t, usdc, jusd uint64) (err error) {
	contract, err := getContract("mng")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}
	txParams.FlatFee = true
	txParams.Fee = types.MicroAlgos(12 * txParams.MinFee)

	signer := future.BasicAccountTransactionSigner{Account: acct}

	mcp := future.AddMethodCallParams{
		AppID:           contract.Networks[net].AppID,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	lqtAddress := crypto.GetApplicationAddress(lqt)
	d4tAddress := crypto.GetApplicationAddress(d4t)
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "config"), []interface{}{lqt, d4t, lqtAddress, d4tAddress, usdc, jusd}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}
	debugAppCall(algodClient, atc, "./dryrun/config.msgp", "./dryrun/response/config.json")
	return
}

// create sub-apps
func CreateApps(algodClient *algod.Client, acct crypto.Account, usdc uint64) (ids [4]uint64, err error) {
	contract, err := getContract("mng")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	lqtClear, err := CompileSmartContractTeal(algodClient, getTeal("clear"))
	if err != nil {
		log.Fatalf("clearState found error, %s", err)
	}
	lqtApproval, err := CompileSmartContractTeal(algodClient, getTeal("lqt"))
	if err != nil {
		log.Fatalf("liquidatorApp found error, %s", err)
	}
	d4tClear, err := CompileSmartContractTeal(algodClient, getTeal("d4tClear"))
	if err != nil {
		log.Fatalf("d4tClear found error, %s", err)
	}
	d4tApproval, err := CompileSmartContractTeal(algodClient, getTeal("d4t"))
	if err != nil {
		log.Fatalf("d4tApp found error, %s", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}
	txParams.FlatFee = true
	txParams.Fee = types.MicroAlgos(2 * txParams.MinFee)

	signer := future.BasicAccountTransactionSigner{Account: acct}

	mcp := future.AddMethodCallParams{
		AppID:           contract.Networks[net].AppID,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	var atc2 future.AtomicTransactionComposer
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "create_liquidator"), []interface{}{lqtApproval, lqtClear}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall, create_liquidator: %+v", err)
	}
	lqt := uint64(0)
	d4t := uint64(0)
	ret := debugAppCall(algodClient, atc, "./dryrun/create_liquidator.msgp", "./dryrun/response/create_liquidator.json")
	lqt = ret[0].ReturnValue.(uint64)

	txParams.Fee = types.MicroAlgos(4 * txParams.MinFee)
	mcp.SuggestedParams = txParams
	err = atc2.AddMethodCall(combine(mcp, getMethod(contract, "create_child"), []interface{}{usdc, d4tApproval, d4tClear, lqt}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall, create_child: %+v", err)
	}

	ret_j := debugAppCall(algodClient, atc2, "./dryrun/create_child.msgp", "./dryrun/response/create_child.json")
	var v []interface{} = ret_j[0].ReturnValue.([]interface{})
	d4t = v[0].(uint64)
	ids[0] = lqt
	ids[1] = d4t
	ids[2] = v[1].(uint64)
	ids[3] = v[2].(uint64)
	updateABI(algodClient, lqt, "lqt")
	updateABI(algodClient, d4t, "d4t")
	return
}

// Fund app
func Fund(algodClient *algod.Client, acct crypto.Account, amt uint64) (err error) {
	contract, err := getContract("mng")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}

	app := contract.Networks[net].AppID
	txn, err := future.MakePaymentTxn(acct.Address.String(), crypto.GetApplicationAddress(app).String(), amt, []byte(""), "", txParams)
	if err != nil {
		log.Fatalf("Failed creating asset: %+v", err)
	}
	signSendWait(algodClient, acct.PrivateKey, txn)
	return
}

// Update smart contract
func Update(algodClient *algod.Client, acct crypto.Account) (err error) {
	contract, err := getContract("mng")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}

	signer := future.BasicAccountTransactionSigner{Account: acct}

	// get approval and clearState as []byte
	clear, err := CompileSmartContractTeal(algodClient, getTeal("clear"))
	if err != nil {
		log.Fatalf("clearState found error, %s", err)
	}
	app, err := CompileSmartContractTeal(algodClient, getTeal("mng"))
	if err != nil {
		log.Fatalf("approval found error, %s", err)
	}

	mcp := future.AddMethodCallParams{
		AppID:           contract.Networks[net].AppID,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.UpdateApplicationOC,
		Signer:          signer,
		ApprovalProgram: app,
		ClearProgram:    clear,
	}

	var atc future.AtomicTransactionComposer
	err = atc.AddMethodCall(mcp)
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}

	debugAppCall(algodClient, atc, "./dryrun/update.msgp", "./dryrun/response/update.json")
	return
}

func SendDusd(algodClient *algod.Client, acct crypto.Account, rec types.Address, jusd uint64) (err error) {
	contract, err := getContract("mng")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}
	txParams.FlatFee = true
	txParams.Fee = 2000

	signer := future.BasicAccountTransactionSigner{Account: acct}

	mcp := future.AddMethodCallParams{
		AppID:           contract.Networks[net].AppID,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "fund"), []interface{}{rec, jusd}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}

	debugAppCall(algodClient, atc, "./dryrun/fund.msgp", "./dryrun/response/fund.json")
	return
}

// Update child smart contract
func ChildUpdate(algodClient *algod.Client, acct crypto.Account, appID uint64, app, clear string) (err error) {
	contract, err := getContract("mng")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}
	txParams.FlatFee = true
	txParams.Fee = 2000

	signer := future.BasicAccountTransactionSigner{Account: acct}

	// get approval and clearState as []byte
	clearState, err := CompileSmartContractTeal(algodClient, clear)
	if err != nil {
		log.Fatalf("clearState found error, %s", err)
	}
	approval, err := CompileSmartContractTeal(algodClient, app)
	if err != nil {
		log.Fatalf("approval found error, %s", err)
	}

	mcp := future.AddMethodCallParams{
		AppID:           contract.Networks[net].AppID,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	var atc future.AtomicTransactionComposer
	atc.AddMethodCall(combine(mcp, getMethod(contract, "update_child_app"), []interface{}{appID, approval, clearState}))

	debugAppCall(algodClient, atc, "./dryrun/update_child.msgp", "./dryrun/response/update_child.json")
	return
}

// Deploy smart contract
func Deploy(algodClient *algod.Client, acct crypto.Account, usdc uint64) (newApp uint64, err error) {
	contract, err := getContract("mng")
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}

	signer := future.BasicAccountTransactionSigner{Account: acct}

	// get approval and clearState as []byte
	clear, err := CompileSmartContractTeal(algodClient, getTeal("clear"))
	if err != nil {
		log.Fatalf("clearState found error, %s", err)
	}
	app, err := CompileSmartContractTeal(algodClient, getTeal("mng"))
	if err != nil {
		log.Fatalf("approval found error, %s", err)
	}

	mcp := future.AddMethodCallParams{
		AppID:           0,
		Sender:          acct.Address,
		SuggestedParams: txParams,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
		ApprovalProgram: app,
		ClearProgram:    clear,
		GlobalSchema:    types.StateSchema{NumUint: 6, NumByteSlice: 0},
		LocalSchema:     types.StateSchema{NumUint: 0, NumByteSlice: 0},
	}

	var atc future.AtomicTransactionComposer
	err = atc.AddMethodCall(combine(mcp, getMethod(contract, "create"), []interface{}{usdc}))
	if err != nil {
		log.Fatalf("Failed to AddMethodCall: %+v", err)
	}

	debugAppCall(algodClient, atc, "./dryrun/create.msgp", "./dryrun/response/create.json")

	// get the created appID
	acctInfo, err := algodClient.AccountInformation(acct.Address.String()).Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch Account information: %+v", err)
	}
	newApp = acctInfo.CreatedApps[len(acctInfo.CreatedApps)-1].Id

	updateABI(algodClient, newApp, "mng")
	return
}

func CreateASA(algodClient *algod.Client, acct crypto.Account, amt uint64, dec uint32, name, url string) (assetID uint64, err error) {
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}
	addr := acct.Address.String()
	txn, err := future.MakeAssetCreateTxn(addr, []byte(""), txParams, amt, dec, false, addr, addr, addr, addr, name, name, url, "")
	if err != nil {
		log.Fatalf("Failed creating asset: %+v", err)
	}
	signSendWait(algodClient, acct.PrivateKey, txn)
	// get the created assetID
	acctInfo, err := algodClient.AccountInformation(acct.Address.String()).Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch Account information: %+v", err)
	}
	assetID = acctInfo.CreatedAssets[len(acctInfo.CreatedAssets)-1].Index
	return
}

// Start sandbox and create USDCa and other NFTs for testing purpose
func Start(algodClient *algod.Client, acct crypto.Account) (assetID uint64, err error) {
	assetID, err = CreateASA(algodClient, acct, 18446744073709551615, 6, "USDC", "https://circle.com/")
	return
}

func getMethod(c *abi.Contract, name string) (m abi.Method) {
	for _, m = range c.Methods {
		if m.Name == name {
			return
		}
	}
	log.Fatalf("No method named: %s", name)
	return
}

func combine(mcp future.AddMethodCallParams, m abi.Method, a []interface{}) future.AddMethodCallParams {
	mcp.Method = m
	mcp.MethodArgs = a
	return mcp
}

func getContract(setContract string) (c *abi.Contract, err error) {
	f, err := os.Open(getAbi(setContract))
	if err != nil {
		log.Fatalf("Failed to open contract file: %+v", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read file: %+v", err)
	}

	c = &abi.Contract{}
	if err := json.Unmarshal(b, c); err != nil {
		log.Fatalf("Failed to marshal contract: %+v", err)
	}
	return
}

func updateABI(algodClient *algod.Client, newApp uint64, setContract string) {
	contract, err := getContract(setContract)
	if err != nil {
		log.Fatalf("can't get contract: %+v\n", err)
	}

	// update appID of contract
	if network, ok := contract.Networks[net]; ok {
		network.AppID = newApp
		contract.Networks[net] = network
	}

	out, err := json.MarshalIndent(contract, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal: %+v", err)
	}
	err = ioutil.WriteFile(getAbi(setContract), out, 0666)
	if err != nil {
		log.Fatalf("Failed to write file: %+v", err)
	}
}

func CompileSmartContractTeal(algodClient *algod.Client, file string) (compiledProgram []byte, err error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	tealFile, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read file: %s\n", err)
	}
	compileResponse, err := algodClient.TealCompile(tealFile).Do(context.Background())
	if err != nil {
		log.Fatalf("Issue with compile: %s\n", err)
	}
	compiledProgram, _ = base64.StdEncoding.DecodeString(compileResponse.Result)
	log.Printf("%s size: %v\n", file, len(compiledProgram))
	return
}

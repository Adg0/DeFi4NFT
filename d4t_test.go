package d4t

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"testing"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
)

var (
	usdc       = uint64(10458941)
	collateral = uint64(97931298)
	mng        = uint64(84436122)
	lqt        = uint64(84436752)
	d4t        = uint64(84436769)
	dUSD       = uint64(84436770)
	inv        = uint64(84436771)
	address    = AlgodAddressPurestake //"http://localhost:4001"
	token      = AlgodTokenPurestake   //strings.Repeat("a", 64)
	note       = "purestake"           //"local"
	accts      = acctsPurestake        // , _ = GetAccounts()
)

func TestConfigASA(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	err = ConfigASA(algodClient, accts[2], mng, d4t, lqt, collateral)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestStart(t *testing.T) {
	// create USDC asset
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	_, err = Start(algodClient, accts[0])
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestCreateASA(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	_, err = CreateASA(algodClient, accts[2], 1000, 0, "LFT", "https://")
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestOptinASA(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	err = OptinASA(algodClient, accts[1], usdc)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestDeploy(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	mng, err = Deploy(algodClient, accts[0], usdc)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestUpdate(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	err = Update(algodClient, accts[0])
	if err != nil {
		t.Errorf("test found error, %s", err)
	}

}

func TestFund(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	err = Fund(algodClient, accts[0], 2000000)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestCreateApps(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	_, err = CreateApps(algodClient, accts[0], usdc)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}

}

func TestConfigureApps(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	err = ConfigureApps(algodClient, accts[0], lqt, d4t, usdc, dUSD)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}

}

func TestUsdc(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get suggeted params: %+v", err)
	}

	txn, _ := future.MakeAssetTransferTxn(accts[0].Address.String(), accts[2].Address.String(), 100000000, nil, txParams, "", usdc)
	signSendWait(algodClient, accts[0].PrivateKey, txn)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestSendDusd(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	rec := crypto.GetApplicationAddress(d4t)
	err = SendDusd(algodClient, accts[0], rec, dUSD)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestChildUpdate(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	err = ChildUpdate(algodClient, accts[0], d4t, "./teal/d4tApp.teal", "./teal/d4tClear.teal")
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestOptinContract(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	err = Optin(algodClient, accts[1], mng, "d4t")
	if err != nil {
		t.Errorf("test found error, %s", err)
	}
}

func TestEarn(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		t.Errorf("Failed to get suggeted params: %+v", err)
	}

	xids := []uint64{collateral, dUSD, inv}
	aamt := uint64(100000000)
	lvr := uint64(172800) + uint64(txParams.FirstRoundValid)

	lsigArgs := make([][]byte, 4)
	var buf [4][8]byte
	binary.BigEndian.PutUint64(buf[0][:], usdc) // USDCa asset ID
	binary.BigEndian.PutUint64(buf[1][:], aamt) // loan available (50 USDCa)
	binary.BigEndian.PutUint64(buf[2][:], lvr)  // Expiring lifespan: 17280 rounds == 1 day
	binary.BigEndian.PutUint64(buf[3][:], d4t)  // d4t appID
	lsigArgs[0] = buf[0][:]
	lsigArgs[1] = buf[1][:]
	lsigArgs[2] = buf[2][:]
	lsigArgs[3] = buf[3][:]

	lsaRaw := CompileToLsig(algodClient, lsigArgs, getTeal("lsig"), "./codec/lender_lsig.codec", accts[3].PrivateKey)
	if lsaRaw.SigningKey == nil {
		t.Errorf("lsig is empty")
	}
	lsa := sha256.Sum256(lsaRaw.Lsig.Logic)

	err = Earn(algodClient, accts[3], xids, aamt, lvr, lsa[:4])
	if err != nil {
		t.Errorf("test found error, %s", err)
	}

}

func TestClaim(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	amt := uint64(10000000)

	err = Claim(algodClient, accts[1], mng, amt, usdc, dUSD)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}

}

func TestBorrow(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	xids := []uint64{collateral}
	camt := []uint64{1}
	lamt := []uint64{1000000}

	err = Borrow(algodClient, accts[4], accts[3].Address, usdc, dUSD, mng, lqt, xids, camt, lamt, "./codec/lender_lsig.codec")
	if err != nil {
		t.Errorf("test found error, %s", err)
	}

}

func TestRepay(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}

	accts, err := GetAccounts()
	if err != nil {
		log.Fatalf("Failed to get accounts: %+v", err)
	}

	xids := []uint64{collateral}
	ramt := []uint64{10000000}

	err = Repay(algodClient, accts[2], mng, lqt, usdc, xids, ramt)
	if err != nil {
		t.Errorf("test found error, %s", err)
	}

}

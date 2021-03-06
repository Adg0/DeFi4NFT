package d4t

import (
	"testing"
)

func TestFetchLsigFromFile(t *testing.T) {
	result, err := FetchLsigFromFile("./codec/dispenserDUSD.codec")
	t.Logf("Lsig : %#v\n", result)
	if err != nil {
		t.Errorf("expecting no errors, got %s", err)
	}
	if result.Lsig.Logic == nil {
		t.Errorf("wrong lsig, %s", result.Lsig.Logic)
	}
}

func TestDispenseAssetDUSD(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}
	err = DispenseAsset(algodClient, ReserveAddr, ToAddr, 10000000, dUSD, "./codec/dispenserdUSD.codec")
	if err != nil {
		t.Errorf("expecting no errors, got %s", err)
	}
}

func TestDispenseAsset(t *testing.T) {
	algodClient, err := InitAlgodClient()
	if err != nil {
		t.Errorf("algodClient found error, %s", err)
	}
	err = DispenseAsset(algodClient, ReserveAddr, BonusAddr, 4, LFT_jina, "./codec/dispenserLFT.codec")
	t.Errorf("expecting no errors, got %s", err)
}

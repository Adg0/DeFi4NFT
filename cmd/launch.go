package main

import (
	"log"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	d "github.com/kalebBd/DeFi4NFT"
)

var (
	usdc           = uint64(10458941)
	sandboxAddress = "http://localhost:4001"
	sandboxToken   = strings.Repeat("a", 64)
)

var defi4Question = []*survey.Question{
	{
		Name: "defi4",
		Prompt: &survey.Select{
			Message: "Configure DeFi4NFT:",
			Options: []string{"Create DeFi4NFT", "Update DeFi4NFT", "Delete DeFi4NFT"},
			Default: "Create DeFi4NFT",
			Help:    "Select action to execute on DeFi4",
		},
	},
}

var networkQuestion = []*survey.Question{
	{
		Name: "network",
		Prompt: &survey.Select{
			Message: "Choose network:",
			Options: []string{"sandbox", "testnet", "betanet", "mainnet"},
			Default: "testnet",
			Help:    "Select which network you want to run DeFi4",
		},
	},
}

var borrowQuestion = []*survey.Question{
	{
		Name:   "collateral",
		Prompt: &survey.Input{Message: "Which NFT/asset to use as collateral?", Help: "Specify the collateral, enter the AssetID of your collateral."},
	},
	{
		Name:   "camt",
		Prompt: &survey.Input{Message: "How much of your asset are youusing as collateral?", Help: "Specify the quantity of collateral."},
	},
	{
		Name:   "lamt",
		Prompt: &survey.Input{Message: "How much loan do you want to take out?", Help: "Amount that will be deposited to your account."},
	},
	{
		Name:   "lenders",
		Prompt: &survey.Input{Message: "Enter which lenders to borrow from:", Help: "The Algorand address of the lenders."},
	},
}

var earnQuestion = []*survey.Question{
	{
		Name: "allows",
		Prompt: &survey.Select{
			Message: "State NFTs/Assets that can borrow your promise:",
			Options: []string{"NFT-1", "NFT-2", "NFT-3", "NFT-4"},
			Default: "NFT-1",
		},
	},
	{
		Name:   "amount",
		Prompt: &survey.Input{Message: "What is the maximum amount your willing to lend?", Help: "Amount that will be withdrawn from your account."},
	},
	{
		Name:   "expire",
		Prompt: &survey.Input{Message: "For how long is this promise active?", Help: "Expiration date length from now."},
	},
}

var configureNFTQuestion = []*survey.Question{
	{
		Name:   "configure",
		Prompt: &survey.Input{Message: "Which NFT are you enabling to be borrowable in DeFi4?", Help: "AssetID of an NFT you are a manager/creator of."},
	},
}

var claimQuestion = []*survey.Question{
	{
		Name:   "claim",
		Prompt: &survey.Input{Message: "How much are you claiming?", Help: "specify the amount of dUSD your going to swap for USDC."},
	},
}

var repayQuestion = []*survey.Question{
	{
		Name:   "repay",
		Prompt: &survey.Input{Message: "Which NFT are you repaying loan for?", Help: "AssetID you plan on repay debt for."},
	},
	{
		Name:   "amount",
		Prompt: &survey.Input{Message: "How much are you repaying?", Help: "Amount you are repay on the loan."},
	},
}

var accountQuestion = []*survey.Question{
	{
		Name: "account",
		Prompt: &survey.Select{
			Message: "Choose which account you would like to use:",
			Options: []string{"sandbox", "testnet", "betanet", "mainnet"},
			Default: "sandbox",
		},
	},
}

func main() {

	accts_fetched, err := d.GetAccounts()
	if err != nil {
		log.Fatalf("Failed to get accounts: %+v", err)
	}
	var acct_list []string
	for i, at := range accts_fetched {
		acct_list[i] = at.Address.String()
	}

	// Prompt questions
	var firstQuestion = []*survey.Question{
		{
			Name: "request",
			Prompt: &survey.Select{
				Message: "Interact with DeFi4NFT:",
				Options: []string{"borrow", "earn", "repay", "claim", "liquidate", "activate NFT", "change Network", "clear out", "Transfer uncollateralized NFT", "Configure DeFi4NFT"},
				Default: "earn",
				Help:    "Select action to execute on DeFi4",
			},
		},
		{
			Name: "account",
			Prompt: &survey.Select{
				Message: "Choose which account you would like to use:",
				Options: acct_list,
			},
		},
	}

	// First question's answer
	answer1 := struct {
		Request string
		Account string
	}{}

	// perform the questions
	err = survey.Ask(firstQuestion, &answer1, survey.WithHelpInput('?'), survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "?"
		icons.Question.Format = "teal+hb"
		icons.Help.Text = "!"
		icons.Help.Format = "cyan"
	}))
	if err != nil {
		log.Fatalln(err.Error())
	}
	// TODO: Switch to account index according to choice selected
	acct := 0
	if answer1.Account == "testnet" {
		acct = 0
	}

	switch answer1.Request {
	case "Configure DeFi4NFT":
		// another prompt to specify among configurations [moreQuestion]
		ans := struct{ Defi4 string }{}
		err := survey.Ask(defi4Question, &ans, survey.WithHelpInput('?'), survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "?"
			icons.Question.Format = "teal+hb"
			icons.Help.Text = "!"
			icons.Help.Format = "cyan"
		}))
		if err != nil {
			log.Fatalln(err.Error())
		}
		switch ans.Defi4 {
		case "Create DeFi4NFT":
			CreateD4T(acct, 1)
		case "Update DeFi4NFT":
			UpdateD4T()
		case "Delete DeFi4NFT":
			DeleteD4T()
		default:
			CreateD4T(acct, 1)
		}
	case "borrow":
		// a prompt to specify parameters for borrow
		ans := struct {
			Collateral string
			Camt       int
			Lamt       int
			Lenders    string
		}{}
		err := survey.Ask(borrowQuestion, &ans, survey.WithHelpInput('?'), survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "?"
			icons.Question.Format = "teal+hb"
			icons.Help.Text = "!"
			icons.Help.Format = "cyan"
		}))
		if err != nil {
			log.Fatalln(err.Error())
		}
		//TODO: get the id of the unit-name of asset with indexer
		col := uint64(1)
		if ans.Collateral == "NFT-1" {
			col = 1
		}
		BorrowD4T(col, uint64(ans.Camt), uint64(ans.Lamt), ans.Lenders)
	case "earn":
		// a prompt to specify parameters for earn
		ans := struct {
			Allows []uint64
			Amount int
			Expire int
		}{}
		err := survey.Ask(earnQuestion, &ans, survey.WithHelpInput('?'), survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "?"
			icons.Question.Format = "teal+hb"
			icons.Help.Text = "!"
			icons.Help.Format = "cyan"
		}))
		if err != nil {
			log.Fatalln(err.Error())
		}
		EarnD4T()
	case "claim":
		// a prompt to specify parameters for earn
		ClaimD4T()
	case "repay":
		// a prompt to specify parameters for earn
		RepayD4T()
	case "collateralChange":
		// a prompt to specify parameters for earn
		ChangeCollateralD4T()
	case "liquidate":
		// a prompt to specify parameters for liquidating
		Liquidate()
	case "Transfer uncollateralized NFT":
		// a prompt to specify parameters for sending frozen NFT
		SendFrozenNFT()
	case "activate NFT":
		// a prompt to specify parameters for earn
		ConfigureNFT()
	case "clear out":
		// a prompt to specify parameters for earn
		CloseOut()
	default:
		log.Println("no selection detected")
	}
}

func CreateD4T(acctIndex int, anotherIndex int) {
	log.Println("Inside CreateD4T")
	return
	algodClient, err := d.InitAlgodClient()
	if err != nil {
		log.Fatalf("algodClient found error: %s", err)
	}

	accts, err := d.GetAccounts()
	if err != nil {
		log.Fatalf("Failed to get accounts: %+v", err)
	}
	// Create USDC asset for sandbox
	usdc, err = d.Start(algodClient, accts[acctIndex])
	if err != nil {
		log.Fatalf("Start found error: %s", err)
	}

	// Deploy manager contract
	mng, err := d.Deploy(algodClient, accts[acctIndex], usdc)
	if err != nil {
		log.Fatalf("Deploying found error: %s", err)
	}
	err = d.Fund(algodClient, accts[acctIndex], 2000000)
	if err != nil {
		log.Fatalf("Funding contract found error: %s", err)
	}

	ids, err := d.CreateApps(algodClient, accts[acctIndex], usdc)
	lqt := ids[0]
	d4t := ids[1]
	dusd := ids[2]
	inv := ids[3]

	err = d.ConfigureApps(algodClient, accts[acctIndex], lqt, d4t, usdc, dusd)
	if err != nil {
		log.Fatalf("Configuring created apps found error: %s", err)
	}

	// Create NFT for testing dapp
	collateral, err := d.CreateASA(algodClient, accts[anotherIndex], 1000, 0, "LFT", "https://")
	if err != nil {
		log.Fatalf("Create NFT found error: %s", err)
	}
	err = d.ConfigASA(algodClient, accts[anotherIndex], mng, d4t, lqt, collateral)
	if err != nil {
		log.Fatalf("Configuring NFT found error: %s", err)
	}
	log.Printf("D4T dapp is launched!\n\tManager contract:%d\n\tLiquidator contract:%d\n\tD4T contract:%d\n", mng, lqt, d4t)
	log.Printf("Asset IDs\n\tusdc assetID:%d\n\tdusd assetID:%d\n\tinv assetID:%d\n\tcollateral demo assetID:%d\n", usdc, dusd, inv, collateral)
}

func EarnD4T() {
	log.Println("Inside EarnD4T")
	return
}

func BorrowD4T(col uint64, camt, lamt uint64, lender string) {
	log.Println("Inside BorrowD4T")
	return
}

func RepayD4T() {
	log.Println("Inside RepayD4T")
	return
}

func ClaimD4T() {
	log.Println("Inside ClaimD4T")
	return
}

func ChangeCollateralD4T() {
	log.Println("Inside ChangeCollateralD4T")
	return
}

func Liquidate() {
	log.Println("Inside Liquidate")
	return
}

func SendFrozenNFT() {
	log.Println("Inside SendFrozenNFT")
	return
}

func ConfigureNFT() {
	log.Println("Inside ConfigureNFT")
	return
}

func UpdateD4T() {
	log.Println("Inside UpdateD4T")
	return
}

func DeleteD4T() {
	log.Println("Inside DeleteD4T")
	return
}

func CloseOut() {
	log.Println("Inside CloseOut")
	return
}

# Technical review

* NFTs used as collateral are frozen in account, only when account takes out loan.
* Frozen NFTs are unfrozen when full loan is paid back.
* There is no interest rate for borrowing USDCa.
* A 3% fee is paid to take out loan.
* Lenders sign a delegated logic signature to allow any account to withdraw USDCa that fullfill the following:
	1. Calls DeFi4NFT contract
	2. Withdraws atmost staked amount
* Any account that holds dUSD can claim 1:1 USDCa by sending the dUSD to DeFi4NFT contract.
* Borrower can borrow from upto 4 lenders
* Liquidation

How liquidation happens?

```{mermaid}
sequenceDiagram
  participant Liquidator
  participant DeFi4NFT
  Liquidator->DeFi4NFT: Liquidate borrower
  loop Check for loan health
      DeFi4NFT->DeFi4NFT: Return loan health
  end
  Note right of DeFi4NFT: If loan is health <br>just reject liquidation request.
  Borrower-->Liquidator: Clawedback NFT
  Liquidator->DeFi4NFT: Sends owed debt
  Borrower-->DeFi4NFT: Unfreeze account for remaning collateral asset
```

* Specify the addresss to liquidate
* Pay 95% of collateral's value to DeFi4NFT contract
* Set an account that will receive the liquidated asset
* You'll be sent the collateral to the address you specified

## Smartcontract

There are three smartcontracts that power DeFi4NFT dapp.

1. DeFi4NFT Contract
DeFi4NFT contract holds the state machine and locks/unlocks NFT in account (freezes/unfreezes  NFT) .
State machine, tracks:
	* `xids` tracks which NFT is used as collateral
	* `camt` tracks how much collateral is used for loan
	* `lamt` tracks how much loan is borrowed
	* `aamt` tracks how much loan is available from lender address

2. Liquidator Contract
Liquidator contract reads current price of NFT from oracle and if loan is more than 90% of collateral it liquidates the NFT locked.
	* liquidator contract is the clawback address of leveragable NFTs on DeFi4NFT.
	* after liquidation completes the remainig asset is unfrozen. This is possible by AVM 1.1 (contract to contract call). Liquidator contract calls DeFi4NFT contract to unfreeze the asset.

3. Manager Contract
Manager contract creates all other contracts on behalf of creator address. It also controls the NFTs that are configured to be borrowable in DeFi4NFT.

```{admonition} Notice
:class: info

We have implemented here the Freeze admin and Clawback of algorand ASA. And IPFS for storing delegated LogicSig. We have also used Circle APIs for blockchain swap.
```

```{mermaid}
sequenceDiagram
  participant Lender X
  participant Borrower
  loop Search for lenders
      Borrower->Borrower: Found Lenders
  end
  Note right of Borrower: If no lenders, <br/>transaction fails
  Lender X->Borrower: Receive promised loan
  DeFi4NFT-->Lender X: Decrement maximum allowed amount
  DeFi4NFT-->Borrower: Freeze collateral NFT
```

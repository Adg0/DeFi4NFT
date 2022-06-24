# DeFi for NFT (DeFi4NFT)
[![Documentation Status](https://readthedocs.org/projects/defi4nft/badge/?version=latest)](https://defi4nft.readthedocs.io/en/latest/?badge=latest)

[![d4t logo](/docs/assets/images/logo_black.png)](https://youtu.be/4n19YhPuku4 "DeFi 4 NFT demo video")

## Frontend Repository

[Link to frontend source](https://github.com/adapole/Defi-for-NFT)

[Documentation](https://defi4nft.readthedocs.io/en/latest)

## About DeFi4NFT

DeFi4NFT dapp is an NFT collateralizing app on Algorand. Powered by circle digital money, which let's you use any NFT you hold in Algorand as leverage for borrowing stable coin in DeFi4NFT.
Leveraged NFT remains locked in your (*borrower*) address until the full loan amount is repaid, which is a **pure non-custodial** transaction.
Assets(NFT) that require KYC or special permissions might benefit from this *pure non-custodial protocol*.

With a 3% return for liquidity providers that supply liquidity to the protocol, DeFi4NFT becomes the go-to site for delivering liquidity for NFTs.
Liquidity providers will keep their liquidity asset in their account until a borrower uses it, at which point the protocol will reward them with an I-O-U stable-coin token plus 3% of the lent amount.

## Using DeFi4NFT dapp

First step is to optin to the smartcontract.

## Optin to DeFi4NFT

### As an NFT creator

Transfer your NFT's admin address to DeFi4NFT.
This will make your NFT leverageable for taking loan in DeFi4NFT dapp.
* This sets manager and freeze admin address to DeFi4NFT smartcontract
* And sets clawback to liquidator smartcontract

### As a liquidity provider

Optin to the I-O-U token of DeFi4NFT dapp **dUSD**, that has 1:1 value with USDCa.
This happens automatically when you create a promise to provide liquidity, via the frontend.

## Earn (Providing Liquidity)

Choose which NFTs can borrow from your account.
* Set maximum amount you are willing to lend.
* Set expiration date for aggrement.

[![Providing liquidity](/docs/assets/images/lend.png)](https://youtu.be/4n19YhPuku4?t=51 "Stake your USDC")

## Borrow (Leveraging NFT)

Use your NFT as collateral, to borrow USDCa stablecoin.
* Set which NFT you want to collateralize
* Set amount of collateral
* Request loan
You'll get requested loan amount in USDCa and your NFT will be locked.

[![Leverage NFT](/docs/assets/images/borrow.png)](https://youtu.be/4n19YhPuku4?t=114 "Borrow in DeFi4NFT")

## Repaying loan

Send USDCa to DeFi4NFT contract.
Your loan amount state will be decremented by sent repaid amount.

[![Repay loan](/docs/assets/images/repay.png)](https://youtu.be/4n19YhPuku4?t=149 "Repaying loan")

If you pay the full loan amount, your collateral assets will be unfrozen.

## Claming USDCa

Send dUSD(I-O-U token of DeFi4NFT contract) to DeFi4NFT contract.
You'll receive a 1:1 USDCa for the dUSD you send.

[![Claim USDCa](/docs/assets/images/claim.png)](https://youtu.be/4n19YhPuku4?t=95 "Claim")

## Circle (Fiat-on ramp)

A fast and easy way to get fiat into our Dapp is using Circle accounts.
You can either use a credit/debit card or send USDC from other chains supported by circle.

### Swap USDC to Algorand

Powered by circle's bridge, we now offer wider options to users that want to interact in our pure non-custodial borrow/lend dapp.
You can transfer USDC from any chain supported by circle bridge to any algorand address you want to.

* First select your circle wallet linked to an algorand address.

[![Select wallet](/docs/assets/images/select_wallet_swap.png)](https://youtu.be/4n19YhPuku4?t=95 "Select wallet")

* Then generate a blockchain address that you'll deposite USD into and get it in the Algorand address linked to that circle wallet.

[![Swap USDCa](/docs/assets/images/generate_address.png)](https://youtu.be/4n19YhPuku4?t=95 "Generate Address")

* Next send USD to the generated addresses

[![Deposit address USDCa](/docs/assets/images/get_address.png)](https://youtu.be/4n19YhPuku4?t=95 "Get Address")

You will recieve the deposited amount in your Algorand address soon.
For testing purposes we recommend using smaller amounts.

### Card deposit

To create a card deposite into your algorand address, follow this steps:

1. Create a circle wallet

[![Create wallet](/docs/assets/images/create_wallet.png)](https://youtu.be/4n19YhPuku4?t=95 "Create wallet")

2. Name your wallet

[![Create wallet](/docs/assets/images/create_wallet2.png)](https://youtu.be/4n19YhPuku4?t=95 "Create wallet")

3. Select the newly created wallet or any wallet you prefer

[![Select wallet](/docs/assets/images/select_wallet.png)](https://youtu.be/4n19YhPuku4?t=95 "Select wallet")

4. Then add a credit/debit card information. It is encrypted on client side, and secured as per the standards

[![Add card](/docs/assets/images/add_card.png)](https://youtu.be/4n19YhPuku4?t=95 "Add card")

5. Add billing information for your card.

[![Add card 1](/docs/assets/images/add_card2.png)](https://youtu.be/4n19YhPuku4?t=95 "Billing Address")

6. Select your card from list of cards.

[![Select card](/docs/assets/images/select_card.png)](https://youtu.be/4n19YhPuku4?t=95 "Select card")

7. Fill in the amount you want to charge into your algorand address, and for fraud protection fill in the `cvv`. For testing all cards use `123`.

[![Make payment](/docs/assets/images/make_payment.png)](https://youtu.be/4n19YhPuku4?t=95 "Make payment")

Then you'll get a deposit arrive in your algorand address linked to the wallet you selected within about 10 to 30 minutes.

For more, look into Techincal review.

## Future

We plan on releasing this project on mainnet.

## Techincal info

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

How liquidation works?

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


# Contact

Discord @1egen#0803

Discord @3spear#9556

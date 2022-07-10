Welcome to DeFi4NFT!
=================

Guide
^^^^^
DeFi4NFT is a borrow/lend platform for NFTs on Algorand. At the core a state machine stores the loan health of a pure non-custodial borrow/lend protocol.

**Pure non-custodial** means assets remain in the owner’s wallet, for borrowers this means collateral assets remain *frozen in address*. And for liquidity providers this means creating a *delegated logicSig* promising to provide a loan when a borrower matches, aka promise. And for liquidators this means a three way transaction, where the end receiver is a third party buyer willing to buy the borrower’s collateral NFT and the liquidator pays the debt of the borrower while a third party buyer sends payment of the collateral NFT to liquidator. This system is set up so permissioned tokens can be open to receive liquidity from non-whitelisted addresses.

.. toctree::
   :maxdepth: 2

   getting_started.md
   neos.md
   technical_review.md
   install.md
   troubleshooting.md
   more.md
   contact.md


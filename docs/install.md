# Install or build project

## Building DeFi4NFT Locally

## Create a sandbox environment
```
$ git clone https://github.com/Adg0/DeFi4NFT.git
$ cd DeFi4NFT
$ git clone https://github.com/algorand/sandbox.git
$ cd sandbox
$ ./sandbox up dev
```
```{admonition} Notice
:class: info

Visit setting up [algorand sandbox](https://github.com/algorand/sandbox#algorand-sandbox) for development.
```

## Create the contracts locally
```
$ cd cmd
$ go install .
$ ./cmd
// scroll down to bottom, where it says Configure DeFi4NFT
>> Interact with DeFi4NFT: 
	Configure DeFi4NFT
// select any
>> Choose which account to you would like to use:
	<algorand address>
// select Create DeFi4NFT
>> Configure DeFi4NFT:
	Create DeFi4NFT
```


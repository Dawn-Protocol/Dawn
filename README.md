
# Welcome to the New Begining

<img src="https://i.postimg.cc/4dvwcKJ9/logo.png" width="340">

#### Dawn is a global, self-governed, and open-sourced protocol that empowers and rewards gamers.</br>
Built with <a href = "https://tendermint.com/">Tendermint </a> and <a href = "https://cosmos.network/sdk">Cosmos SDK</a>

This repo is for Dawn's testnet build.


## Requirements
 - Go 1.13

## DEPLOYMENT OF SMART CONTRACT

### Set-up

```bash
Into the Project Directory=>

cd testnet-contracts/

# Create .env with sample environment variables
cp .env.example .env

# Add environment variable MNEMONIC from your MetaMask account in .env file.

# Add environment variable INFURA_PROJECT_ID from your Infura account in .env file.

```

### Running the bridge on the Ropsten testnet

cd testnet-contracts/

run the following commands..
```bash

# Deploy contract to ropsten network
yarn migrate --network ropsten

# Get contract's address
yarn peggy:address --network ropsten

```

## SET UP OF NODES/VALIDATORS

### ON NODE-1

### Terminal-1

```bash

# Install tools (golangci-lint v1.18)
make tools-clean
make tools

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands, confirming the build is successful:
dawnd help
dawncli help
dawnrelayer help
```

## Running and testing the application

First, initialize a chain and create accounts.

```bash
# Initialize the genesis.json file that will help you to bootstrap the network
dawnd init node1 --chain-id=dawn-protocol

# Create a key to hold your validator account and for another test account
dawncli keys add validator1
# Enter password

dawncli keys add testuser
# Enter password

# Edit the genesis.json file for customised stake denom
--> on terminal go to folder ~/.dawnd/config/genesis.json
--> Edit the staking section having bond_denom key from "stake" to "dawn"
--> Save and close the file.


# Initialize the genesis account and transaction
dawnd add-genesis-account $(dawncli keys show validator1 -a) 1000000000dawn,1000000000vspr

# Create genesis transaction
dawnd gentx --name validator1 --amount 1000000dawn
# Enter password

# Collect genesis transaction
dawnd collect-gentxs

#Add Customisation to genesis file
--> on terminal go to folder ~/.dawnd/config/genesis.json
Replace '"staking"' section with the json mentioned in customised_genesis.json (which is in Project Directory)

# Now its safe to start `dawnd`
dawnd start

```
## Terminal-2:Start the Relayer service

For automated relaying, there is a relayer service that can be run that will automatically watch and relay events (local web socket and deployed address parameters may vary).

```bash
# Check dawnrelayer connection to ebd
dawnrelayer status

# Start dawnrelayer on the contract's deployed address with [PEGGY_DEPLOYED_ADDRESS]

dawnrelayer init wss://ropsten.infura.io/ws [PEGGY_DEPLOYED_ADDRESS] LogLock\(bytes32,address,bytes,address,string,uint256,uint256\) validator1 --chain-id=dawn-protocol
# Enter password and press enter
# You should see a message like: Started ethereum websocket with provider: wss://ropsten.infura.io/ws \ Subscribed to contract events on address: [PEGGY_DEPLOYED_ADDRESS]
# The relayer will now watch the contract on Ropsten and create a claim whenever it detects a lock event.

#Using the application from rest-server

dawncli rest-server --trust-node

```


### ON NODE-2

### Terminal-1

```bash

# Install tools (golangci-lint v1.18)
make tools-clean
make tools

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands, confirming the build is successful:
dawnd help
dawncli help
dawnrelayer help
```

## Running and testing the application

```bash
# Initialize the genesis.json file with another moniker and same namechain
dawnd init node2 --chain-id=dawn-protocol

# Create a key to hold your validator account and for another test account
dawncli keys add validator2
# Enter password

overwrite ~/.dawnd/config/genesis.json with first nodes genesis.json

#change persistent_peers
#run `dawncli status` on first node to get id.
go to ~/.dawnd/config/config.toml
persistent_peers = "<id_of_node1>@<ip>:<26656>"


# Now its safe to start `dawnd`
dawnd start
```

### FROM NODE-1

```bash

# Then, wait 10 seconds and from first node send tokens to valdator2 address for testing
dawncli tx send validator1 <validator2> 10000dawn --chain-id=dawn-protocol --yes
run "dawncli keys show validator2" on second node to get validator2 address.

```
### FROM NODE-2

## Terminal-2

```bash

# Next, setup the staking module prerequisites
# First, create a validator and stake
dawncli tx staking create-validator \
  --amount=1000dawn \
  --pubkey=$(dawnd tendermint show-validator) \
  --moniker="node2" \
  --chain-id=dawn-protocol \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas=200000 \
  --gas-prices="0.001dawn" \
  --from=validator2

```

## Terminal-3:Start the Relayer service

For automated relaying, there is a relayer service that can be run that will automatically watch and relay events (local web socket and deployed address parameters may vary).

```bash
# Check dawnrelayer connection to ebd
dawnrelayer status

# Start dawnrelayer on the contract's deployed address with [PEGGY_DEPLOYED_ADDRESS]

dawnrelayer init wss://ropsten.infura.io/ws [PEGGY_DEPLOYED_ADDRESS] LogLock\(bytes32,address,bytes,address,string,uint256,uint256\) validator2 --chain-id=dawn-protocol
# Enter password and press enter
# You should see a message like: Started ethereum websocket with provider: wss://ropsten.infura.io/ws \ Subscribed to contract events on address: [PEGGY_DEPLOYED_ADDRESS]
# The relayer will now watch the contract on Ropsten and create a claim whenever it detects a lock event.

#Using the application from rest-server

dawncli rest-server --trust-node

```

### ON NODE-3

### Terminal-1

```bash

# Install tools (golangci-lint v1.18)
make tools-clean
make tools

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands, confirming the build is successful:
dawnd help
dawncli help
dawnrelayer help
```

## Running and testing the application

First, initialize a chain and create accounts.

```bash
# Initialize the genesis.json file with another moniker and same namechain
dawnd init node3 --chain-id=dawn-protocol

# Create a key to hold your validator account and for another test account
dawncli keys add validator3
# Enter password

overwrite ~/.dawnd/config/genesis.json with first nodes genesis.json

#change persistent_peers

go to ~/.dawnd/config/config.toml
persistent_peers = "id@first_node_ip:26656, id@second_node_ip:26659"
run "dawncli status" on first node to get first_node id.
run "dawncli status" on second node to get second_node id.

# Now its safe to start `dawnd`
dawnd start

```
### FROM NODE-1

```bash

# Then, wait 10 seconds and from first node send tokens to valdator3 address for testing
dawncli tx send validator1 <validator3> 10000dawn --chain-id=dawn-protocol --yes
run "dawncli keys show validator3" on third node to get validator3 address.

```
### FROM NODE-3

## Terminal-2

```bash

# Next, setup the staking module prerequisites
# First, create a validator and stake
dawncli tx staking create-validator \
  --amount=900dawn \
  --pubkey=$(dawnd tendermint show-validator) \
  --moniker="node3" \
  --chain-id=dawn-protocol \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas=200000 \
  --gas-prices="0.001dawn" \
  --from=validator3

```
## Terminal-3:Start the Relayer service

For automated relaying, there is a relayer service that can be run that will automatically watch and relay events (local web socket and deployed address parameters may vary).

```bash
# Check dawnrelayer connection to ebd
dawnrelayer status

# Start dawnrelayer on the contract's deployed address with [PEGGY_DEPLOYED_ADDRESS]

dawnrelayer init wss://ropsten.infura.io/ws [PEGGY_DEPLOYED_ADDRESS] LogLock\(bytes32,address,bytes,address,string,uint256,uint256\) validator3 --chain-id=dawn-protocol
# Enter password and press enter
# You should see a message like: Started ethereum websocket with provider: wss://ropsten.infura.io/ws \ Subscribed to contract events on address: [PEGGY_DEPLOYED_ADDRESS]
# The relayer will now watch the contract on Ropsten and create a claim whenever it detects a lock event.

#Using the application from rest-server

dawncli rest-server --trust-node

```


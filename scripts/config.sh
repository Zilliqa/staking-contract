#!/usr/bin/env bash

########################################################################################
# Update the values below
STAKING_PRIVKEY_ADMIN="38f3715e7ef9b5a5080171dca4cb37b05eaa7e3b0d9a9427a11e021e1029525d"
STAKING_PRIVKEY_VERIF="38f3715e7ef9b5a5080171dca4cb37b05eaa7e3b0d9a9427a11e021e1029525d"
STAKING_ADDRESS_ADMIN="09710e00256db2e3db4b44f597f17f3d97f06318"
STAKING_ADDRESS_PROXY="09710e00256db2e3db4b44f597f17f3d97f06318"
STAKING_ADDRESS_SSNLIST="09710e00256db2e3db4b44f597f17f3d97f06318"
STAKING_ADDRESS_VERIF="09710e00256db2e3db4b44f597f17f3d97f06318"
TESTNET_API_URL="https://dev-api.zilliqa.com/"
TESTNET_CHAINID="2"
########################################################################################

# Replace testnet settings in go-zli wallet file
sed -i "s|\"api\":[^,]*,|\"api\":\"${TESTNET_API_URL}\",|g" ~/.zilliqa
sed -i "s|\"chain_id\":[^,]*,|\"chain_id\":${TESTNET_CHAINID},|g" ~/.zilliqa

# Replace admin address in proxy.json file for deploy_proxy.sh
sed -i "8s|.*|  \"value\": \"0x${STAKING_ADDRESS_ADMIN}\"|" proxy.json

# Replace admin address and proxy address in ssnlist.json file for deploy_ssnlist.sh
sed -i "8s|.*|  \"value\": \"0x${STAKING_ADDRESS_ADMIN}\"|" ssnlist.json
sed -i "12s|.*|  \"value\": \"0x${STAKING_ADDRESS_PROXY}\"|" ssnlist.json
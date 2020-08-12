#!/usr/bin/env bash
source ../config.sh

# Update the values below
FUNDS_IN_QA="10000000000000"

zli contract call -a ${STAKING_ADDRESS_PROXY} -t update_maxstake -k ${STAKING_PRIVKEY_ADMIN} -r "[{\"vname\":\"max_stake\",\"type\":\"Uint128\",\"value\":\"${FUNDS_IN_QA}\"}]" -f true

#!/usr/bin/env bash
source config.sh

# Update the values below
FUNDS_IN_QA="10000000000"

zli contract call -a ${STAKING_ADDRESS_PROXY} -t update_minstake -r "[{\"vname\":\"min_stake\",\"type\":\"Uint128\",\"value\":\"${FUNDS_IN_QA}\"}]" -f true

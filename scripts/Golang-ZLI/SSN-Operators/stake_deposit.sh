#!/usr/bin/env bash
source ../config.sh

# Update the values below
STAKING_PRIVKEY_SSN="38f3715e7ef9b5a5080171dca4cb37b05eaa7e3b0d9a9427a11e021e1029525d"
FUNDS_IN_QA="100000000000000"

zli contract call -a ${STAKING_ADDRESS_PROXY} -t stake_deposit -k ${STAKING_PRIVKEY_SSN} -r "[]" -m ${FUNDS_IN_QA} -f true
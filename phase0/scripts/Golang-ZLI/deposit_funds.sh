#!/usr/bin/env bash
source config.sh

# Update the values below
FUNDS_IN_QA="100000000000000"

zli contract call -a ${STAKING_ADDRESS_PROXY} -t deposit_funds -k ${STAKING_PRIVKEY_ADMIN} -r "[]" -m ${FUNDS_IN_QA} -f true
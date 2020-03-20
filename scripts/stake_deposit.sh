#!/usr/bin/env bash
source config.sh

# Update the values below
FUNDS_IN_QA="100000000000000"

zli contract call -a ${STAKING_ADDRESS_PROXY} -t stake_deposit -r "[]" -m ${FUNDS_IN_QA} -f true
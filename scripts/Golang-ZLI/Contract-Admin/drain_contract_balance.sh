#!/usr/bin/env bash
zli contract call -a ${STAKING_ADDRESS_PROXY} -t drain_contract_balance -r "[]" -s ${KEY_STORE_PATH} -f true 

#!/usr/bin/env bash
zli contract call -a ${STAKING_ADDRESS_PROXY} -s ${KEY_STORE_PATH} -t unpause -r "[]" -f true

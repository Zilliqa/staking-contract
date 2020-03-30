#!/usr/bin/env bash
zli contract call -a ${MULTISIG_WALLET_ADDRESS} -t ExecuteTransaction -k ${OWNER_KEY} -r "[{\"vname\":\"transactionId\",\"type\":\"Uint32\",\"value\":\"0\"}]" -f true
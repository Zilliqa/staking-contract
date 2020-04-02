#!/usr/bin/env bash
zli contract call -a ${MULTISIG_WALLET_ADDRESS} -t SignTransaction -k ${OWNER_KEY} -r "[{\"vname\":\"transactionId\",\"type\":\"Uint32\",\"value\":\"${TRANSACTION_ID}\"}]" -f true
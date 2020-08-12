#!/usr/bin/env bash
zli contract call -a ${MULTISIG_WALLET_ADDRESS} -t SubmitCustomUpgradeToTransaction -s ${KEY_STORE_PATH} -r "[{\"vname\":\"proxyContract\",\"type\":\"ByStr20\",\"value\":\"${PROXY_CONTRACT}\"},{\"vname\":\"newImplementation\",\"type\":\"ByStr20\",\"value\":\"${NEW_IMPL}\"}]" -f true

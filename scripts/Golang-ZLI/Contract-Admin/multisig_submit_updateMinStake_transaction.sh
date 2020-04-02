#!/usr/bin/env bash
zli contract call -a ${MULTISIG_WALLET_ADDRESS} -t SubmitUpdateMinStakeTransaction -k ${OWNER_KEY} -r "[{\"vname\":\"proxyContract\",\"type\":\"ByStr20\",\"value\":\"${PROXY_CONTRACT}\"},{\"vname\":\"min_stake\",\"type\":\"Uint128\",\"value\":\"${MIN_STAKE}\"]" -f true

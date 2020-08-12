#!/usr/bin/env bash
zli contract call -a ${STAKING_ADDRESS_PROXY} -s ${KEY_STORE_PATH} -t update_admin -r "[{\"vname\":\"admin\",\"type\":\"ByStr20\",\"value\":\"0xb2e51878722d8b6d2c0f97e995a7276d64c1618b\"}]" -f true

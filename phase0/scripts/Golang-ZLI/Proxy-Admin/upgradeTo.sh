#!/usr/bin/env bash
source ../config.sh

zli contract call -a ${STAKING_ADDRESS_PROXY} -t upgradeTo -s ${KEY_STORE_PATH} -r "[{\"vname\":\"newImplementation\",\"type\":\"ByStr20\",\"value\":\"0x${STAKING_ADDRESS_SSNLIST}\"}]" -f true
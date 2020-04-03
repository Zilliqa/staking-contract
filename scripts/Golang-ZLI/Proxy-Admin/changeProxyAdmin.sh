#!/usr/bin/env bash
source ../config.sh

zli contract call -a ${STAKING_ADDRESS_PROXY} -t changeProxyAdmin -r "[{\"vname\":\"newAdmin\",\"type\":\"ByStr20\",\"value\":\"0x${STAKING_PRIVKEY_ADMIN}\"}]" -f true

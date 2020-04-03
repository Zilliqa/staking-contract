#!/usr/bin/env bash
source ../config.sh

zli contract call -a ${STAKING_ADDRESS_PROXY} -t update_verifier -k ${STAKING_PRIVKEY_ADMIN} -r "[{\"vname\":\"verif\",\"type\":\"ByStr20\",\"value\":\"0x${STAKING_ADDRESS_VERIF}\"}]" -f true

#!/usr/bin/env bash
source config.sh

zli contract call -a ${STAKING_ADDRESS_PROXY} -t update_verifier -r "[{\"vname\":\"verif\",\"type\":\"ByStr20\",\"value\":\"0x${STAKING_ADDRESS_VERIF}\"}]" -f true

#!/usr/bin/env bash
source ../config.sh

# Update the values below
SSN_ADDR="0xb2e51878722d8b6d2c0f97e995a7276d64c1618b"

zli contract call -a ${STAKING_ADDRESS_PROXY} -t remove_ssn -r "[{\"vname\":\"ssnaddr\",\"type\":\"ByStr20\",\"value\":\"${SSN_ADDR}\"}]" -f true

#!/usr/bin/env bash
source ../config.sh

# Update the values below
SSN_ADDR="0xb2e51878722d8b6d2c0f97e995a7276d64c1618b"
STAKE_AMOUNT="0"
REWARDS="0"
URLRAW="https://dev-api.zilliqa.com/"
URLAPI="https://dev-api.zilliqa.com/"
BUFFERED_DEPOSIT="0"

zli contract call -a ${STAKING_ADDRESS_PROXY} -t add_ssn -k ${STAKING_PRIVKEY_ADMIN} -r "[{\"vname\":\"ssnaddr\",\"type\":\"ByStr20\",\"value\":\"${SSN_ADDR}\"},{\"vname\":\"stake_amount\",\"type\":\"Uint128\",\"value\":\"${STAKE_AMOUNT}\"},{\"vname\":\"rewards\",\"type\":\"Uint128\",\"value\":\"${REWARDS}\"},{\"vname\":\"urlraw\",\"type\":\"String\",\"value\":\"${URLRAW}\"},{\"vname\":\"urlapi\",\"type\":\"String\",\"value\":\"${URLAPI}\"},{\"vname\":\"buffered_deposit\",\"type\":\"Uint128\",\"value\":\"${BUFFERED_DEPOSIT}\"}]" -f true

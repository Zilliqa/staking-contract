#!/usr/bin/env bash
source config.sh

zli contract deploy -c ../../contracts/ssnlist.scilla -i ssnlist.json -k ${STAKING_PRIVKEY_ADMIN}
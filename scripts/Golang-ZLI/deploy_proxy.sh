#!/usr/bin/env bash
source config.sh

zli contract deploy -c ../../contracts/proxy.scilla -i proxy.json -k ${STAKING_PRIVKEY_ADMIN}
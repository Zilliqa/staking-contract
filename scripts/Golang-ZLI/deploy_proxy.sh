#!/usr/bin/env bash
source config.sh

zli contract deploy -c ../../contracts/proxy.scilla -i proxy.json -s ${KEY_STORE_PATH}
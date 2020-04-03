#!/usr/bin/env bash
zli contract call -a 09710e00256db2e3db4b44f597f17f3d97f06318 -t withdraw_stake_amount -r "[{\"vname\":\"amount\",\"type\":\"Uint128\",\"value\":\"500000000000\"}]"  -f true

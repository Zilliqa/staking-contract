#!/usr/bin/env bash
zli contract call -a 09710e00256db2e3db4b44f597f17f3d97f06318 -t update_contractmaxstake -r "[{\"vname\":\"max_stake\",\"type\":\"Uint128\",\"value\":\"10000000000000000\"}]" -f true

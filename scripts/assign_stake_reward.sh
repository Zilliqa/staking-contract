#!/usr/bin/env bash
zli contract call -a 09710e00256db2e3db4b44f597f17f3d97f06318 -t assign_stake_reward -r "[{\"vname\":\"ssnreward_list\",\"type\":\"List SsnRewardShare\",\"value\":[{\"constructor\":\"SsnRewardShare\",\"argtypes\":[],\"arguments\":[\"0xb2e51878722d8b6d2c0f97e995a7276d64c1618b\",\"50000000\"]}]},{\"vname\":\"reward_blocknum\",\"type\":\"Uint32\",\"value\":\"50000\"}]" -f true
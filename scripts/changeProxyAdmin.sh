#!/usr/bin/env bash
zli contract call -a 09710e00256db2e3db4b44f597f17f3d97f06318 -t changeProxyAdmin -r "[{\"vname\":\"newAdmin\",\"type\":\"ByStr20\",\"value\":\"0xb2e51878722d8b6d2c0f97e995a7276d64c1618b\"}]" -f true

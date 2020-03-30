#!/usr/bin/env bash
OWNER_KEY=40a08154418fcc0026e9f93f6ed16c6c6a499cbcda1335b581084f18105d1c7b
zli contract deploy -c ../contracts/multisig_wallet.scilla -i wallet.json -k ${OWNER_KEY}
#!/usr/bin/env bash
MULTISIG_WALLET_ADDRESS=78a8d5b16c9ddf0a00d0640e83ad07239c1e9acc
OWNER_KEY=40a08154418fcc0026e9f93f6ed16c6c6a499cbcda1335b581084f18105d1c7b
zli contract call -a ${MULTISIG_WALLET_ADDRESS} -t SubmitDrainContractBalanceTransaction -k ${OWNER_KEY} -r "[{\"vname\":\"proxyContract\",\"type\":\"ByStr20\",\"value\":\"0xe82d1b7f8fdd879ba2709fcd98c6491c84add3f9\"}]" -f true
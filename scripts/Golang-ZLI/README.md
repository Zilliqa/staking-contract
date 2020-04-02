### Contents
The scripts are divided into folders based on the following roles
- Contract admin
- Verifier
- Staked seed node operator
- Proxy contract admin


### 1. Ensure you have [zli](https://github.com/Zilliqa/zli) installed

can use following command to test:

```shell script
zli -h
```

### 2. Init your wallet config

#### 1) Init a new fresh one

```shell script
zli wallet init
```

Check your wallet config:

```shell script
zli wallet echo
```

Transfer some zils your account before making any transaction.

#### 2) Init for an exist private key

If you already have one with zils, then can use that private key to init your wallet config:

```shell script
zli wallet from -p your_private_key
```

Above two commands will generate a file located at ~/.zilliqa, try to edit it directly if you want to change api url or chain
id.

<b>All following commands will use the private key you generated above, if you want to override it, just use `-k another_private_key`</b>

like `zli contract deploy -c ../contracts/proxy.scilla -i proxy.json -k 38f3715e7ef9b5a5080171dca4cb37b05eaa7e3b0d9a9427a11e021e1029525d`.

But we still need to make `~/.zilliqa` exist, because we need `api url` and `chain id`.

### 3. Deploy proxy contract

Run `sh ./deploy_proxy.sh`

Make sure `proxy.json` contain an correct `init_admin`, usually from your config wallet, can use `zli wallet echo` to get.
No need to modify `init_implementation`, will `upgrade` later. Both transaction id and contract address will be print on
standard output.

#### 4. Deploy sshlist contract

Run `sh ./deploy_sshlist.sh`

Also need to make sure you put correct `init_admin` and `proxy_address`.

#### 5. Upgrade proxy contract

Run `zli contract call -a 0x09710e00256db2e3db4b44f597f17f3d97f06318 -t upgradeTo -r "[{\"vname\":\"newImplementation\",\"type\":\"ByStr20\",\"value\":\"0x1256e7c364d4f5b4b579541d1483f4be9ab5bc3d\"}]" -f true`

or `sh ./upgradeTo.sh`

This allow you to make your proxy point to actual sshlist contract.

#### 6. Others

After doing deploy and upgrade, you can run any others script to call transitions. Just keep in mind <b>do not mess up

contract address, always modify the scripts to your own contract address</b>

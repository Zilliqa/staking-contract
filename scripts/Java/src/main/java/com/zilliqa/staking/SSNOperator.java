package com.zilliqa.staking;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.firestack.laksaj.account.Wallet;
import com.firestack.laksaj.contract.CallParams;
import com.firestack.laksaj.contract.Contract;
import com.firestack.laksaj.contract.ContractFactory;
import com.firestack.laksaj.contract.Value;
import com.firestack.laksaj.crypto.KeyTools;
import com.firestack.laksaj.jsonrpc.HttpProvider;
import com.firestack.laksaj.transaction.Transaction;
import com.firestack.laksaj.transaction.TxStatus;
import com.firestack.laksaj.utils.Bech32;
import com.firestack.laksaj.utils.Validation;

import java.util.Arrays;
import java.util.List;

import static com.firestack.laksaj.account.Wallet.pack;

public class SSNOperator {
    // the Zilliqa api endpoint
    private String api;
    // chain id for small scale test net or community devnet(333) or mainnet(1)
    private int chainId;
    // basically this is the private key for ssn operator
    private String SSNPrivateKey;
    // the staking smart contract address resides on Zilliqa network, should be bech32 format
    private String stakingProxyAddress;

    private HttpProvider provider;
    private Wallet wallet;
    private Contract SSNContractInstance;


    public SSNOperator(String api, int chainId, String SSNPrivateKey, String stakingProxyAddress) throws Exception {
        this.api = api;
        this.chainId = chainId;
        this.SSNPrivateKey = SSNPrivateKey;
        if (!Validation.isBech32(stakingProxyAddress)) {
            this.stakingProxyAddress = Bech32.toBech32Address(stakingProxyAddress);
        } else {
            this.stakingProxyAddress = stakingProxyAddress;
        }

        this.provider = new HttpProvider(api);
        this.wallet = new Wallet();
        wallet.addByPrivateKey(this.SSNPrivateKey);
        wallet.setProvider(this.provider);

        List<Value> init = Arrays.asList();
        ContractFactory factory = ContractFactory.builder().provider(this.provider).signer(wallet).build();
        this.SSNContractInstance = factory.atContract(this.stakingProxyAddress, "", (Value[]) init.toArray(), "");
    }

    /**
     * @param amount   staking amount
     * @param attempts attempt times for polling transaction
     * @param interval interval time in seconds between each polling
     * @return transaction id
     * @throws Exception
     */
    public String stakeDeposit(String amount, int attempts, int interval) throws Exception {
        String normalAddr = KeyTools.getAddressFromPrivateKey(this.SSNPrivateKey);
        String publicKey = KeyTools.getPublicKeyFromPrivateKey(this.SSNPrivateKey, true);
        Integer nonce = Integer.valueOf(this.provider.getBalance(normalAddr).getResult().getNonce());
        CallParams params = CallParams.builder().nonce(String.valueOf(nonce + 1)).version(String.valueOf(pack(this.chainId, 1))).gasPrice("1000000000").gasLimit("30000").senderPubKey(publicKey).amount(amount).build();
        List<Value> values = Arrays.asList();
        Transaction tx = this.SSNContractInstance.call("stake_deposit", (Value[]) values.toArray(), params, attempts, interval);
        if (tx.getStatus() == TxStatus.Confirmed && tx.getReceipt().isSuccess()) {
            // check receipt and event logs
            List<Object> eventLogs = tx.getReceipt().getEvent_logs();
            if (eventLogs != null && eventLogs.size() != 0) {
                ObjectMapper objectMapper = new ObjectMapper();
                String events = objectMapper.writeValueAsString(eventLogs.get(0));
                EventLogEntry eventLogEntry = objectMapper.reader().forType(EventLogEntry.class).readValue(events);
                if (!eventLogEntry._eventname.equals("SSN updated stake") || !eventLogEntry._eventname.equals("SSN updated buffered stake")) {
                    // todo carefully handle this
                    System.out.println(events);
                } else {
                    // todo everything is fine
                }
            }

        } else {
            // todo please carefully handle this situation
        }
        return tx.getID();
    }
}

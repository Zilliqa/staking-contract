package com.zilliqa.staking;

import org.junit.Test;

public class SSNOperatorTest {
    @Test
    public void stakeDeposit() throws Exception {
        SSNOperator ssnOperator = new SSNOperator("https://zilliqa-isolated-server.zilliqa.com/", 1,
                "40a08154418fcc0026e9f93f6ed16c6c6a499cbcda1335b581084f18105d1c7b",
                "40a57198730c58a59eb67d7d299e55dd958090ff");

        String tx = ssnOperator.stakeDeposit("20000000", 100, 3);
        System.out.println(tx);
    }
}

package com.zilliqa.staking;

import org.junit.Test;

public class SSNOperatorTest {

    private static String api = "https://dev-api.zilliqa.com/";
    private static int chainId = 333;
    private static String ssnPrivateKey = "";
    private static String proxyAddress = "";
    SSNOperator ssnOperator = new SSNOperator(api, chainId, ssnPrivateKey, proxyAddress);

    public SSNOperatorTest() throws Exception {
    }


    @Test
    public void stakeDeposit() throws Exception {
        String tx = ssnOperator.stakeDeposit("1000", 100, 3);
        System.out.println(tx);
    }

    @Test
    public void withdrawStakeAmount() throws Exception {
        String tx = ssnOperator.withdrawStakeAmount("1000", 100, 3);
        System.out.println(tx);
    }

    @Test
    public void withdrawStakeReward() throws Exception {
        String tx = ssnOperator.withdrawStakeRewards(100, 3);
        System.out.println(tx);
    }

    @Test
    public void getStakeAmount() throws Exception {
        String state = ssnOperator.getStakeAmount();
        System.out.println(state);
    }

    @Test
    public void getStakeBufferedAmount() throws Exception {
        String state = ssnOperator.getStakeBufferedAmount();
        System.out.println(state);
    }

    @Test
    public void getStakeRewards() throws Exception {
        String state = ssnOperator.getStakeRewards();
        System.out.println(state);
    }

    @Test
    public void getActiveStatue() throws Exception {
        System.out.println(ssnOperator.getActiveStatus());
    }

    @Test
    public void getNodeStatus() throws Exception {
        System.out.println(ssnOperator.getNodeStatus("https://dev-api.zilliqa.com",10));
    }
}

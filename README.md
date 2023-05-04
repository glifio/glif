# glif cli

## Getting started - Creating your Agent

The Agent is a crucial component of the underlying [GLIF Pools Protocol](https://glif.io/docs) (the Protocol on which the Infinity Pool is built) - the Agent is a wrapper contract around one or more [Miner Actors](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/miner/miner_actor.go). The Agent is primarily responsible for:

1. Adding miner(s) to the Agent
2. Borrowing FIL from a pool
3. Making a payment to a pool
4. Pushing funds to any of the Agent's miners
5. Pulling funds up from any of the Agent's miners
6. Withdrawing funds from the Agent
7. Removing a miner
8. Changing operator addresses

**Owner** - The Owner Address owns your agent. The Agent's Owner is like your Miner Owner - it has the permission to call any method on your Agent. Additionally, it is the only address that is able to: (1) Borrow funds from pools into the Agent, (2) Withdraw funds from the Agent to a recipient, and (3) Remove a miner from the agent. We recommend keeping your Owner private key cold and not kept on a machine that's directly connected to the internet.<br />
**Operator** - The Operator Address is primarily useful for automation - the operator can (1) Make a payment to the pool, and (2) Push/pull funds from the Agent back and forth with the Agent's miner actors.<br />
**Requester** - The requester signs your requests to the Agent Data Oracle (ADO) to ensure that only _you_ can request a signed credential for your Agent. 

To create your Agent, first, we need to create 3 new addresses:

`glif agent keys new`<br />

This command will randomly generate 3 new private keys that represent your `owner`, `operator`, and `requester`. It will store the keys in `$HOME/.glif/config/keys.toml`

You should see the output:

```
➜ ✗ glif agent keys new
2023/05/03 19:53:37 Owner address: 0x8b35624Ed57789D18445142a51A4a51eFb375F26 (ETH), f410frm2wetwvo6e5dbcfcqvfdjffd35toxzgtslcqja (FIL)
2023/05/03 19:53:37 Operator address: 0x69D2CE31DDABF4A7098a2547147c13cF10F5Ea7b (ETH), f410fnhjm4mo5vp2kocmkevdri7atz4ipl2t3zcdiani (FIL)
2023/05/03 19:53:37 Request key: 0x115fFf27B67875032D6E6D2F28a3aAbC31A69f54 (ETH), f410fcfp76j5wpb2qgllonuxsri5kxqy2nh2uqbvcnoy (FIL)
```

Next, we need to fund our owner key. To do this, please navigate over to the [GLIF Wallet](https://glif.io/wallet), and send some funds to your owner address. **IMPORTANT** - do NOT manually craft and send a `method 0` send transaction to an EVM address, passing it `value`. Use [fil-forwarder](https://docs.filecoin.io/smart-contracts/filecoin-evm-runtime/filforwader/) instead.

Once you've funded your owner key, verify:

```
➜ ✗ glif agent keys balance
2023/05/03 20:55:13 owner balance: 0 FIL - (FEVM) f410fr...cqja (EVM) 0x8b35...5F26
2023/05/03 20:55:13 operator balance: 0 FIL - (FEVM) f410fn...iani (EVM) 0x69D2...Ea7b
2023/05/03 20:55:13 request balance: 0 FIL - (FEVM) f410fc...cnoy (EVM) 0x115f...9f54
```

Lastly, you can go ahead and create your agent:

```
➜ ✗ glif agent create
```
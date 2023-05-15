# glif cli

## Getting started - Installing the code
First, clone the repo:
`git clone git@github.com:glifio/cli.git`<br />
`cd cli`<br />
`sudo make install`<br />
`make calibnet-config`<br />

Now when you run:
```
➜ ✗ glif wallet new
```

You should see:

```
➜  cli git:(main) ✗ glif --help
Usage:
  glif [command]

Available Commands:
  agent         Commands for interacting with the Glif Agent
  completion    Generate the autocompletion script for the specified shell
  help          Help about any command
  ifil          Commands for interacting with the Infinity Pool Liquid Staking Token (iFIL)
  infinity-pool Commands for interacting with the Infinity Pool
  pools         Commands for interacting with the GLIF Pools Protocol
  wallet        Manage Glif wallets

Flags:
      --config string   config file (default is $HOME/.config/glif/config.toml)
  -h, --help            help for glif
  -t, --toggle          Help message for toggle

Use "glif [command] --help" for more information about a command.
```

## Creating your Agent

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

`glif wallet new`<br />

This command will randomly generate 3 new private keys that represent your `owner`, `operator`, and `requester`. It will store the keys in `$HOME/.glif/config/keys.toml`

You should see the output:

```
➜ ✗ glif wallet new
2023/05/03 19:53:37 Owner address: 0x8b35624Ed57789D18445142a51A4a51eFb375F26 (ETH), f410frm2wetwvo6e5dbcfcqvfdjffd35toxzgtslcqja (FIL)
2023/05/03 19:53:37 Operator address: 0x69D2CE31DDABF4A7098a2547147c13cF10F5Ea7b (ETH), f410fnhjm4mo5vp2kocmkevdri7atz4ipl2t3zcdiani (FIL)
2023/05/03 19:53:37 Request key: 0x115fFf27B67875032D6E6D2F28a3aAbC31A69f54 (ETH), f410fcfp76j5wpb2qgllonuxsri5kxqy2nh2uqbvcnoy (FIL)
```

Next, we need to fund our owner key. To do this, please navigate over to the [GLIF Wallet](https://glif.io/wallet), and send some funds to your owner address. **IMPORTANT** - do NOT manually craft and send a `method 0` send transaction to an EVM address, passing it `value`. Use [fil-forwarder](https://docs.filecoin.io/smart-contracts/filecoin-evm-runtime/filforwader/) instead.

Once you've funded your owner key, verify:

```
➜ ✗ glif wallet balance
2023/05/15 10:57:53 owner balance: 2.00 FIL
2023/05/15 10:57:53 operator balance: 0.00 FIL
2023/05/15 10:57:53 request balance: 0.00 FIL
```

Lastly, you can go ahead and create your agent:

```
➜ ✗ glif agent create
```

This should


```
2023/05/15 17:38:18 pools sdk: agent create: failed to estimate gas: CallWithGas failed: call raw get actor: resolution lookup failed (t410f52ogwtdgdaafchdrj54tiftjdjpbn3kiemnaiay): resolve address t410f52ogwtdgdaafchdrj54tiftjdjpbn3kiemnaiay: actor not found
```

## Add a miner to your Agent
```

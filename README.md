# glif cli

## Getting started (Calibnet)
First, clone the repo:
`git clone git@github.com:glifio/cli.git`<br />
`cd cli`<br />
`make calibnet`<br />
`sudo make install`<br />
`make calibnet-config`<br />

## Getting started (Mainnet)
First, clone the repo:
`git clone git@github.com:glifio/cli.git`<br />
`cd cli`<br />
`make glif`<br />
`sudo make install`<br />
`make config`<br />

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

If all goes successfully, you can run:

```
➜ ✗ cli git:(main) glif agent info
```

And you should see something like:

```
➜  cli git:(main) glif agent info

BASIC INFO
Agent Address: 0xbf11d189D528736d25D3a342C826cF60253Df41c
Agent ID: 4
Agent Version: 1

AGENT ASSETS
0.000000 FIL
|
INFINITY POOL ACCOUNT
No account exists with the Infinity Pool
```

If you see an error that looks like:

```
2023/05/15 17:38:18 pools sdk: agent create: failed to estimate gas: CallWithGas failed: call raw get actor: resolution lookup failed (t410f52ogwtdgdaafchdrj54tiftjdjpbn3kiemnaiay): resolve address t410f52ogwtdgdaafchdrj54tiftjdjpbn3kiemnaiay: actor not found
```

It means that your owner key is not properly funded. You must send FIL to this actor before creating an Agent.

## Add a miner to your Agent

Adding a miner to your Agent occurs in two steps:

1. Setting the owner address on your miner actor to point to your agent
2. Adding the miner to your Agent

### Step 1 - Proposing an Ownership change

**NOTE** - if your owner key is the default wallet selected on your Lotus daemon, you can adjust your `~/.glif/config.toml` to point to your own Lotus daemon to use our built-in change-miner-owner command. For example, with a `~/.glif/config.toml` that looks like:

```
[daemon]
rpc-url = 'http://localhost:1234/rpc/v1'
token = 'eyJh...om49Vu1w'
```

Then you can simply run:<br />
`glif agent miners change-owner <miner addr>`<br />

This will propose the ownership change. Alternatively, you can call:<br />
`glif agent info`<br />

To retrieve your delegated `f4` address of your Agent. Then you can go ahead and manually propose an ownership change, setting your Agent's `f4` address as the new owner.

### Step 2 - Adding miner to your Agent

Once you've successfully proposed an ownership change to your agent, you can then call `glif agent miners add <miner addr>` to pledge your miner.

You should see:

```
2023/05/15 14:05:44 Adding miner f0xxx to agent 0x...
|Transaction: 0x....
Successfully added miner f0xxx to agent
```

### Confirm - Miner added successfully:

You can call `glif agent miners list` and you should see your new miner in the returned list!

# Command Reference

## Preview Flag

Several of the critical commands have a `--preview` flag that allows you to perform a dry-run and get an indication of the impact of that operation on the financial position of your agent. Check whether a command has a `--preview` flag by calling the `--help` flag on that command.

## Borrow funds

You can borrow funds from the Infinity Pool by calling `glif agent borrow <amount>`. For example, to borrow 1 FIL, you can call:

```
➜ ✗ glif agent borrow 1
2023/05/15 14:08:44 Borrowing 1 FIL from Infinity Pool
|Transaction: 0x....
Successfully borrowed 1 FIL from Infinity Pool
```

## Make a payment

### To-Current payment

You can pay all current fees by calling `glif agent pay to-current`. For example, to pay the to-current amount, you can call:

```
➜ ✗ glif agent pay to-current
2023/05/15 14:08:44 Making a to-current payment to Infinity Pool
|Transaction: 0x....
Successfully paid 5 FIL
```

### Principal payment

To make a payment against the principal, make the following call `glif agent pay principal`. Note, that the amount specified is how much principle you wish to pay off, but the call will also pay off an outstanding fees, so the total amount paid will be equal to having called `glif agent pay to-current` plus the amount specified against the principle. For example:

```
➜ ✗ glif agent pay principal 10
2023/05/15 14:08:44 Paying fees of 5 FIL to the Infinity Pool
2023/05/15 14:08:45 Paying principle of 10 FIL to the Infinity Pool
|Transaction: 0x....
Successfully paid principal amount to Infinity Pool
```

### Custom payment

You can make a payment to the Infinity Pool by calling `glif agent pay custom <amount>`. For example, to pay 1 FIL, you can call:

```
➜ ✗ glif agent pay custom 1
2023/05/15 14:08:44 Paying 1 FIL to the Infinity Pool
|Transaction: 0x....
Successfully paid 1 FIL
```

## Push funds to a miner

You can push funds to a miner by calling `glif agent push <miner addr> <amount>`. For example, to push 1 FIL to a miner, you can call:

```
➜ ✗ glif agent push f0xxx 1
|Transaction: 0x....
Successfully pushed funds down to miner f0xxx
```

## Pull funds from a miner

You can pull funds from a miner by calling `glif agent pull <miner addr> <amount>`. For example, to pull 1 FIL from a miner, you can call:

```
➜ ✗ glif agent pull f0xxx 1
2023/05/15 14:08:44 Pulling 1 FIL from miner f0xxx
|Transaction: 0x....
Successfully pulled funds up from miner f0xxx
```

## Withdraw funds from your Agent

You can withdraw funds from your Agent by calling `glif agent withdraw <amount>`. For example, to withdraw 1 FIL from your Agent, you can call:

```
➜ ✗ glif agent withdraw 1
2023/05/15 14:08:44 Withdrawing 1 FIL from Agent
|Transaction: 0x....
Successfully withdrew 1 FIL
```

## Remove a miner from your Agent

You can remove a miner from your Agent by calling `glif agent miners remove <miner addr>`. For example, to remove a miner from your Agent, you can call:

```
➜ ✗ glif agent miners remove f0xxx f3xxx
2023/05/15 14:08:44 Removing miner f0xxx from agent
|Transaction: 0x....
Successfully removed miner f0xxx from Agent
```

# Glif Autopilot

Glif Autopilot allows for the automatic payment of fees and/or principal to the Infinity Pool. It is a daemon that runs in the background and will automatically make payments to the Infinity Pool on your behalf. It is configured via a `config.toml` file that is located in your `~/.glif` directory.

## Configuration

The `config.toml` file is located in your `~/.glif` directory. It is a TOML file that looks like:

```
[autopilot]
# payment-type can be one of "to-current", "principal", or "custom", each working as described by the matching cli commands
payment-type = "principal"
# amount is in FIL and only applies to payment-type = "custom" or "principal"
amount = 5
# frequency is the frequency at which payments are made in days
frequency = 5
```


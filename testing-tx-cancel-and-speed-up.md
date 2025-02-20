# Testing the tx cancel and speed-up commands

The `glif tx cancel` and `glif tx speed-up` commands can be used to replace pending
transactions with either a dummy transaction with no effect (for cancel) or a new
transaction with a 25% increate in gas premium to try to boost it for inclusion
when there is on-chain congestion.

Typically, there is not much on-chain congestion on mainnet, so it can be difficult
to test these commands as submitted transactions tend to go through in 30 seconds
and don't get "stuck" waiting for a miner to select them.

## Fee Cap workaround

A workaround is to submit a transaction with a "Fee Cap" that is below the current
gas price ... that way the transaction will get stuck in the mempool until a
replacement transaction is made with a higher fee cap. This is useful for testing.

However, a standard Lotus node performs a sanity check on transactions submitted
via it's RPC API and will refuse to post the transaction to the mempool if the
fee cap is set too low.

A hack is to modify the Lotus node to remove this check:

```
diff --git a/chain/messagepool/messagepool.go b/chain/messagepool/messagepool.go
index 80dd15849..89c84724e 100644
--- a/chain/messagepool/messagepool.go
+++ b/chain/messagepool/messagepool.go
@@ -744,9 +744,11 @@ func (mp *MessagePool) checkMessage(ctx context.Context, m *types.SignedMessage)
                return ErrMessageValueTooHigh
        }
 
-       if m.Message.GasFeeCap.LessThan(minimumBaseFee) {
-               return ErrGasFeeCapTooLow
-       }
+       /*
+               if m.Message.GasFeeCap.LessThan(minimumBaseFee) {
+                       return ErrGasFeeCapTooLow
+               }
+       */
 
        if err := mp.VerifyMsgSig(m); err != nil {
                return xerrors.Errorf("signature verification failed: %s", err)
```

For convenience, there is a patched Lotus node that may be available ... modify ~/.glif/config.toml to
use these values:

```
rpc-url = 'http://lotus.v6z.me:1234/rpc/v1'
token = ''
```

## Testing: Replacing a transaction using a nonce

Make a transaction that will intentionally get stuck:

```
glif infinity-pool deposit-fil 0.00001 --from operator --gas-premium=0 --gas-fee-cap=1
```

This transaction will get "stuck" in the mempook, so you will need to interrupt the command with Ctrl-C.

Note: We are avoiding agent transactions that require a credential, as it is necessary to wait 5 minutes before a replacement credential can be issued.

Go to: `https://www.glif.io/en/tx/<tx hash>` and you will see the transaction is "pending". 

If you have access to a Lotus node, you can watch the transaction directly in the mempool like this:

```
while true; do lotus mpool pending --to f410f45skz4bnrn6cduvwvdyks3dykqpa3q75ofrfdty; sleep 5; echo; echo; done
```

You can see the pending transactions using the cli like this:

```
$ ./glif tx list-pending operator
Nonce  Transaction                                                         Gas Premium  Gas Fee Cap  
208    0x8d401c728ebeacac8fc2dc5f9793ba4feaedc47473562ad7c330ad26b9dd36b7  0            1 
```

In this case, the fee cap is too low, so it is stuck.

To replace it using a nonce, do this:

```
glif infinity-pool deposit-fil 0.00001 --from operator --gas-limit 119265312 --gas-premium-multiply 2 --nonce <nonce>
```

(Replace <nonce> with the nonce shown in the list-pending table)

For this particular example, it is necessary to set the gas-limit as sometimes Lotus will make a bad estimate causing the transaction to fail.

If you use the command above, the transaction should be executed and will no longer be in the mempool.

Go to: `https://www.glif.io/en/tx/<replacement tx hash>` and you will see the transaction has executed with the updated gas premium. 

## Testing: Canceling a transaction

Again, we will create a stuck transaction:

```
glif infinity-pool deposit-fil 0.00001 --from operator --gas-premium=0 --gas-fee-cap=1
```

Use Ctrl-C to exit as the transaction will never complete.

Look at the pending transactions in the mempool:

```
$ ./glif tx list-pending operator
Nonce  Transaction                                                         Gas Premium  Gas Fee Cap  
209    0x160fee15fc5a05993a7ffd9cb389b8a2a1ef4927117e04bf66e125fb324890b4  0            1  
```

Cancel it:

```
$ ./glif tx cancel 0x160fee15fc5a05993a7ffd9cb389b8a2a1ef4927117e04bf66e125fb324890b4
Replacement transaction sent: 0xc06e198c0c8262efcfab164d000c56fe06858b56982a8917c1a743960f56e487
```

Go to: `https://www.glif.io/en/tx/<replacement tx hash>` and you will see the transaction has executed,
but the ethereum call has been replaced with a dummy transaction (a zero value transfer to itself).

## Testing: Speeding up a transaction

Again, we will create a stuck transaction:

```
glif infinity-pool deposit-fil 0.00001 --from operator --gas-premium=0 --gas-fee-cap=1
```

Use Ctrl-C to exit as the transaction will never complete.

Look at the pending transactions in the mempool:

```
$ ./glif tx list-pending operator
Nonce  Transaction                                                         Gas Premium  Gas Fee Cap  
210    0x2d1ccbfcfce29ac45e3788a8bfbf76bb83299e44a983ef50509754ec01634c9d  0            1 
```

Speed it up ... this replaces the ethereum transaction with another ethereum transaction with the same
parameters, but with 25% higher premium, so hopefully it will execute faster when there is network
congestion:

```
$ ./glif tx speed-up 0x2d1ccbfcfce29ac45e3788a8bfbf76bb83299e44a983ef50509754ec01634c9d
Replacement transaction sent: 0xe0d2f0b3f8e2bf7109a710ebeaf5293cec4ec3b633016f61fa9a13ef67bf9d69
```

(Note: there is a --gas-premium flag that can be set to add even more premium than the automatically
computed 25% increase)

Go to: `https://www.glif.io/en/tx/<replacement tx hash>` and you will see the transaction has executed,
but the ethereum call has been replaced with an equivalent call, but with a higher gas premium and gas fee cap.






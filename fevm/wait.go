package fevm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WaitTx(ctx context.Context, c *ethclient.Client, hash common.Hash, ch chan *types.Receipt) {
	for {
		time.Sleep(time.Millisecond * 5000)

		tx, err := c.TransactionReceipt(ctx, hash)
		if err == nil && tx != nil {
			ch <- tx
			return
		}
	}
}

func WaitForNextBlock(ctx context.Context, c *ethclient.Client, current *big.Int, ch chan bool) {
	target := current.Uint64() + 1
	for {
		time.Sleep(time.Millisecond * 5000)

		b, err := c.BlockNumber(ctx)
		if err == nil && b >= target {
			ch <- true
			return
		}
	}
}

func waitForReceipt(hash common.Hash, client *ethclient.Client) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 480*time.Second)
	defer cancel()

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	ch := make(chan *types.Receipt)
	go WaitTx(ctx, client, hash, ch)

	var receipt *types.Receipt
	select {
	case <-ctx.Done():
		s.Stop()
		return nil, fmt.Errorf("timed out waiting for transaction")
	case receipt = <-ch:
		msg := fmt.Sprintf(" Receipt received at block %v. Waiting for next block.\n", receipt.BlockNumber)
		s.Suffix = msg

		ch := make(chan bool)
		go WaitForNextBlock(ctx, client, receipt.BlockNumber, ch)
		select {
		case <-ctx.Done():
			s.Stop()
			return nil, fmt.Errorf("timed out waiting for transaction")
		case <-ch:
			s.Stop()
			fmt.Println("Transaction receipt received.")
			fmt.Printf("Status: %v\n", receipt.Status)
			return receipt, nil
		}
	}
}

func WaitForReceipt(hash common.Hash) error {
	eapi, err := Connection().ConnectEthClient()
	if err != nil {
		return err
	}
	defer eapi.Close()

	if _, err := waitForReceipt(hash, eapi); err != nil {
		return err
	}
	return nil
}

func WaitReturnReceipt(hash common.Hash) (*types.Receipt, error) {
	eapi, err := Connection().ConnectEthClient()
	if err != nil {
		return nil, err
	}
	defer eapi.Close()

	var receipt *types.Receipt
	if receipt, err = waitForReceipt(hash, eapi); err != nil {
		return nil, err
	}
	return receipt, nil
}

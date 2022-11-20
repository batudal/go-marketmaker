package wrapper

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/decoded-labs/go-marketmaker/abis/weth"
	"github.com/decoded-labs/go-marketmaker/helpers/logger"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// swaps all tokens into wbnb
func Wrap(
	_bot_id int,
	_weth_address common.Address,
	_wallet_address common.Address,
	_wallet_key *ecdsa.PrivateKey,
	_rpc string,
	_chainid int,
) {
	// connect to devm
	client, err := ethclient.Dial(_rpc)
	if err != nil {
		log.Fatal(err)
	}
	// load contract
	weth, err := weth.NewWeth(_weth_address, client)
	if err != nil {
		log.Fatal(err)
	}
	// nonce & gas price
	nonce, err := client.PendingNonceAt(context.Background(), _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// set transactOps
	chain_id := big.NewInt(int64(_chainid))
	auth, err := bind.NewKeyedTransactorWithChainID(_wallet_key, chain_id)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(1000000000000000000) // payable method
	auth.GasLimit = uint64(0)
	auth.GasPrice = gas_price

	// do tx, wait, and log
	wrap_tx, err := weth.Deposit(auth)
	if err != nil {
		log.Fatal(err)
	}
	res, err := bind.WaitMined(context.Background(), client, wrap_tx)
	if err != nil {
		log.Fatal(err)
	}
	logger.Wrap(
		"wrap",
		_bot_id,
		_wallet_address,
		wrap_tx.Hash(),
		res.Status,
		auth.Value,
	)
}

// swaps all tokens into wbnb
func UnWrap(
	_bot_id int,
	_weth_address common.Address,
	_wallet_key *ecdsa.PrivateKey,
	_wallet_address common.Address,
	_rpc string,
	_chain_id int,
) {
	// connect to devm
	client, err := ethclient.Dial(_rpc)
	if err != nil {
		log.Fatal(err)
	}
	// load contract
	weth, err := weth.NewWeth(_weth_address, client)
	if err != nil {
		log.Fatal(err)
	}
	// nonce & gas price
	nonce, err := client.PendingNonceAt(context.Background(), _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// set transactOps
	chain_id := big.NewInt(int64(_chain_id))
	auth, err := bind.NewKeyedTransactorWithChainID(_wallet_key, chain_id)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(0)
	auth.GasPrice = gas_price

	// amount
	weth_balance, err := weth.BalanceOf(nil, _wallet_address)
	if err != nil {
		log.Fatal(err)
	}

	// do tx, wait, and log
	wrap_tx, err := weth.Withdraw(auth, weth_balance)
	if err != nil {
		log.Fatal(err)
	}
	res, err := bind.WaitMined(context.Background(), client, wrap_tx)
	if err != nil {
		log.Fatal(err)
	}
	logger.Wrap(
		"unwrap",
		_bot_id,
		_wallet_address,
		wrap_tx.Hash(),
		res.Status,
		weth_balance,
	)
}

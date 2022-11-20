package approve

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/decoded-labs/go-marketmaker/abis/erc20"
	"github.com/decoded-labs/go-marketmaker/helpers/logger"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Approve(
	_wallet_address common.Address,
	_wallet_key *ecdsa.PrivateKey,
	_token0_address common.Address,
	_token1_address common.Address,
	_router common.Address,
	_rpc string,
	_chainid int,
	_bot_id int,
) {
	client, err := ethclient.Dial(_rpc)
	if err != nil {
		log.Fatal(err)
	}
	token0, err := erc20.NewErc20(_token0_address, client)
	if err != nil {
		log.Fatal(err)
	}
	token1, err := erc20.NewErc20(_token1_address, client)
	if err != nil {
		log.Fatal(err)
	}
	chain_id := big.NewInt(int64(_chainid))
	auth, err := bind.NewKeyedTransactorWithChainID(_wallet_key, chain_id)
	if err != nil {
		log.Fatal(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(0)
	auth.GasPrice = gas_price

	approve_tx0, err := token0.Approve(auth, _router, abi.MaxUint256)
	if err != nil {
		log.Fatal(err)
	}

	tx0_res, err := bind.WaitMined(context.Background(), client, approve_tx0)
	if err != nil {
		log.Fatal(err)
	}

	logger.Approve(
		_bot_id,
		_wallet_address,
		approve_tx0.Hash(),
		tx0_res.Status,
	)

	nonce, err = client.PendingNonceAt(context.Background(), _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	approve_tx1, err := token1.Approve(auth, _router, abi.MaxUint256)
	if err != nil {
		log.Fatal(err)
	}
	tx1_res, err := bind.WaitMined(context.Background(), client, approve_tx1)
	if err != nil {
		log.Fatal(err)
	}
	logger.Approve(
		_bot_id,
		_wallet_address,
		approve_tx1.Hash(),
		tx1_res.Status,
	)
}

package funding

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/decoded-labs/go-marketmaker/helpers/logger"
	"github.com/decoded-labs/go-marketmaker/helpers/wrapper"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Fund(
	_bot_id int,
	_admin_address common.Address,
	_admin_key *ecdsa.PrivateKey,
	_wallet_address common.Address,
	_wallet_key *ecdsa.PrivateKey,
	_weth_address common.Address,
	_rpc string,
	_chainid int,
) {
	client, err := ethclient.Dial(_rpc)
	if err != nil {
		log.Fatal(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), _admin_address)
	if err != nil {
		log.Fatal(err)
	}
	gas_limit := uint64(21000)
	gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	eth_needed := big.NewInt(1100000000000000000) // 1.1 eth
	eth_tx := types.NewTransaction(nonce, _wallet_address, eth_needed, gas_limit, gas_price, nil)
	signed_eth_tx, err := types.SignTx(eth_tx, types.NewEIP155Signer(big.NewInt(int64(_chainid))), _admin_key)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signed_eth_tx)
	if err != nil {
		log.Fatal(err)
	}
	eth_tx_res, err := bind.WaitMined(context.Background(), client, signed_eth_tx)
	if err != nil {
		log.Fatal(err)
	}
	logger.Fund(
		_bot_id,
		eth_needed,
		_wallet_address,
		eth_tx.Hash(),
		eth_tx_res.Status,
	)

	time.Sleep(time.Duration(20) * time.Second) // wait 20 secs

	wrapper.Wrap(
		_bot_id,
		_weth_address,
		_wallet_address,
		_wallet_key,
		_rpc,
		_chainid,
	)
}

func UnFund(
	_bot_id int,
	_admin_address common.Address,
	_wallet_address common.Address,
	_wallet_key *ecdsa.PrivateKey,
	_rpc string,
	_weth_address common.Address,
	_chainid int,
) {
	client, err := ethclient.Dial(_rpc)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Duration(20) * time.Second) // wait 20 secs

	wrapper.UnWrap(_bot_id, _weth_address, _wallet_key, _wallet_address, _rpc, _chainid)

	time.Sleep(time.Duration(20) * time.Second) // wait 20 secs

	// unfund remaining eth
	nonce, err := client.PendingNonceAt(context.Background(), _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	gas_limit := uint64(21000)
	gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	eth_remaining, err := client.PendingBalanceAt(context.Background(), _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	eth_remaining_after_tx := new(big.Int).Sub(eth_remaining, new(big.Int).Mul(gas_price, big.NewInt(int64(gas_limit))))
	eth_tx := types.NewTransaction(nonce, _admin_address, eth_remaining_after_tx, gas_limit, gas_price, nil)
	signed_eth_tx, err := types.SignTx(eth_tx, types.NewEIP155Signer(big.NewInt(int64(_chainid))), _wallet_key)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signed_eth_tx)
	if err != nil {
		log.Fatal(err)
	}

	eth_tx_res, err := bind.WaitMined(context.Background(), client, signed_eth_tx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Unfund(
		_bot_id,
		eth_remaining_after_tx,
		_wallet_address,
		eth_tx.Hash(),
		eth_tx_res.Status,
	)
}

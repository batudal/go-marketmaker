package strategy

import (
	"crypto/ecdsa"
	"time"

	"github.com/decoded-labs/go-marketmaker/helpers/approve"
	"github.com/decoded-labs/go-marketmaker/helpers/funding"
	"github.com/decoded-labs/go-marketmaker/helpers/generate_wallet"
	"github.com/decoded-labs/go-marketmaker/helpers/logger"
	"github.com/decoded-labs/go-marketmaker/helpers/swap"

	"github.com/ethereum/go-ethereum/common"
)

type Bot_config struct {
	router, factory, weth, token      common.Address
	rpc                               string
	chainid, swap_interval, max_swaps int
}

func New_Dumb(
	router common.Address,
	factory common.Address,
	weth common.Address,
	token common.Address,
	rpc string,
	chainid int,
	swap_interval int,
	max_swaps int,
) Bot_config {
	return Bot_config{
		router:        router,
		factory:       factory,
		weth:          weth,
		token:         token,
		rpc:           rpc,
		chainid:       chainid,
		swap_interval: swap_interval,
		max_swaps:     max_swaps,
	}
}

func (_instance_config Bot_config) Dumb(
	_admin_address common.Address,
	_admin_key *ecdsa.PrivateKey,
	_live_instances *int,
	_progress chan bool,
	_bot_id int,
) {
	wallet_key, wallet_address := generate_wallet.Generate()
	logger.Wallet(wallet_address, wallet_key, _bot_id)

	funding.Fund(
		_bot_id,
		_admin_address,
		_admin_key,
		wallet_address,
		wallet_key,
		_instance_config.weth,
		_instance_config.rpc,
		_instance_config.chainid,
	)

	_progress <- true
	approve.Approve(
		wallet_address,
		wallet_key,
		_instance_config.weth,
		_instance_config.token,
		_instance_config.router,
		_instance_config.rpc,
		_instance_config.chainid,
		_bot_id,
	)

	for swap_counter := 1; swap_counter < _instance_config.max_swaps; swap_counter++ {
		swap.Swap(
			_bot_id,
			swap_counter,
			_instance_config.router,
			_instance_config.factory,
			_instance_config.weth,
			_instance_config.token,
			wallet_key,
			wallet_address,
			_instance_config.rpc,
			_instance_config.chainid,
		)
		time.Sleep(time.Duration(_instance_config.swap_interval) * time.Second) // add randomization
		if swap_counter+1 == _instance_config.max_swaps {
			swap.SwapFinal(
				_bot_id,
				swap_counter+1,
				_instance_config.router,
				_instance_config.factory,
				_instance_config.weth,
				_instance_config.token,
				wallet_key,
				wallet_address,
				_instance_config.rpc,
				_instance_config.chainid,
			)
			funding.UnFund(
				_bot_id,
				_admin_address,
				wallet_address,
				wallet_key,
				_instance_config.rpc,
				_instance_config.weth,
				_instance_config.chainid,
			)
			*_live_instances--
		}
	}
}

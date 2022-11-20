package swap

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/decoded-labs/go-marketmaker/abis/erc20"
	"github.com/decoded-labs/go-marketmaker/abis/factory"
	"github.com/decoded-labs/go-marketmaker/abis/pair"
	"github.com/decoded-labs/go-marketmaker/abis/router"
	"github.com/decoded-labs/go-marketmaker/helpers/logger"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Swap(
	_bot_id int,
	_swap_id int,
	_router common.Address,
	_factory common.Address,
	_weth_address common.Address,
	_token_address common.Address,
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

	// create factory & router instances
	router, err := router.NewRouter(_router, client)
	if err != nil {
		log.Fatal(err)
	}
	factory, err := factory.NewFactory(_factory, client)
	if err != nil {
		log.Fatal(err)
	}

	// derive pubkey
	public_key := _wallet_key.Public()
	public_key_ecdsa, ok := public_key.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// get nonce
	from_address := crypto.PubkeyToAddress(*public_key_ecdsa)
	nonce, err := client.PendingNonceAt(context.Background(), from_address)
	if err != nil {
		log.Fatal(err)
	}

	// set msg.value and get gas_price
	value := big.NewInt(0)
	gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// get balances
	weth, err := erc20.NewErc20(_weth_address, client)
	if err != nil {
		log.Fatal(err)
	}
	weth_balance, err := weth.BalanceOf(nil, _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	token, err := erc20.NewErc20(_token_address, client)
	if err != nil {
		log.Fatal(err)
	}
	token_balance, err := token.BalanceOf(nil, _wallet_address)
	if err != nil {
		log.Fatal(err)
	}

	// set pair addresses (token/weth)
	pair_address, err := factory.GetPair(nil, _token_address, _weth_address)
	if err != nil {
		log.Fatal(err)
	}

	// get reserves
	pair, err := pair.NewPair(pair_address, client)
	if err != nil {
		log.Fatal(err)
	}
	reserves, err := pair.GetReserves(nil)
	if err != nil {
		log.Fatal(err)
	}

	// float division of reserves for price
	reserve_weth := new(big.Float)
	reserve_weth.SetInt(reserves.Reserve1)
	reserve_token := new(big.Float)
	reserve_token.SetInt(reserves.Reserve0)
	token_price_f := new(big.Float)
	token_price_f.Quo(reserve_weth, reserve_token)

	// derive values of held in wallet
	weth_balance_f := new(big.Float)
	weth_balance_f.SetInt(weth_balance)

	token_balance_f := new(big.Float)
	token_balance_f.SetInt(token_balance)

	token_value := new(big.Float)
	token_value.Mul(token_price_f, token_balance_f)

	direction := token_value.Cmp(weth_balance_f)
	addresses := make([]common.Address, 2)

	// swap amount
	prandom := new(big.Int)
	prandom.SetInt64(time.Now().Unix() % 19)
	nominator := new(big.Int)
	nominator.Add(big.NewInt(80), prandom)
	denominator := big.NewInt(100)

	// set path according to direction
	amount_in := new(big.Int)
	var amount_quoted *big.Int
	if direction == 1 {
		addresses[0] = _token_address
		addresses[1] = _weth_address
		amount_in.Mul(token_balance, nominator)
		amount_in.Div(amount_in, denominator)
		amount_quoted, err = router.GetAmountOut(
			nil,
			amount_in,
			reserves.Reserve0, // reserve in
			reserves.Reserve1, // reserve out
		)
		if err != nil {
			log.Fatal(err)
		}
	} else { // equal also goes here
		addresses[0] = _weth_address
		addresses[1] = _token_address
		amount_in.Mul(weth_balance, nominator)
		amount_in.Div(amount_in, denominator)
		amount_quoted, err = router.GetAmountOut(
			nil,
			amount_in,
			reserves.Reserve1, // reserve in
			reserves.Reserve0, // reserve out
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 10% max slippage
	quotient := new(big.Int).Quo(amount_quoted, big.NewInt(10))
	amount_out_min := new(big.Int).Mul(quotient, big.NewInt(9))

	// set deadline
	deadline := new(big.Int)
	deadline.SetInt64(time.Now().Unix() + 600)

	// set transactOps
	chain_id := big.NewInt(int64(_chain_id))
	auth, err := bind.NewKeyedTransactorWithChainID(_wallet_key, chain_id)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(0)
	auth.GasPrice = gas_price

	// swap tx
	swap_tx, err := router.SwapExactTokensForTokens(
		auth,
		amount_in,
		amount_out_min,
		addresses,
		_wallet_address,
		deadline)
	if err != nil {
		log.Fatal(err)
	}
	res, err := bind.WaitMined(context.Background(), client, swap_tx)
	if err != nil {
		log.Fatal(err)
	}

	weth_balance, err = weth.BalanceOf(nil, _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	token_balance, err = token.BalanceOf(nil, _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	logger.Swap(
		_bot_id,
		_wallet_address,
		_swap_id,
		swap_tx.Hash(),
		res.Status,
		weth_balance,
		token_balance,
	)
}

// swaps all tokens into wbnb
func SwapFinal(
	_bot_id int,
	_swap_id int,
	_router common.Address,
	_factory common.Address,
	_weth_address common.Address,
	_token_address common.Address,
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

	// create factory & router instances
	router, err := router.NewRouter(_router, client)
	if err != nil {
		log.Fatal(err)
	}
	factory, err := factory.NewFactory(_factory, client)
	if err != nil {
		log.Fatal(err)
	}

	// get nonce
	nonce, err := client.PendingNonceAt(context.Background(), _wallet_address)
	if err != nil {
		log.Fatal(err)
	}

	// set msg.value and get gas_price
	value := big.NewInt(0)
	gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// get balance
	token, err := erc20.NewErc20(_token_address, client)
	if err != nil {
		log.Fatal(err)
	}
	token_balance, err := token.BalanceOf(nil, _wallet_address)
	if err != nil {
		log.Fatal(err)
	}

	// set pair addresses
	pair_address, err := factory.GetPair(nil, _token_address, _weth_address)
	if err != nil {
		log.Fatal(err)
	}

	// get reserves
	pair, err := pair.NewPair(pair_address, client)
	if err != nil {
		log.Fatal(err)
	}
	reserves, err := pair.GetReserves(nil)
	if err != nil {
		log.Fatal(err)
	}

	addresses := make([]common.Address, 2)
	addresses[0] = _token_address
	addresses[1] = _weth_address

	// set path according to direction
	amount_quoted, err := router.GetAmountOut(
		nil,
		token_balance,
		reserves.Reserve0, // reserve in
		reserves.Reserve1, // reserve out
	)
	if err != nil {
		log.Fatal(err)
	}

	// 10% max slippage - change in production(?)
	quotient := new(big.Int).Quo(amount_quoted, big.NewInt(10))
	amount_out_min := new(big.Int).Mul(quotient, big.NewInt(9))

	// set deadline
	deadline := new(big.Int)
	deadline.SetInt64(time.Now().Unix() + 600)

	// set transactOps
	chain_id := big.NewInt(int64(_chain_id))
	auth, err := bind.NewKeyedTransactorWithChainID(_wallet_key, chain_id)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(0)
	auth.GasPrice = gas_price

	// swap tx
	swap_tx, err := router.SwapExactTokensForTokens(
		auth,
		token_balance,
		amount_out_min,
		addresses,
		_wallet_address,
		deadline)
	if err != nil {
		log.Fatal(err)
	}
	res, err := bind.WaitMined(context.Background(), client, swap_tx)
	if err != nil {
		log.Fatal(err)
	}
	weth, err := erc20.NewErc20(_weth_address, client)
	if err != nil {
		log.Fatal(err)
	}
	weth_balance, err := weth.BalanceOf(nil, _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	token_balance, err = token.BalanceOf(nil, _wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	logger.Swap(
		_bot_id,
		_wallet_address,
		_swap_id,
		swap_tx.Hash(),
		res.Status,
		weth_balance,
		token_balance,
	)
}

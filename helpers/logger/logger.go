package logger

import (
	// "crypto"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Approve(
	_instance_id int,
	_instance_address common.Address,
	_tx_hash common.Hash,
	_tx_status uint64,
) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/approve.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Recorded transaction.",
		"transaction_type", "approve",
		"instance_id", _instance_id,
		"instance_address", _instance_address,
		"tx_hash", _tx_hash,
		"tx_status", _tx_status,
	)
}

func Swap(
	_instance_id int,
	_instance_address common.Address,
	_tx_id int,
	_tx_hash common.Hash,
	_tx_status uint64,
	_weth_balance *big.Int,
	_token_balance *big.Int,
) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/swap.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Recorded transaction.",
		"transaction_type", "swap",
		"instance_id", _instance_id,
		"instance_address", _instance_address,
		"swap_count", _tx_id,
		"tx_hash", _tx_hash,
		"tx_status", _tx_status,
		"weth_balance", _weth_balance,
		"token_balance", _token_balance,
	)
}

func Wrap(
	_event_type string,
	_instance_id int,
	_instance_address common.Address,
	_tx_hash common.Hash,
	_tx_status uint64,
	_amount *big.Int,
) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/wrap.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Recorded transaction.",
		"transaction_type", _event_type,
		"instance_id", _instance_id,
		"instance_address", _instance_address,
		"tx_hash", _tx_hash,
		"tx_status", _tx_status,
		"amount", _amount,
	)
}

func Fund(
	_instance_id int,
	_eth_amount *big.Int,
	_instance_address common.Address,
	_tx_hash common.Hash,
	_tx_status uint64,
) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/funding.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Recorded transaction.",
		"transaction_type", "fund_eth",
		"instance_id", _instance_id,
		"eth_amount", _eth_amount,
		"instance_address", _instance_address,
		"tx_hash", _tx_hash,
		"tx_status", _tx_status,
	)
}

func Unfund(
	_instance_id int,
	_eth_amount *big.Int,
	_instance_address common.Address,
	_tx_hash common.Hash,
	_tx_status uint64,
) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/funding.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Recorded transaction.",
		"transaction_type", "unfund_eth",
		"instance_id", _instance_id,
		"eth_amount", _eth_amount,
		"instance_address", _instance_address,
		"tx_hash", _tx_hash,
		"tx_status", _tx_status,
	)
}

func Error(
	_instance_id int,
	_instance_address common.Address,
	_tx_id int,
	_tx_hash common.Hash,
) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/tx.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		"instance_id", _instance_id,
		"instance_address", _instance_address,
		"tx_id", _tx_id,
		"tx_hash", _tx_hash,
	)
}

func Wallet(
	_wallet_address common.Address,
	_wallet_key *ecdsa.PrivateKey,
	_bot_id int,
) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/wallet.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	private_key_bytes := crypto.FromECDSA(_wallet_key)
	private_key := hexutil.Encode(private_key_bytes)[2:]
	logger := zap.New(core)
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Wallet generated.",
		"wallet_address", _wallet_address,
		"wallet_key", private_key,
		"bot_id", _bot_id,
	)
}

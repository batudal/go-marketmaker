package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/decoded-labs/go-marketmaker/helpers/config"
	"github.com/decoded-labs/go-marketmaker/strategy"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
)

func get_env_var(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv(key)
}

func main() {
	// channels
	progress := make(chan bool)
	defer close(progress)

	// config -- move to helpers/config/config.go
	conf, _ := config.ParseConfig("./config.yaml")
	max_live_instances := conf.Admin.MaxLiveInstances
	instance_interval := conf.Admin.InstanceInterval
	max_total_bots := conf.Admin.MaxTotalBots
	admin_key_hex := get_env_var("PRIVATE_KEY")
	instance_config := strategy.New_Dumb(
		common.HexToAddress(conf.Bot.Address.Router),  // router
		common.HexToAddress(conf.Bot.Address.Factory), // factory
		common.HexToAddress(conf.Bot.Address.Weth),    // weth
		common.HexToAddress(conf.Bot.Address.Token),   // token
		get_env_var("BINANCE_RPC"),
		conf.Bot.Config.ChainId,
		conf.Bot.Config.SwapInterval,
		conf.Bot.Config.MaxSwaps,
	)

	// get admin address from key -- move to helpers/wallet/derive_address.go
	admin_key_ecdsa, err := crypto.HexToECDSA(admin_key_hex)
	if err != nil {
		log.Fatal(err)
	}
	public_key_hex := admin_key_ecdsa.Public()
	public_key_ecdsa, ok := public_key_hex.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	admin_address := crypto.PubkeyToAddress(*public_key_ecdsa)

	bot_counter := 0
	live_instances := 0
	for {
		if live_instances < max_live_instances && bot_counter < max_total_bots {
			go instance_config.Dumb(
				admin_address,
				admin_key_ecdsa,
				&live_instances,
				progress,
				bot_counter,
			)
			<-progress
			live_instances++
			bot_counter++
			time.Sleep(time.Duration(instance_interval) * time.Second)
		} else if live_instances == 0 { // all bots have refunded
			fmt.Println("Exiting gracefully.")
			os.Exit(0)
		} else {
			time.Sleep(time.Duration(instance_interval) * time.Second)
		}
	}
}

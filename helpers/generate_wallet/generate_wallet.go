package generate_wallet

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Generate() (*ecdsa.PrivateKey, common.Address) {
	private_key, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	public_key := private_key.Public()
	public_key_ecdsa, ok := public_key.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address_hex := crypto.PubkeyToAddress(*public_key_ecdsa)
	return private_key, address_hex
}

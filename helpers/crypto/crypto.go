package crypto_

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

func PrivateKeyFromD(d *big.Int) string {
	c := crypto.S256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = d
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(d.Bytes())
	return hex.EncodeToString(crypto.FromECDSA(priv))
}

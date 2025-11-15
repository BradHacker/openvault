package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/BradHacker/openvault/cryptolib"
)

// helper to make deterministic secret key for reproducible test vector
func fixedSecretKey() *cryptolib.SecretKey {
	sk := &cryptolib.SecretKey{}
	sk.Version = cryptolib.LatestSecretKeyVersion
	copy(sk.AccountID[:], []byte("ABCDEF"))
	copy(sk.Secret[:], []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZA")) // 27, truncate to 26
	return sk
}

func randomSalt16() cryptolib.Salt {
	var s cryptolib.Salt
	rand.Read(s[:])
	return s
}

func main() {
	salt, err := cryptolib.NewSalt()
	if err != nil {
		panic(err)
	}
	a, err := cryptolib.DeriveAUK(&cryptolib.AUKParams{
		Email:    "user@example.com",
		Password: "correcthorsebatterystaple",
		Salt:     salt,
		Secret:   fixedSecretKey(),
		Rounds:   100000,
	})
	fmt.Printf("%+v\n", a)
	if err != nil {
		panic(err)
	}
	auk, err := cryptolib.NewKey(cryptolib.AccountUnlockKeyID, a, cryptolib.KeyUseEncryption)
	if err != nil {
		panic(err)
	}
	keyJson, _ := auk.MarshalJSON()
	println("Derived AUK Key:")
	fmt.Printf("%s\n", keyJson)
	wrapped, err := auk.Wrap(auk)
	if err != nil {
		panic(err)
	}
	wrappedKeyJson, _ := json.Marshal(wrapped)
	println("Derived AUK Key (wrapped):")
	fmt.Printf("%s\n", wrappedKeyJson)

	// ks, err := cryptolib.GenerateKeySet(auk)
	// if err != nil {
	// 	panic(err)
	// }

	// ksJson, err := json.Marshal(ks)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Key Set:\n%s\n", ksJson)
}

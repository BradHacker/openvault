package cryptolib

import (
	"fmt"
	"testing"
)

func fixedAUK() *JWK {
	k, err := NewKey(AccountUnlockKeyID, []byte{137, 128, 127, 106, 182, 58, 62, 101, 34, 109, 135, 181, 175, 218, 234, 51, 151, 11, 135, 165, 45, 89, 31, 64, 161, 97, 19, 19, 71, 148, 85, 211}, KeyUseEncryption)
	if err != nil {
		panic(fmt.Sprintf("failed to create AUK: %v", err))
	}
	return k
}

func fixedKeySet(auk *JWK) *KeySet {
	ks, err := GenerateKeySet(auk, &Salt{1, 2, 3, 4, 5, 6, 7, 8}, 650000)
	if err != nil {
		panic(fmt.Sprintf("failed to create KeySet: %v", err))
	}
	return ks
}

func TestDecryptSymmetricKey(t *testing.T) {
	auk := fixedAUK()
	ks := fixedKeySet(auk)
	_, err := ks.SymmetricKey(auk)
	if err != nil {
		t.Fatalf("failed to decrypt symmetric key: %v", err)
	}
}

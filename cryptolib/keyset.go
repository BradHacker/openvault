package cryptolib

import (
	"crypto/elliptic"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	AES_BYTES   int            = 32
	RSA_BITS    int            = 2048
	ECDSA_CURVE elliptic.Curve = elliptic.P521()
)

var (
	ErrInvalidAUK = errors.New("invalid account unlock key")
)

type KeySet struct {
	// Unique identifier for the key set
	ID string `json:"id"`
	// Master Key encrypted with AUK (used to encrypt/decrypt private + signing key)
	EncSymKey *JWE `json:"enc_sym_key"`
	// Public Encryption Key (used to encrypt vault keys)
	PubKey *JWK `json:"pub_key"`
	// Private Encryption Key encrypted with Master Key (used to decrypt vault keys)
	EncPriKey *JWE `json:"enc_pri_key"`
	// Public Signing Key (reserved for future use)
	PubSignKey *JWK `json:"pub_sign_key,omitempty"`
	// Private Signing Key encrypted with Master Key (reserved for future use)
	EncSignKey *JWE `json:"enc_sign_key,omitempty"`
}

// SymmetricKey unwraps the key set symmetric key using the account unlock key (AUK).
//
// The key set symmetric key (a.k.a. Master Key) is used to encrypt/decrypt the private and signing key
func (ks *KeySet) SymmetricKey(accountUnlockKey *JWK) (*JWK, error) {
	if accountUnlockKey.KeyID != AccountUnlockKeyID {
		return nil, fmt.Errorf("%w: invalid AUK ID", ErrInvalidAUK)
	}
	symKey, err := ks.EncSymKey.Unwrap(accountUnlockKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap symmetric key: %w", err)
	}
	return symKey, nil
}

// PrivateKey unwraps the key set private key using the account unlock key (AUK).
//
// The key set private key is used to decrypt vault keys.
func (ks *KeySet) PrivateKey(accountUnlockKey *JWK) (*JWK, error) {
	symKey, err := ks.EncSymKey.Unwrap(accountUnlockKey)
	if err != nil {
		return nil, err
	}
	privKey, err := ks.EncPriKey.Unwrap(symKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap symmetric key: %w", err)
	}
	return privKey, nil
}

// SigningKey unwraps the key set signing key using the account unlock key (AUK).
//
// The key set signing key is reserved for future use.
func (ks *KeySet) SigningKey(accountUnlockKey *JWK) (*JWK, error) {
	symKey, err := ks.EncSymKey.Unwrap(accountUnlockKey)
	if err != nil {
		return nil, err
	}
	signKey, err := ks.EncSignKey.Unwrap(symKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap symmetric key: %w", err)
	}
	return signKey, nil
}

func GenerateKeySet(accountUnlockKey *JWK, aukSalt *Salt, aukRounds int) (*KeySet, error) {
	ks := &KeySet{
		ID: uuid.New().String(),
	}
	symKey, err := generateSymmetricKey(AES_BYTES)
	if err != nil {
		return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	privKey, pubKey, err := generatePrivateKey(RSA_BITS)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	privSignKey, pubSignKey, err := generateSigningKey(ECDSA_CURVE)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signing key: %w", err)
	}
	ks.EncSymKey, err = symKey.Wrap(accountUnlockKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt symmetric key: %w", err)
	}
	// Add PBKDF2 parameters to symmetric key headers (used for AUK derivation)
	ks.EncSymKey.P2Salt = aukSalt[:]
	ks.EncSymKey.P2Rounds = &aukRounds
	ks.EncPriKey, err = privKey.Wrap(symKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}
	ks.PubKey = pubKey
	ks.EncSignKey, err = privSignKey.Wrap(symKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private signing key: %w", err)
	}
	ks.PubSignKey = pubSignKey
	return ks, nil
}

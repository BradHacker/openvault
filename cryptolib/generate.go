package cryptolib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/google/uuid"
)

func generateSymmetricKey(bytes int) (*JWK, error) {
	data := make([]byte, bytes)
	_, err := rand.Read(data[:])
	if err != nil {
		return nil, fmt.Errorf("failed to read random key bytes: %w", err)
	}
	return NewKey(uuid.New().String(), data, KeyUseEncryption)
}

func generatePrivateKey(bytes int) (privateKey *JWK, publicKey *JWK, err error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	// Public and private keys must have the same key ID for encryption and decryption to work properly.
	keyId := uuid.New().String()
	privateKey, err = NewKey(keyId, privKey, KeyUseEncryption)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err = NewKey(keyId, &privKey.PublicKey, KeyUseEncryption)
	if err != nil {
		return nil, nil, err
	}
	return
}

func generateSigningKey(curve elliptic.Curve) (privateKey *JWK, publicKey *JWK, err error) {
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	// Public and private keys must have the same key ID for encryption and decryption to work properly.
	keyId := uuid.New().String()
	privateKey, err = NewKey(keyId, privKey, KeyUseSignature)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err = NewKey(keyId, &privKey.PublicKey, KeyUseSignature)
	if err != nil {
		return nil, nil, err
	}
	return
}

package cryptolib

import (
	"crypto/hkdf"
	"crypto/pbkdf2"
	"crypto/sha256"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/text/unicode/norm"
)

var (
	hash         = sha256.New
	aukAlgorithm = "PBES2g-HS256"
)

type AUKParams struct {
	Email    string
	Password string
	Salt     *Salt
	Secret   *SecretKey
	Rounds   int
}

func DeriveAUK(params *AUKParams) (auk *JWK, err error) {
	if err := validateParams(params); err != nil {
		return nil, err
	}

	lowerEmail := strings.ToLower(params.Email)

	// Password Preprocessing (8.2.2)
	strippedPass := strings.TrimSpace(params.Password)
	normalizedPass := norm.NFKD.Bytes([]byte(strippedPass))

	// Preparing the Salt (8.2.3)
	emailSaltedSalt, err := hkdf.Extract(hash, params.Salt[:], []byte(lowerEmail))
	if err != nil {
		return nil, fmt.Errorf("failed to extract email salt: %w", err)
	}
	expandedSalt, err := hkdf.Expand(hash, emailSaltedSalt, aukAlgorithm, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to expand email salt with HKDF: %w", err)
	}

	// Slow Hashing (8.2.4)
	pKey, err := pbkdf2.Key(hash, string(normalizedPass), expandedSalt, params.Rounds, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key with PBKDF2: %w", err)
	}

	// Combining with the Secret Key (8.2.5)
	acctSaltedSecret, err := hkdf.Extract(hash, params.Secret.Secret[:], params.Secret.AccountID[:])
	if err != nil {
		return nil, fmt.Errorf("failed to extract secret salt with HKDF: %w", err)
	}
	sKey, err := hkdf.Expand(hash, acctSaltedSecret, string(params.Secret.Version[:]), 32)
	if err != nil {
		return nil, fmt.Errorf("failed to expand secret salt with HKDF: %w", err)
	}

	data := make([]byte, 32)
	for i := range data {
		data[i] = pKey[i] ^ sKey[i]
	}

	return NewKey(AccountUnlockKeyID, data, KeyUseEncryption)
}

func validateParams(params *AUKParams) error {
	if params.Rounds <= 0 {
		return fmt.Errorf("rounds must be > 0")
	}
	if _, err := mail.ParseAddress(params.Email); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}
	if len(params.Password) == 0 {
		return fmt.Errorf("password cannot be empty")
	}
	if params.Salt == nil {
		return fmt.Errorf("salt is required")
	}
	if params.Secret == nil {
		return fmt.Errorf("secret key is required")
	}
	return nil
}

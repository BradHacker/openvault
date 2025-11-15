package structs

import (
	"fmt"

	"github.com/BradHacker/openvault/cryptolib"
)

type Account struct {
	ID        string               `json:"id"`
	Email     string               `json:"user_email"`
	FirstName string               `json:"user_first_name"`
	LastName  string               `json:"user_last_name"`
	SecretKey *cryptolib.SecretKey `json:"secret_key"`
}

// TryUnlock attempts to unlock the account using the provided password and symmetric key.
// If successful, it returns the derived Account Unlock Key (AUK).
func (a *Account) TryUnlock(password string, encSymKey *cryptolib.JWE) (auk *cryptolib.JWK, err error) {
	fmt.Printf("Trying to unlock account %s with email %s and secret key %s\n", a.ID, a.Email, a.SecretKey)
	if encSymKey.P2Salt == nil {
		return nil, fmt.Errorf("missing p2s parameter in symmetric key headers")
	}
	if encSymKey.P2Rounds == nil {
		return nil, fmt.Errorf("missing p2c parameter in symmetric key headers")
	}
	fmt.Printf("Extracted salt: %x\n", encSymKey.P2Salt)
	fmt.Printf("Extracted rounds: %d\n", *encSymKey.P2Rounds)
	salt := cryptolib.Salt(encSymKey.P2Salt)
	fmt.Printf("Using salt (len: %d): %x\n", len(salt), salt)
	// Derive the AUK
	aukParams := &cryptolib.AUKParams{
		Email:    a.Email,
		Password: password,
		Salt:     &salt,
		Rounds:   *encSymKey.P2Rounds,
		Secret:   a.SecretKey,
	}
	auk, err = cryptolib.DeriveAUK(aukParams)
	if err != nil {
		return nil, fmt.Errorf("failed to derive AUK: %w", err)
	}
	// Try to unwrap the symmetric key to verify the AUK is correct
	k, err := encSymKey.Unwrap(auk)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap symmetric key with derived AUK: %w", err)
	}
	k.Close()
	return auk, nil
}

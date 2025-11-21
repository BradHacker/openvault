package main

import (
	"fmt"

	"github.com/BradHacker/openvault/openvault/internal/fs"
	"github.com/BradHacker/openvault/openvault/internal/structs"

	"github.com/BradHacker/openvault/cryptolib"
)

type State struct {
	// Whether the application has been initialized with at least one account
	IsInitialized bool
	// Accounts mapped by their IDs
	Accounts fs.AccountStore
	// Keysets mapped by their associated account IDs
	KeySets fs.KeySetStore
	// Vaults mapped by their associated account IDs
	Vaults fs.VaultStore
	// The current set of decrypted vault item overviews
	ItemOverviews fs.ItemOverviewsStore
	// The current decrypted vault item details (only one at a time)
	ItemDetails fs.ItemDetailsStore
	// The account unlock keys for each account
	AUK map[string]*cryptolib.JWK
}

func (s *State) LookupVaultCrypto(vaultId string) (keySet *cryptolib.KeySet, auk *cryptolib.JWK, vault *structs.Vault, err error) {
	// Determine which account the vault belongs to
	vault, ok := s.Vaults[vaultId]
	if !ok {
		return nil, nil, nil, fmt.Errorf("vault %s not found", vaultId)
	}
	auk, ok = s.AUK[vault.AccountID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("account %q is locked", vault.AccountID)
	}
	keySet, ok = s.KeySets[vault.AccountID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("no keyset found for active account %q", vault.AccountID)
	}
	return keySet, auk, vault, nil
}

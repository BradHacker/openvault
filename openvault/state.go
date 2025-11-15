package main

import (
	"fmt"

	"github.com/BradHacker/openvault/openvault/internal/fs"
	"github.com/BradHacker/openvault/openvault/internal/structs"

	"github.com/BradHacker/openvault/openvault/cryptolib"
)

type State struct {
	// Whether the application has been initialized with at least one account
	IsInitialized bool
	// Accounts mapped by their IDs
	Accounts fs.AccountStore
	// Keysets mapped by their associated account IDs
	KeySets fs.KeySetStore
	// Vaults mapped by their IDs
	Vaults fs.VaultStore
	// The current set of decrypted vault item overviews
	ItemOverviews fs.ItemOverviewsStore
	// The current decrypted vault item details (only one at a time)
	ItemDetails fs.ItemDetailsStore
	// The account unlock keys for each account
	AUK map[string]*cryptolib.JWK
}

func (s *State) LookupVaultID(vaultId string) (accountId string, auk *cryptolib.JWK, vault *structs.Vault, err error) {
	// Determine which account the vault belongs to
	for accId, vaults := range s.Vaults {
		for _, v := range vaults {
			if v.ID == vaultId {
				accountId = accId
				vault = v
				break
			}
		}
	}
	if vault == nil {
		return "", nil, nil, fmt.Errorf("vault %s not found", vaultId)
	}
	auk, ok := s.AUK[accountId]
	if !ok {
		return "", nil, nil, fmt.Errorf("account %s is locked", accountId)
	}
	return accountId, auk, vault, nil
}

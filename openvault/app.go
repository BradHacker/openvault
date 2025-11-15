package main

import (
	"context"
	"fmt"

	"github.com/BradHacker/openvault/openvault/internal/fs"
	"github.com/BradHacker/openvault/openvault/internal/structs"

	"github.com/BradHacker/openvault/cryptolib"
	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	ctx   context.Context
	state *State
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		state: &State{
			IsInitialized: false,
			Accounts:      make(fs.AccountStore),
			KeySets:       make(fs.KeySetStore),
			Vaults:        make(fs.VaultStore),
			ItemOverviews: make(fs.ItemOverviewsStore),
			ItemDetails:   make(fs.ItemDetailsStore),
			AUK:           make(map[string]*cryptolib.JWK),
		},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Check if the application is initialized
	a.state.IsInitialized = fs.IsInitializesd()
	if a.state.IsInitialized {
		var err error
		// Load accounts
		a.state.Accounts, err = fs.LoadAccounts()
		if err != nil {
			fmt.Println("Error loading accounts:", err)
			return
		}
		// Load keysets
		a.state.KeySets, err = fs.LoadKeySets()
		if err != nil {
			fmt.Println("Error loading keysets:", err)
			return
		}
		// Load vaults
		a.state.Vaults, err = fs.LoadVaults()
		if err != nil {
			fmt.Println("Error loading vaults:", err)
			return
		}
		// Load item overviews
		a.state.ItemOverviews, err = fs.LoadItemOverviews()
		if err != nil {
			fmt.Println("Error loading item overviews:", err)
			return
		}
		// Load item details
		a.state.ItemDetails, err = fs.LoadItemDetails()
		if err != nil {
			fmt.Println("Error loading item details:", err)
			return
		}
	}

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// IsInitialized returns whether the application has been initialized
func (a *App) IsInitialized() bool {
	a.state.IsInitialized = fs.IsInitializesd()
	return a.state.IsInitialized
}

// Initialize initializes the application with the given options. If the application
// is already initialized, it does nothing.
func (a *App) Initialize(opts fs.InitOptions) error {
	if a.state.IsInitialized {
		return nil
	}
	var err error
	a.state.Accounts, a.state.KeySets, a.state.Vaults, a.state.ItemOverviews, a.state.ItemDetails, err = fs.InitializeAccount(&opts)
	if err != nil {
		return fmt.Errorf("failed to initialize account: %w", err)
	}
	return nil
}

func (a *App) IsLocked() bool {
	return len(a.state.AUK) == 0
}

func (a *App) Lock() error {
	for _, auk := range a.state.AUK {
		auk.Close()
	}
	a.state.AUK = make(map[string]*cryptolib.JWK)
	return nil
}

func (a *App) TryUnlock(password string) error {
	for _, account := range a.state.Accounts {
		// Get the keyset for the account
		keySet, ok := a.state.KeySets[account.ID]
		if !ok {
			return fmt.Errorf("no keyset found for account %s", account.ID)
		}
		fmt.Printf("%+v\n", keySet.EncSymKey)
		// Try to unlock the account
		auk, err := account.TryUnlock(password, keySet.EncSymKey)
		if err == nil {
			// If it unlocks, the AUK matches
			a.state.AUK[account.ID] = auk
			logrus.Printf("Successfully unlocked account %s", account.ID)
			return nil
		}
		logrus.Printf("Account %s did not unlock: %v", account.ID, err)
	}
	return fmt.Errorf("password did not match any account")
}

type AccountWithUnlockStatus struct {
	*structs.Account
	IsUnlocked bool `json:"is_unlocked"`
}

// GetAccounts returns the accounts for the application.
func (a *App) GetAccounts() ([]*AccountWithUnlockStatus, error) {
	if !a.state.IsInitialized {
		return nil, fmt.Errorf("application not initialized")
	}
	var accounts []*AccountWithUnlockStatus
	for _, account := range a.state.Accounts {
		accounts = append(accounts, &AccountWithUnlockStatus{
			Account:    account,
			IsUnlocked: a.state.AUK[account.ID] != nil,
		})
	}
	return accounts, nil
}

// GetVaultMetadatas returns the vault metadata for the given account IDs.
func (a *App) GetVaultMetadatas(accountIds []string) ([]*structs.VaultMetadata, error) {
	if a.IsLocked() {
		return nil, fmt.Errorf("application not unlocked")
	}
	var vaultMetadatas []*structs.VaultMetadata
	for _, accountId := range accountIds {
		auk, ok := a.state.AUK[accountId]
		if !ok {
			return nil, fmt.Errorf("account %s is locked", accountId)
		}
		// Decrypt the private key from the keyset
		keySet, ok := a.state.KeySets[accountId]
		if !ok {
			return nil, fmt.Errorf("no keyset found for active account %s", accountId)
		}
		privKey, err := keySet.PrivateKey(auk)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt private key: %w", err)
		}
		defer privKey.Close()
		for _, vault := range a.state.Vaults[accountId] {
			meta, err := vault.DecryptMetadata(privKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt vault metadata for vault %s: %w", vault.ID, err)
			}
			vaultMetadatas = append(vaultMetadatas, meta)
		}
	}
	return vaultMetadatas, nil
}

type DecryptedVaultItemOverview struct {
	*structs.EncryptedVaultItemOverview
	*structs.VaultItemOverview
}

func (a *App) GetVaultItemOverviews(vaultId string) ([]*DecryptedVaultItemOverview, error) {
	if a.IsLocked() {
		return nil, fmt.Errorf("application not unlocked")
	}
	// Determine which account the vault belongs to
	var accountId string
	var vault *structs.Vault
	for accId, vaults := range a.state.Vaults {
		for _, v := range vaults {
			if v.ID == vaultId {
				accountId = accId
				vault = v
				break
			}
		}
	}
	if vault == nil {
		return nil, fmt.Errorf("vault %s not found", vaultId)
	}
	auk, ok := a.state.AUK[accountId]
	if !ok {
		return nil, fmt.Errorf("account %s is locked", accountId)
	}
	// Decrypt the private key from the keyset
	keySet, ok := a.state.KeySets[accountId]
	if !ok {
		return nil, fmt.Errorf("no keyset found for active account %s", accountId)
	}
	privKey, err := keySet.PrivateKey(auk)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}
	defer privKey.Close()
	// Get the encrypted overviews for the vault
	encOverviews, ok := a.state.ItemOverviews[vaultId]
	if !ok {
		return nil, fmt.Errorf("no item overviews found for vault %s", vaultId)
	}
	// Decrypt the overviews
	overviews, err := vault.DecryptItemOverviews(privKey, encOverviews...)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt item overviews for vault %s: %w", vaultId, err)
	}
	var decryptedOverviews []*DecryptedVaultItemOverview
	for i, ov := range overviews {
		decryptedOverviews = append(decryptedOverviews, &DecryptedVaultItemOverview{
			EncryptedVaultItemOverview: encOverviews[i],
			VaultItemOverview:          ov,
		})
	}
	return decryptedOverviews, nil
}

func (a *App) GetAllVaultItemOverviews() ([]*DecryptedVaultItemOverview, error) {
	if a.IsLocked() {
		return nil, fmt.Errorf("application not unlocked")
	}
	var overviews []*DecryptedVaultItemOverview
	for accountId, vaults := range a.state.Vaults {
		auk, ok := a.state.AUK[accountId]
		if !ok {
			return nil, fmt.Errorf("account %s is locked", accountId)
		}
		// Decrypt the private key from the keyset
		keySet, ok := a.state.KeySets[accountId]
		if !ok {
			return nil, fmt.Errorf("no keyset found for active account %s", accountId)
		}
		privKey, err := keySet.PrivateKey(auk)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt private key: %w", err)
		}
		defer privKey.Close()
		for _, vault := range vaults {
			encOverviews, ok := a.state.ItemOverviews[vault.ID]
			if !ok {
				return nil, fmt.Errorf("no item overviews found for vault %s", vault.ID)
			}
			// Decrypt the overviews
			itemOverviews, err := vault.DecryptItemOverviews(privKey, encOverviews...)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt item overviews for vault %s: %w", vault.ID, err)
			}
			for i, ov := range itemOverviews {
				overviews = append(overviews, &DecryptedVaultItemOverview{
					EncryptedVaultItemOverview: encOverviews[i],
					VaultItemOverview:          ov,
				})
			}
		}
	}
	return overviews, nil
}

type DecryptedVaultItemDetails struct {
	*structs.EncryptedVaultItemDetails
	*structs.VaultItemDetails
}

func (a *App) GetVaultItemDetails(vaultId string, itemId string) (*DecryptedVaultItemDetails, error) {
	accountId, auk, vault, err := a.state.LookupVaultID(vaultId)
	if err != nil {
		return nil, err
	}
	// Decrypt the private key from the keyset
	keySet, ok := a.state.KeySets[accountId]
	if !ok {
		return nil, fmt.Errorf("no keyset found for active account %s", accountId)
	}
	privKey, err := keySet.PrivateKey(auk)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}
	defer privKey.Close()
	// Get the encrypted details for the item
	vaultEncDetails, ok := a.state.ItemDetails[vaultId]
	if !ok {
		return nil, fmt.Errorf("no item details found for item %s", itemId)
	}
	var encDetails *structs.EncryptedVaultItemDetails
	for _, ed := range vaultEncDetails {
		if ed.ItemID == itemId {
			encDetails = ed
			break
		}
	}
	if encDetails == nil {
		return nil, fmt.Errorf("no item details found for item %s", itemId)
	}
	// Decrypt the details
	details, err := vault.DecryptItemDetails(privKey, encDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt item details for item %s: %w", itemId, err)
	}
	return &DecryptedVaultItemDetails{
		EncryptedVaultItemDetails: encDetails,
		VaultItemDetails:          details,
	}, nil
}

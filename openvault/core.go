package main

import (
	"fmt"
	"slices"

	"github.com/BradHacker/openvault/openvault/internal/fs"
	"github.com/BradHacker/openvault/openvault/internal/structs"

	"github.com/BradHacker/openvault/cryptolib"
	"github.com/sirupsen/logrus"
)

// CoreService struct
type CoreService struct {
	state *State
}

// NewCoreService creates a new CoreService struct
func NewCoreService() *CoreService {
	core := &CoreService{
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
	core.startup()
	return core
}

func (a *CoreService) startup() {
	// Check if the application is initialized
	a.state.IsInitialized = fs.IsInitialized()
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

// IsInitialized returns whether the application has been initialized
func (a *CoreService) IsInitialized() bool {
	a.state.IsInitialized = fs.IsInitialized()
	return a.state.IsInitialized
}

// Initialize initializes the application with the given options. If the application
// is already initialized, it does nothing.
func (a *CoreService) Initialize(opts fs.InitOptions) error {
	if a.state.IsInitialized {
		return nil
	}
	var err error
	a.state.Accounts, a.state.KeySets, a.state.Vaults, a.state.ItemOverviews, a.state.ItemDetails, err = fs.InitializeAccount(&opts)
	if err != nil {
		return fmt.Errorf("failed to initialize account: %w", err)
	}
	a.state.IsInitialized = true
	return nil
}

func (a *CoreService) IsLocked() bool {
	return len(a.state.AUK) == 0
}

func (a *CoreService) Lock() error {
	for _, auk := range a.state.AUK {
		auk.Close()
	}
	a.state.AUK = make(map[string]*cryptolib.JWK)
	return nil
}

func (a *CoreService) TryUnlock(password string) error {
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
func (a *CoreService) GetAccounts() ([]*AccountWithUnlockStatus, error) {
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

func (a *CoreService) GetAccount(accountId string) (*AccountWithUnlockStatus, error) {
	if !a.state.IsInitialized {
		return nil, fmt.Errorf("application not initialized")
	}
	account, ok := a.state.Accounts[accountId]
	if !ok {
		return nil, fmt.Errorf("account %s not found", accountId)
	}
	return &AccountWithUnlockStatus{
		Account:    account,
		IsUnlocked: a.state.AUK[account.ID] != nil,
	}, nil
}

// GetVaultMetadatas returns the vault metadata for the given account IDs.
func (a *CoreService) ListVaultMetadatas(accountIds []string) ([]*structs.VaultMetadata, error) {
	if a.IsLocked() {
		return nil, fmt.Errorf("application not unlocked")
	}
	vaultByAccount := make(map[string][]*structs.Vault)
	for _, vault := range a.state.Vaults {
		if vaultByAccount[vault.AccountID] == nil {
			vaultByAccount[vault.AccountID] = []*structs.Vault{}
		}
		vaultByAccount[vault.AccountID] = append(vaultByAccount[vault.AccountID], vault)
	}

	var vaultMetadatas []*structs.VaultMetadata
	for accountId, vaults := range vaultByAccount {
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
			meta, err := vault.DecryptMetadata(privKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt vault metadata for vault %s: %w", vault.VaultID, err)
			}
			vaultMetadatas = append(vaultMetadatas, meta)
		}
	}
	return vaultMetadatas, nil
}

// GetVaultMetadata returns the vault metadata for the given vault ID.
func (a *CoreService) GetVaultMetadata(vaultId string) (*structs.VaultMetadata, error) {
	if a.IsLocked() {
		return nil, fmt.Errorf("application not unlocked")
	}
	keySet, auk, vault, err := a.state.LookupVaultCrypto(vaultId)
	if err != nil {
		return nil, err
	}
	privKey, err := keySet.PrivateKey(auk)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}
	defer privKey.Close()
	meta, err := vault.DecryptMetadata(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault metadata for vault %s: %w", vault.VaultID, err)
	}
	return meta, nil
}

type DecryptedVaultItemOverview struct {
	*structs.EncryptedVaultItemOverview
	*structs.VaultItemOverview
}

func (a *CoreService) ListVaultItemOverviews(vaultId string) ([]*DecryptedVaultItemOverview, error) {
	if a.IsLocked() {
		return nil, fmt.Errorf("application not unlocked")
	}
	encItemOverviews := slices.Collect(func(yield func(*structs.EncryptedVaultItemOverview) bool) {
		for _, encOverview := range a.state.ItemOverviews {
			if encOverview.VaultID == vaultId {
				if !yield(encOverview) {
					break
				}
			}
		}
	})
	keySet, auk, vault, err := a.state.LookupVaultCrypto(vaultId)
	if err != nil {
		return nil, err
	}
	privKey, err := keySet.PrivateKey(auk)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}
	defer privKey.Close()
	// Decrypt the overviews
	overviews, err := vault.DecryptItemOverviews(privKey, encItemOverviews...)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt item overviews for vault %s: %w", vaultId, err)
	}
	var decryptedOverviews []*DecryptedVaultItemOverview
	for i, ov := range overviews {
		decryptedOverviews = append(decryptedOverviews, &DecryptedVaultItemOverview{
			EncryptedVaultItemOverview: encItemOverviews[i],
			VaultItemOverview:          ov,
		})
	}
	return decryptedOverviews, nil
}

func (a *CoreService) GetItemOverview(itemId string) (*DecryptedVaultItemOverview, error) {
	if a.IsLocked() {
		return nil, fmt.Errorf("application not unlocked")
	}
	var encItemOverview *structs.EncryptedVaultItemOverview
	for _, encOverview := range a.state.ItemOverviews {
		if encOverview.ItemID == itemId {
			encItemOverview = encOverview
			break
		}
	}
	if encItemOverview == nil {
		return nil, fmt.Errorf("no item overview found for item %s", itemId)
	}
	keySet, auk, vault, err := a.state.LookupVaultCrypto(encItemOverview.VaultID)
	if err != nil {
		return nil, err
	}
	privKey, err := keySet.PrivateKey(auk)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}
	defer privKey.Close()

	// Decrypt the overviews
	overviews, err := vault.DecryptItemOverviews(privKey, encItemOverview)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt item overview for item %s: %w", itemId, err)
	}
	return &DecryptedVaultItemOverview{
		EncryptedVaultItemOverview: encItemOverview,
		VaultItemOverview:          overviews[0],
	}, nil
}

func (a *CoreService) ListAllItemOverviews() ([]*DecryptedVaultItemOverview, error) {
	if a.IsLocked() {
		return nil, fmt.Errorf("application not unlocked")
	}

	encItemsByVault := make(map[string][]*structs.EncryptedVaultItemOverview)
	for _, encOverview := range a.state.ItemOverviews {
		if encItemsByVault[encOverview.VaultID] == nil {
			encItemsByVault[encOverview.VaultID] = []*structs.EncryptedVaultItemOverview{}
		}
		encItemsByVault[encOverview.VaultID] = append(encItemsByVault[encOverview.VaultID], encOverview)
	}

	var decryptedOverviews []*DecryptedVaultItemOverview
	for vaultId, encItemOverviews := range encItemsByVault {
		keySet, auk, vault, err := a.state.LookupVaultCrypto(vaultId)
		if err != nil {
			return nil, err
		}
		privKey, err := keySet.PrivateKey(auk)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt private key: %w", err)
		}
		defer privKey.Close()
		overviews, err := vault.DecryptItemOverviews(privKey, encItemOverviews...)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt item overviews for vault %s: %w", vaultId, err)
		}
		for i, ov := range overviews {
			decryptedOverviews = append(decryptedOverviews, &DecryptedVaultItemOverview{
				EncryptedVaultItemOverview: encItemOverviews[i],
				VaultItemOverview:          ov,
			})
		}
	}
	return decryptedOverviews, nil
}

type DecryptedVaultItemDetails struct {
	*structs.EncryptedVaultItemDetails
	*structs.VaultItemDetails
}

func (a *CoreService) GetVaultItemDetails(itemId string) (*DecryptedVaultItemDetails, error) {
	// Get the encrypted details for the item
	encItemDetails, ok := a.state.ItemDetails[itemId]
	if !ok {
		return nil, fmt.Errorf("no item details found for item %s", itemId)
	}
	keySet, auk, vault, err := a.state.LookupVaultCrypto(encItemDetails.VaultID)
	if err != nil {
		return nil, err
	}
	privKey, err := keySet.PrivateKey(auk)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}
	defer privKey.Close()
	// Decrypt the details
	details, err := vault.DecryptItemDetails(privKey, encItemDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt item details for item %s: %w", itemId, err)
	}
	return &DecryptedVaultItemDetails{
		EncryptedVaultItemDetails: encItemDetails,
		VaultItemDetails:          details,
	}, nil
}

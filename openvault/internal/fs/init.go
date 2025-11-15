package fs

import (
	"fmt"

	"github.com/BradHacker/openvault/openvault/internal/constants"

	"os"
	"time"

	"github.com/BradHacker/openvault/openvault/internal/structs"

	"github.com/BradHacker/openvault/cryptolib"
	"github.com/google/uuid"
)

var initialized bool = false

func init() {
	// Check if all of the required files and directories exist
	files := []string{
		accountFile,
		keySetFile,
		vaultFile,
		itemOverviewsFile,
		itemDetailsFile,
	}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			initialized = false
			return
		}
	}
	initialized = true
}

func IsInitializesd() bool {
	return initialized
}

type InitOptions struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func InitializeAccount(opts *InitOptions) (AccountStore, KeySetStore, VaultStore, ItemOverviewsStore, ItemDetailsStore, error) {
	// Check that the data directory exists
	if _, err := os.Stat(constants.DATA_DIR); os.IsNotExist(err) {
		if err := os.MkdirAll(constants.DATA_DIR, 0755); err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("failed to create data directory: %w", err)
		}
	}
	// Generate a random secret key
	secretKey, err := cryptolib.NewSecretKey()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to generate secret key: %w", err)
	}
	fmt.Printf("Generated secret key: %s\n", secretKey)
	// Generate a new account
	account := &structs.Account{
		ID:        fmt.Sprintf("%x", secretKey.AccountID),
		Email:     opts.Email,
		FirstName: opts.FirstName,
		LastName:  opts.LastName,
		SecretKey: secretKey,
	}
	// Save the account to the filesystem
	accountStore := make(AccountStore)
	accountStore[account.ID] = account
	if err := SaveAccounts(accountStore); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to save accounts: %w", err)
	}
	// Generate a PBKDF2 salt
	salt, err := cryptolib.NewSalt()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	fmt.Printf("Generated salt (len: %d): %x\n", len(salt), salt)
	// Derive the Account Unlock Key (AUK)
	aukParams := &cryptolib.AUKParams{
		Email:    opts.Email,
		Password: opts.Password,
		Salt:     salt,
		Secret:   account.SecretKey,
		Rounds:   constants.PBKDF2_ROUNDS,
	}
	auk, err := cryptolib.DeriveAUK(aukParams)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to derive AUK: %w", err)
	}

	// Create an initial KeySet
	ks, err := cryptolib.GenerateKeySet(auk, salt, constants.PBKDF2_ROUNDS)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to generate keyset: %w", err)
	}

	// Save the keyset to the filesystem
	keySetStore := make(KeySetStore)
	keySetStore[account.ID] = ks
	if err := SaveKeySets(keySetStore); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to save keysets: %w", err)
	}

	// Create a new key for the default vault
	vaultKey, err := cryptolib.GenerateVaultKey()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to generate vault key: %w", err)
	}
	defer vaultKey.Close()

	// Create a new default vault
	vaultId := uuid.New().String()
	vaultMeta := &structs.VaultMetadata{
		VaultID:     vaultId,
		Name:        "Default",
		Description: "Welcome to OpenVault!",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	encVaultMeta, err := vaultMeta.Encrypt(vaultKey)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to encrypt vault metadata: %w", err)
	}
	encVaultKey, err := vaultKey.Wrap(ks.PubKey)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to wrap vault key: %w", err)
	}
	vault := &structs.Vault{
		ID:                vaultId,
		EncryptedMetadata: encVaultMeta,
		EncryptedVaultKey: encVaultKey,
	}
	// Save the vault to the filesystem
	vaultStore := make(VaultStore)
	vaultStore[account.ID] = []*structs.Vault{vault}
	if err := SaveVaults(vaultStore); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to save vaults: %w", err)
	}

	// Create a sample item in the default vault
	itemId := uuid.New().String()
	itemOverview := &structs.EncryptedVaultItemOverview{
		ItemID:    itemId,
		VaultID:   vault.ID,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	err = itemOverview.Update(vaultKey, &structs.VaultItemOverview{
		Title: fmt.Sprintf("OpenVault (%s)", account.Email),
		URL:   "https://openvault.io",
	})
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to encrypt item overview: %w", err)
	}
	itemDetails := &structs.EncryptedVaultItemDetails{
		ItemID:    itemId,
		VaultID:   vault.ID,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	err = itemDetails.Update(vaultKey, &structs.VaultItemDetails{
		Username: opts.Email,
		Password: opts.Password,
		Notes:    "**Welcome to OpenVault!** This is your first item. It has your OpenVault login details.",
	})
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to encrypt item details: %w", err)
	}

	// Save the item to the filesystem
	overviewStore := make(ItemOverviewsStore)
	overviewStore[vault.ID] = []*structs.EncryptedVaultItemOverview{itemOverview}
	if err := SaveItemOverviews(overviewStore); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to save overviews: %w", err)
	}
	detailsStore := make(ItemDetailsStore)
	detailsStore[vault.ID] = []*structs.EncryptedVaultItemDetails{itemDetails}
	if err := SaveItemDetails(detailsStore); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to save details: %w", err)
	}

	return accountStore, keySetStore, vaultStore, overviewStore, detailsStore, nil
}

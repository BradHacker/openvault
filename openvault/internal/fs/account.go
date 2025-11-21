package fs

import (
	"path"

	"github.com/BradHacker/openvault/openvault/internal/constants"
	"github.com/BradHacker/openvault/openvault/internal/structs"
)

var accountFile = path.Join(constants.DATA_DIR, "accounts.json")

// AccountStore is a map of accounts by their IDs
type AccountStore map[string]*structs.Account

// AccountsExists checks if the accounts file exists
func AccountsExists() bool {
	return exists(accountFile)
}

// LoadAccounts loads the accounts from the filesystem
func LoadAccounts() (AccountStore, error) {
	var as AccountStore
	if err := load(accountFile, &as); err != nil {
		return nil, err
	}
	return as, nil
}

// SaveAccounts saves the accounts to the filesystem
func SaveAccounts(as AccountStore) error {
	return save(accountFile, as)
}

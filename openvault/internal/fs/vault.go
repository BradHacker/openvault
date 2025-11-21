package fs

import (
	"path"

	"github.com/BradHacker/openvault/openvault/internal/constants"
	"github.com/BradHacker/openvault/openvault/internal/structs"
)

var vaultFile = path.Join(constants.DATA_DIR, "vaults.json")

// VaultStore is a map of vaults by their associated account IDs
type VaultStore map[string]*structs.Vault

// VaultsExists checks if the vaults file exists
func VaultsExists() bool {
	return exists(vaultFile)
}

// LoadVaults loads the vaults from the filesystem
func LoadVaults() (VaultStore, error) {
	var vs VaultStore
	if err := load(vaultFile, &vs); err != nil {
		return nil, err
	}
	return vs, nil
}

// SaveVaults saves the vaults to the filesystem
func SaveVaults(vs VaultStore) error {
	return save(vaultFile, vs)
}

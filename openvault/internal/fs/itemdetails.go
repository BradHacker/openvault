package fs

import (
	"path"

	"github.com/BradHacker/openvault/openvault/internal/constants"
	"github.com/BradHacker/openvault/openvault/internal/structs"
)

var itemDetailsFile = path.Join(constants.DATA_DIR, "item_details.json")

// ItemDetailsStore is a map of item details by their associated vault IDs
type ItemDetailsStore map[string][]*structs.EncryptedVaultItemDetails

// ItemDetailsExists checks if the item details file exists
func ItemDetailsExists() bool {
	return exists(itemDetailsFile)
}

// LoadItemDetails loads the item details from the filesystem
func LoadItemDetails() (ItemDetailsStore, error) {
	var ids ItemDetailsStore
	if err := load(itemDetailsFile, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// SaveItemDetails saves the item details to the filesystem
func SaveItemDetails(ids ItemDetailsStore) error {
	return save(itemDetailsFile, ids)
}

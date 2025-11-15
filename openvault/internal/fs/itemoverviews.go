package fs

import (
	"path"

	"github.com/BradHacker/openvault/openvault/internal/constants"
	"github.com/BradHacker/openvault/openvault/internal/structs"
)

var itemOverviewsFile = path.Join(constants.DATA_DIR, "item_overviews.json")

// ItemOverviewsStore is a map of item overviews by their associated vault IDs
type ItemOverviewsStore map[string][]*structs.EncryptedVaultItemOverview

// ItemOverviewsExists checks if the item overviews file exists
func ItemOverviewsExists() bool {
	return exists(itemOverviewsFile)
}

// LoadItemOverviews loads the item overviews from the filesystem
func LoadItemOverviews() (ItemOverviewsStore, error) {
	var ios ItemOverviewsStore
	if err := load(itemOverviewsFile, &ios); err != nil {
		return nil, err
	}
	return ios, nil
}

// SaveItemOverviews saves the item overviews to the filesystem
func SaveItemOverviews(ios ItemOverviewsStore) error {
	return save(itemOverviewsFile, ios)
}

package fs

import (
	"path"

	"github.com/BradHacker/openvault/openvault/internal/constants"

	"github.com/BradHacker/openvault/cryptolib"
)

var keySetFile = path.Join(constants.DATA_DIR, "keysets.json")

// KeySetStore is a map of keysets by their associated account IDs
type KeySetStore map[string]*cryptolib.KeySet

// KeySetsExists checks if the keysets file exists
func KeySetsExists() bool {
	return exists(keySetFile)
}

// LoadKeySets loads the keysets from the filesystem
func LoadKeySets() (KeySetStore, error) {
	var ks KeySetStore
	if err := load(keySetFile, &ks); err != nil {
		return nil, err
	}
	return ks, nil
}

// SaveKeySets saves the keysets to the filesystem
func SaveKeySets(ks KeySetStore) error {
	return save(keySetFile, ks)
}

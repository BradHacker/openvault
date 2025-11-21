package structs

import (
	"encoding/json"
	"time"

	"github.com/BradHacker/openvault/cryptolib"
)

type Vault struct {
	VaultID           string         `json:"vault_id"`
	AccountID         string         `json:"account_id"`
	EncryptedMetadata *cryptolib.JWE `json:"encrypted_metadata"`
	EncryptedVaultKey *cryptolib.JWE `json:"encrypted_vault_key"`
}

// DecryptVaultKey decrypts the vault key using the vault's encrypted vault key and the keyset's private key
func (v *Vault) decryptVaultKey(privKey *cryptolib.JWK) (*cryptolib.JWK, error) {
	// Decrypt the vault key using the vault's encrypted vault key and the keyset's private key
	vaultKey, err := v.EncryptedVaultKey.Unwrap(privKey)
	if err != nil {
		return nil, err
	}
	return vaultKey, nil
}

// DecryptMetadata decrypts the vault metadata using the provided private key
func (v *Vault) DecryptMetadata(privKey *cryptolib.JWK) (*VaultMetadata, error) {
	// Decrypt the vault key
	vaultKey, err := v.decryptVaultKey(privKey)
	if err != nil {
		return nil, err
	}
	defer vaultKey.Close()
	// Decrypt the metadata
	metadataBytes, err := vaultKey.Decrypt(v.EncryptedMetadata)
	if err != nil {
		return nil, err
	}
	var metadata VaultMetadata
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

// DecryptItemOverviews decrypts the vault item overviews using the provided private key
func (v *Vault) DecryptItemOverviews(privKey *cryptolib.JWK, encryptedOverviews ...*EncryptedVaultItemOverview) ([]*VaultItemOverview, error) {
	// Decrypt the vault key
	vaultKey, err := v.decryptVaultKey(privKey)
	if err != nil {
		return nil, err
	}
	defer vaultKey.Close()
	// Decrypt each overview
	var overviews []*VaultItemOverview
	for _, encOverview := range encryptedOverviews {
		overviewBytes, err := vaultKey.Decrypt(encOverview.EncryptedOverview)
		if err != nil {
			return nil, err
		}
		var overview VaultItemOverview
		err = json.Unmarshal(overviewBytes, &overview)
		if err != nil {
			return nil, err
		}
		overviews = append(overviews, &overview)
	}
	return overviews, nil
}

// DecryptItemDetails decrypts the vault item details using the provided private key
func (v *Vault) DecryptItemDetails(privKey *cryptolib.JWK, encryptedDetails *EncryptedVaultItemDetails) (*VaultItemDetails, error) {
	// Decrypt the vault key
	vaultKey, err := v.decryptVaultKey(privKey)
	if err != nil {
		return nil, err
	}
	defer vaultKey.Close()
	// Decrypt the item details
	detailsBytes, err := vaultKey.Decrypt(encryptedDetails.EncryptedDetails)
	if err != nil {
		return nil, err
	}
	var details *VaultItemDetails
	err = json.Unmarshal(detailsBytes, &details)
	if err != nil {
		return nil, err
	}
	return details, nil
}

type VaultMetadata struct {
	AccountID   string `json:"account_id"`
	VaultID     string `json:"vault_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	// Members     []string `json:"members"`
}

func (vm *VaultMetadata) Encrypt(vaultKey *cryptolib.JWK) (*cryptolib.JWE, error) {
	encryptedMetadata, err := vaultKey.EncryptJSON(vm)
	if err != nil {
		return nil, err
	}
	return encryptedMetadata, nil
}

type EncryptedVaultItemOverview struct {
	ItemID            string         `json:"item_id"`
	VaultID           string         `json:"vault_id"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
	EncryptedOverview *cryptolib.JWE `json:"encrypted_overview"`
}

type VaultItemOverview struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (vio *EncryptedVaultItemOverview) Update(vaultKey *cryptolib.JWK, data *VaultItemOverview) (err error) {
	vio.EncryptedOverview, err = vaultKey.EncryptJSON(data)
	vio.UpdatedAt = time.Now().Format(time.RFC3339)
	return err
}

func (vio *EncryptedVaultItemOverview) Read(vaultKey *cryptolib.JWK) (data *VaultItemOverview, err error) {
	data = &VaultItemOverview{}
	err = vaultKey.DecryptJSON(vio.EncryptedOverview, &data)
	return data, err
}

type EncryptedVaultItemDetails struct {
	ItemID           string         `json:"item_id"`
	VaultID          string         `json:"vault_id"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
	EncryptedDetails *cryptolib.JWE `json:"encrypted_details"`
}

type VaultItemDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Notes    string `json:"notes"`
}

func (vio *EncryptedVaultItemDetails) Update(vaultKey *cryptolib.JWK, data *VaultItemDetails) (err error) {
	vio.EncryptedDetails, err = vaultKey.EncryptJSON(data)
	vio.UpdatedAt = time.Now().Format(time.RFC3339)
	return err
}

func (vio *EncryptedVaultItemDetails) Read(vaultKey *cryptolib.JWK) (data *VaultItemDetails, err error) {
	data = &VaultItemDetails{}
	err = vaultKey.DecryptJSON(vio.EncryptedDetails, &data)
	return data, err
}

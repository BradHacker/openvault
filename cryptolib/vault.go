package cryptolib

// GenerateVaultKey generates a new random symmetric key for use as a vault key.
//
// The vault key is used to encrypt/decrypt vault data. It should be encrypted with the key set private key and stored in the vault.
func GenerateVaultKey() (*JWK, error) {
  return generateSymmetricKey(AES_BYTES)
}

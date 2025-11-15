package cryptolib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-jose/go-jose/v4"
)

var (
	AccountUnlockKeyID = "auk" // Special key ID for AUK (Account Unlock Key)
)

type JWK struct {
	jose.JSONWebKey
	// Is set once the Key has been zeroed in memory and can no longer be used
	cleared bool
}

// Close will zero the memory used for private portions of this key, rendering it unusable.
// It is a good practice to call Close on keys containing sensitive material when they are no longer needed, such as in a defer statement.
//
// IMPORTANT: This key will no longer be usable for operations and MUST be disposed of.
//
// For public keys, this is a no-op.
func (k *JWK) Close() {
	switch key := k.Key.(type) {
	// Clear RSA private key
	case *rsa.PrivateKey:
		key.D.SetInt64(0)
		clear(key.Primes)
	// Clear ECDSA key
	case *ecdsa.PrivateKey:
		key.D.SetInt64(0)
	// Clear AES symmetric key
	case []byte:
		clear(key)
	default:
		return
	}
	k.cleared = true
}

// IsCleared returns if the key has been previously cleared
func (k *JWK) IsCleared() bool {
	return k.cleared
}

type KeyUse string

var (
	KeyUseSignature  KeyUse = "sig"
	KeyUseEncryption KeyUse = "enc"
)

type ContentType string

var (
	ContentTypeJWK ContentType = "jwk+json"
)

var (
	ErrUnsupportedAlg = errors.New("unsupported algorithm")
	ErrKeyCleared     = errors.New("key has been cleared")
)

// NewKey creates a new Key with the given parameters.
//
// Key is the Go in-memory representation of this key. It must have one
// of these types:
//   - ed25519.PublicKey
//   - ed25519.PrivateKey
//   - *ecdsa.PublicKey
//   - *ecdsa.PrivateKey
//   - *rsa.PublicKey
//   - *rsa.PrivateKey
//   - []byte (a symmetric key)
func NewKey(keyID string, key interface{}, keyUse KeyUse) (*JWK, error) {
	keyAlg, err := keyAlgFromType(key)
	if err != nil {
		return nil, err
	}
	return &JWK{
		JSONWebKey: jose.JSONWebKey{Key: key, KeyID: keyID, Use: string(keyUse), Algorithm: string(keyAlg)},
	}, nil
}

// keyAlgFromType determines the appropriate JOSE algorithm for the given key type.
//
// Only a subset of key types and algorithms are supported.
func keyAlgFromType(key interface{}) (jose.KeyAlgorithm, error) {
	var keyAlg jose.KeyAlgorithm
	// Validate key type
	switch key := key.(type) {
	case []byte:
		// Symmetric key
		if len(key) == 16 {
			keyAlg = jose.KeyAlgorithm("A128GCM")
		} else if len(key) == 24 {
			keyAlg = jose.KeyAlgorithm("A192GCM")
		} else if len(key) == 32 {
			keyAlg = jose.KeyAlgorithm("A256GCM")
		} else {
			return "", fmt.Errorf("%w: unsupported symmetric key length %d", ErrUnsupportedAlg, len(key))
		}
	case *ecdsa.PublicKey, *ecdsa.PrivateKey:
		// Signing key
		keyAlg = jose.KeyAlgorithm("ECDH-ES")
	case *rsa.PublicKey, *rsa.PrivateKey:
		// RSA key
		keyAlg = jose.KeyAlgorithm("RSA-OAEP")
	default:
		return "", fmt.Errorf("%w: unsupported key type %T", ErrUnsupportedAlg, key)
	}
	return keyAlg, nil
}

var (
	ErrInvalidWrap    = errors.New("invalid wrapping configuration")
	ErrAlreadyWrapped = errors.New("key already wrapped")
	ErrNotWrapped     = errors.New("key not wrapped")
)

type JWE struct {
	// ContentType indicates the type of the wrapped content
	ContentType ContentType `json:"cty"`
	// EncryptedData is the encrypted key data
	EncryptedData []byte `json:"data"`
	// IV is the initialization vector or nonce used during encryption (if applicable)
	IV []byte `json:"iv,omitempty"`
	// EncryptionAlg is the algorithm used to encrypt the key
	EncryptionAlg string `json:"enc"`
	// Hint at which key was used to wrap (encrypt) this key
	KeyID string `json:"kid"`
	// Optional header for PBKDF2 salt
	P2Salt []byte `json:"p2s,omitempty"`
	// Optional header for PBKDF2 rounds
	P2Rounds *int `json:"p2c,omitempty"`
}

// Encrypt encrypts the given data using this Key and returns a JWE containing the encrypted data.
//
// The encryption algorithm used depends on the type of key contained in this Key. If the key is a
// symmetric key (a []byte), then AES GCM will be used. If the key is an RSA public key, then
// RSA-OAEP will be used. For any other key type, an error will be returned.
func (k *JWK) Encrypt(data []byte) (jwe *JWE, err error) {
	if k.cleared {
		return nil, ErrKeyCleared
	}
	switch key := k.Key.(type) {
	case *rsa.PublicKey:
		return k.encryptRSA(key, data)
	case []byte:
		return k.encryptAES(key, data)
	default:
		return nil, fmt.Errorf("%w: cannot use algorithm \"%s\" for encrypting", ErrUnsupportedAlg, k.Algorithm)
	}
}

// EncryptJSON is a convenience method that marshals the given value to JSON.
//
// This is equivalent to calling json.Marshal on the value and then calling Encrypt on the resulting bytes.
func (k *JWK) EncryptJSON(v interface{}) (jwe *JWE, err error) {
	if k.cleared {
		return nil, ErrKeyCleared
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return k.Encrypt(data)
}

func (k *JWK) encryptAES(key []byte, data []byte) (*JWE, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}
	ct := gcm.Seal(nil, nil, data, nil)
	iv := ct[:aesNonceSize]
	encData := ct[aesNonceSize:]
	return &JWE{
		ContentType:   ContentTypeJWK,
		KeyID:         k.KeyID,
		EncryptedData: encData,
		IV:            iv,
		EncryptionAlg: fmt.Sprintf("A%dGCM", len(key)*8),
	}, nil
}

func (k *JWK) encryptRSA(key *rsa.PublicKey, data []byte) (*JWE, error) {
	encData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, data, nil)
	if err != nil {
		return nil, err
	}
	return &JWE{
		ContentType: ContentTypeJWK,
		KeyID:       k.KeyID,
		// IV:            nil,
		EncryptedData: encData,
		EncryptionAlg: string(jose.RSA_OAEP),
	}, nil
}

func (k *JWK) Decrypt(jwe *JWE) ([]byte, error) {
	if k.cleared {
		return nil, ErrKeyCleared
	}
	switch key := k.Key.(type) {
	case *rsa.PrivateKey:
		return k.decryptRSA(key, jwe)
	case []byte:
		return k.decryptAES(key, jwe)
	default:
		return nil, fmt.Errorf("%w: cannot use algorithm \"%s\" for decrypting", ErrUnsupportedAlg, k.Algorithm)
	}
}

// EncryptJSON is a convenience method that marshals the given value to JSON.
//
// This is equivalent to calling json.Marshal on the value and then calling Encrypt on the resulting bytes.
func (k *JWK) DecryptJSON(data *JWE, v interface{}) (err error) {
	if k.cleared {
		return ErrKeyCleared
	}
	decData, err := k.Decrypt(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(decData, v)
	if err != nil {
		return err
	}
	return nil
}

func (k *JWK) decryptAES(key []byte, jwe *JWE) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}
	// Prep the combined ciphertext (IV/nonce + raw_ciphertext + auth_tag)
	ct := append(jwe.IV, jwe.EncryptedData...)
	data, err := gcm.Open(nil, nil, ct, nil)
	// data, err := gcm.Open(nil, nil, jwe.EncryptedData, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (k *JWK) decryptRSA(key *rsa.PrivateKey, jwe *JWE) ([]byte, error) {
	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, jwe.EncryptedData, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Wraps this key using wrapKey and returns an encrypted WrappedKey.
//
// wrapKey must utilize one of the following algorithms:
//   - RSA-OAEP
//   - AES GCM
func (k *JWK) Wrap(wrapKey *JWK) (*JWE, error) {
	if k.cleared {
		return nil, ErrKeyCleared
	}
	switch key := wrapKey.Key.(type) {
	case *rsa.PublicKey:
		return k.wrapRSA(key, wrapKey.KeyID)
	case []byte:
		return k.wrapAES(key, wrapKey.KeyID)
	default:
		return nil, fmt.Errorf("%w: cannot use algorithm \"%s\" for key wrapping", ErrUnsupportedAlg, wrapKey.Algorithm)
	}
}

const (
	// Standard nonce size used when random nonce GCM mode is used (see SealWithRandomNonce in crypto/internal/fips140/aes/gcm/gcm_nonces.go)
	aesNonceSize = 12
)

func (k *JWK) wrapAES(key []byte, kid string) (*JWE, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}
	keyBytes, err := k.MarshalJSON()
	if err != nil {
		return nil, err
	}
	ct := gcm.Seal(nil, nil, keyBytes, nil)
	return &JWE{
		ContentType:   ContentTypeJWK,
		KeyID:         kid,
		EncryptedData: ct,
		EncryptionAlg: fmt.Sprintf("A%dGCM", len(key)*8),
	}, nil
}

func (k *JWK) wrapRSA(key *rsa.PublicKey, kid string) (*JWE, error) {
	keyBytes, err := k.MarshalJSON()
	if err != nil {
		return nil, err
	}
	data, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, keyBytes, nil)
	if err != nil {
		return nil, err
	}
	return &JWE{
		ContentType: ContentTypeJWK,
		KeyID:       kid,
		// IV:            nil,
		EncryptedData: data,
		EncryptionAlg: string(jose.RSA_OAEP),
	}, nil
}

// Unwraps this wrapped key using unwrapKey and returns an decrypted Key.
//
// unwrapKey must utilize one of the following algorithms:
//   - RSA-OAEP
//   - AES GCM
func (e *JWE) Unwrap(unwrapKey *JWK) (*JWK, error) {
	if unwrapKey.cleared {
		return nil, ErrKeyCleared
	}
	switch key := unwrapKey.Key.(type) {
	case *rsa.PrivateKey:
		return e.unwrapRSA(key)
	case []byte:
		return e.unwrapAES(key)
	default:
		return nil, fmt.Errorf("%w: cannot use algorithm \"%s\" for key unwrapping", ErrUnsupportedAlg, unwrapKey.Algorithm)
	}
}

func (e *JWE) unwrapAES(key []byte) (*JWK, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}
	data, err := gcm.Open(nil, nil, e.EncryptedData, nil)
	if err != nil {
		return nil, err
	}
	k := JWK{}
	err = k.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (e *JWE) unwrapRSA(key *rsa.PrivateKey) (*JWK, error) {
	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, e.EncryptedData, nil)
	if err != nil {
		return nil, err
	}
	k := JWK{}
	err = k.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

package cryptolib

import (
  "bytes"
  "crypto/rand"
  "crypto/rsa"
  "fmt"
  "testing"

  "github.com/go-jose/go-jose/v4"
  "github.com/google/uuid"
)

func randomAES256Key() []byte {
  keyData := make([]byte, 32)
  _, err := rand.Read(keyData[:])
  if err != nil {
    panic(fmt.Sprintf("failed to read random bytes: %v", err))
  }
  return keyData
}

func randomRSAKey() (*rsa.PrivateKey, *rsa.PublicKey) {
  privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
  if err != nil {
    panic(fmt.Sprintf("failed to generate RSA key: %v", err))
  }
  return privateKey, &privateKey.PublicKey
}

func TestSymmetricKeyCreation(t *testing.T) {
  keyData := randomAES256Key()
  k, err := NewKey("", keyData, KeyUseEncryption)
  if err != nil {
    t.Fatalf("failed to create new key: %v", err)
  }
  if k.Algorithm != "A256GCM" {
    t.Fatalf("unexpected key algorithm: %s", k.Algorithm)
  }
}

func TestAsymmetricKeyCreation(t *testing.T) {
  privKey, pubKey := randomRSAKey()
  privK, err := NewKey("", privKey, KeyUseEncryption)
  if err != nil {
    t.Fatalf("failed to create new key: %v", err)
  }
  if privK.Algorithm != "RSA-OAEP" {
    t.Fatalf("unexpected key algorithm: %s", privK.Algorithm)
  }
  pubK, err := NewKey("", pubKey, KeyUseEncryption)
  if err != nil {
    t.Fatalf("failed to create new key: %v", err)
  }
  if pubK.Algorithm != "RSA-OAEP" {
    t.Fatalf("unexpected key algorithm: %s", pubK.Algorithm)
  }
}

func randomSymmetricKey() *JWK {
  keyData := randomAES256Key()
  k, err := NewKey("", keyData, KeyUseEncryption)
  if err != nil {
    panic(fmt.Sprintf("failed to create new key: %v", err))
  }
  return k
}

func randomAsymmetricKey() (privateKey *JWK, publicKey *JWK) {
  privK, pubK := randomRSAKey()
  var err error
  privateKey, err = NewKey("", privK, KeyUseEncryption)
  if err != nil {
    panic(fmt.Sprintf("failed to create new key: %v", err))
  }
  publicKey, err = NewKey("", pubK, KeyUseEncryption)
  if err != nil {
    panic(fmt.Sprintf("failed to create new key: %v", err))
  }
  return privateKey, publicKey
}

func TestSymmetricKeyWrap(t *testing.T) {
  originalKey := randomSymmetricKey()
  wrapKey := randomSymmetricKey()
  wrappedKey, err := originalKey.Wrap(wrapKey)
  if err != nil {
    t.Fatalf("failed to wrap key: %v", err)
  }
  unwrappedKey, err := wrappedKey.Unwrap(wrapKey)
  if err != nil {
    t.Fatalf("failed to unwrap key: %v", err)
  }
  if !bytes.Equal(unwrappedKey.Key.([]byte), originalKey.Key.([]byte)) {
    t.Fatalf("unwrapped key (%x) does not match original key (%x)", unwrappedKey.Key.([]byte), originalKey.Key.([]byte))
  }
}

func TestAsymmetricKeyWrap(t *testing.T) {
  originalKey := randomSymmetricKey()
  privWrapKey, pubWrapKey := randomAsymmetricKey()
  wrappedKey, err := originalKey.Wrap(pubWrapKey)
  if err != nil {
    t.Fatalf("failed to wrap key: %v", err)
  }
  unwrappedKey, err := wrappedKey.Unwrap(privWrapKey)
  if err != nil {
    t.Fatalf("failed to unwrap key: %v", err)
  }
  if !bytes.Equal(unwrappedKey.Key.([]byte), originalKey.Key.([]byte)) {
    t.Fatalf("unwrapped key (%x) does not match original key (%x)", unwrappedKey.Key.([]byte), originalKey.Key.([]byte))
  }
}

func TestInvalidWrap(t *testing.T) {
  originalKey := randomSymmetricKey()
  invalidWrapKey := &JWK{
    JSONWebKey: jose.JSONWebKey{
      Key:       "invalid key type",
      KeyID:     uuid.New().String(),
      Algorithm: "invalid-alg",
      Use:       string(KeyUseEncryption),
    },
  }
  _, err := originalKey.Wrap(invalidWrapKey)
  if err == nil {
    t.Fatal("expected error when wrapping with invalid key type, but got none")
  }
}

func TestInvalidUnwrap(t *testing.T) {
  originalKey := randomSymmetricKey()
  wrapKey := randomSymmetricKey()
  wrappedKey, err := originalKey.Wrap(wrapKey)
  if err != nil {
    t.Fatalf("failed to wrap key: %v", err)
  }
  invalidUnwrapKey := &JWK{
    JSONWebKey: jose.JSONWebKey{
      Key:       "invalid key type",
      KeyID:     uuid.New().String(),
      Algorithm: "invalid-alg",
      Use:       string(KeyUseEncryption),
    },
  }
  _, err = wrappedKey.Unwrap(invalidUnwrapKey)
  if err == nil {
    t.Fatal("expected error when unwrapping with invalid key type, but got none")
  }
}

func TestUnwrapWithWrongKey(t *testing.T) {
  originalKey := randomSymmetricKey()
  wrapKey := randomSymmetricKey()
  wrappedKey, err := originalKey.Wrap(wrapKey)
  if err != nil {
    t.Fatalf("failed to wrap key: %v", err)
  }
  wrongUnwrapKey := randomSymmetricKey()
  _, err = wrappedKey.Unwrap(wrongUnwrapKey)
  if err == nil {
    t.Fatal("expected error when unwrapping with wrong key, but got none")
  }
}

func TestUnwrapWithWrongKeyType(t *testing.T) {
  originalKey := randomSymmetricKey()
  wrapKey := randomSymmetricKey()
  wrappedKey, err := originalKey.Wrap(wrapKey)
  if err != nil {
    t.Fatalf("failed to wrap key: %v", err)
  }
  wrongUnwrapKey := &JWK{
    JSONWebKey: jose.JSONWebKey{
      Key:       "invalid key type",
      KeyID:     uuid.New().String(),
      Algorithm: "invalid-alg",
      Use:       string(KeyUseEncryption),
    },
  }
  _, err = wrappedKey.Unwrap(wrongUnwrapKey)
  if err == nil {
    t.Fatal("expected error when unwrapping with wrong key type, but got none")
  }
}

func TestKeyEncryptSymmetric(t *testing.T) {
  symKey := randomSymmetricKey()
  ciphertext, err := symKey.Encrypt([]byte("test plaintext"))
  if err != nil {
    t.Fatalf("failed to encrypt data: %v", err)
  }
  plaintext, err := symKey.Decrypt(ciphertext)
  if err != nil {
    t.Fatalf("failed to decrypt data: %v", err)
  }
  if !bytes.Equal(plaintext, []byte("test plaintext")) {
    t.Fatalf("decrypted plaintext (%s) does not match original plaintext (%s)", plaintext, "test plaintext")
  }
}

func TestKeyEncryptAsymmetric(t *testing.T) {
  privKey, pubKey := randomAsymmetricKey()
  ciphertext, err := pubKey.Encrypt([]byte("test plaintext"))
  if err != nil {
    t.Fatalf("failed to encrypt data: %v", err)
  }
  plaintext, err := privKey.Decrypt(ciphertext)
  if err != nil {
    t.Fatalf("failed to decrypt data: %v", err)
  }
  if !bytes.Equal(plaintext, []byte("test plaintext")) {
    t.Fatalf("decrypted plaintext (%s) does not match original plaintext (%s)", plaintext, "test plaintext")
  }
}

package cryptolib

import (
  "crypto/rand"
  "encoding/base64"
  "encoding/json"
  "fmt"
)

type Salt [16]byte

// NewSalt generates a new random salt.
func NewSalt() (*Salt, error) {
  var salt Salt
  _, err := rand.Read(salt[:])
  return &salt, err
}

// String returns the base64 URL-encoded representation of the salt.
func (s *Salt) String() string {
  return base64.RawURLEncoding.EncodeToString(s[:])
}

func (s *Salt) Bytes() []byte {
  return s[:]
}

func (s *Salt) MarshalJSON() ([]byte, error) {
  return json.Marshal(s.Bytes())
}

func (s *Salt) UnmarshalJSON(data []byte) error {
  var saltBytes []byte
  if err := json.Unmarshal(data, &saltBytes); err != nil {
    return fmt.Errorf("failed to unmarshal salt: %w", err)
  }
  copy(s[:], saltBytes)
  return nil
}

func (s *Salt) MarshalText() ([]byte, error) {
  return []byte(s.String()), nil
}

func (s *Salt) UnmarshalText(data []byte) error {
  saltBytes, err := base64.StdEncoding.DecodeString(string(data))
  if err != nil {
    saltBytes, err = base64.RawURLEncoding.DecodeString(string(data))
    if err != nil {
      return fmt.Errorf("failed to decode salt: %w", err)
    }
  }
  copy(s[:], saltBytes)
  return nil
}

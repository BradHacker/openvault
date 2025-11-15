package cryptolib

import (
  "crypto/rand"
  "encoding/hex"
  "encoding/json"
  "errors"
  "fmt"
  "math/big"
  "strings"
)

var (
  ErrBadSecretKeyFormat  = errors.New("invalid secret key format, expected XX-YYYYYY-ZZZZZZ-ZZZZZ-ZZZZZ-ZZZZZ-ZZZZZ")
  ErrBadSecretKeyVersion = errors.New("invalid secret key version")
)

type SecretKeyVersion [2]byte

var (
  SecretKeyVersion1      SecretKeyVersion = [2]byte{'O', '1'}
  LatestSecretKeyVersion                  = SecretKeyVersion1
)

type SecretKey struct {
  Version   SecretKeyVersion
  AccountID [6]byte
  Secret    [26]byte
}

// Alphabet: {2-9, A-H, J-N, P-T, V-Z} (avoids easily confusable characters)
const SecretKeyAlphabet = "23456789ABCDEFGHJKLMNPQRSTVWXYZ"

func randomWithAlphabet(dst []byte) error {
  var maxAlphabetIndex = big.NewInt(int64(len(SecretKeyAlphabet)))
  for i := range dst {
    n, err := rand.Int(rand.Reader, maxAlphabetIndex)
    if err != nil {
      return fmt.Errorf("failed to read random bytes to generate secret key: %w", err)
    }
    dst[i] = SecretKeyAlphabet[n.Int64()]
  }
  return nil
}

func NewSecretKey() (*SecretKey, error) {
  acctId := [6]byte{}
  if err := randomWithAlphabet(acctId[:]); err != nil {
    return nil, fmt.Errorf("failed to read random bytes to generate secret key: %w", err)
  }
  secret := [26]byte{}
  if err := randomWithAlphabet(secret[:]); err != nil {
    return nil, fmt.Errorf("failed to read random bytes to generate secret key: %w", err)
  }
  return &SecretKey{
    Version:   LatestSecretKeyVersion,
    AccountID: acctId,
    Secret:    secret,
  }, nil
}

func (sk *SecretKey) MarshalText() ([]byte, error) {
  secretKey := strings.ToUpper(string(sk.Version[:])) + "-"
  secretKey += strings.ToUpper(string(sk.AccountID[:])) + "-"
  secretKey += strings.ToUpper(string(sk.Secret[0:6])) + "-"
  secretKey += strings.ToUpper(string(sk.Secret[6:11])) + "-"
  secretKey += strings.ToUpper(string(sk.Secret[11:16])) + "-"
  secretKey += strings.ToUpper(string(sk.Secret[16:21])) + "-"
  secretKey += strings.ToUpper(string(sk.Secret[21:26]))
  return []byte(secretKey), nil
}

func (sk *SecretKey) UnmarshalText(data []byte) error {
  keyStr := string(data)

  // Check if the key is obfuscated
  if strings.HasSuffix(keyStr, "obfus") {
    deobfuscated, err := deobfuscateSecretKey(keyStr)
    if err != nil {
      return fmt.Errorf("failed to deobfuscate secret key: %w", err)
    }
    keyStr = deobfuscated
  }

  parts := strings.SplitN(keyStr, "-", 7)
  if len(parts) != 7 {
    return ErrBadSecretKeyFormat
  }
  // Uppercase all string
  for i := range parts {
    parts[i] = strings.ToUpper(parts[i])
  }
  // Check version
  switch parts[0] {
  case string(SecretKeyVersion1[:]):
    sk.Version = SecretKeyVersion1
  default:
    return fmt.Errorf("%w: %v", ErrBadSecretKeyFormat, fmt.Errorf("invalid version %q", parts[0]))
  }
  // Check Account ID
  if len(parts[1]) != 6 {
    return fmt.Errorf("%w: %v", ErrBadSecretKeyFormat, fmt.Errorf("invalid account id %q", parts[1]))
  }
  copy(sk.AccountID[:], parts[1])
  // Check Secret Key
  if len(parts[2]) != 6 {
    return fmt.Errorf("%w: %v", ErrBadSecretKeyFormat, fmt.Errorf("invalid secret part %q", parts[2]))
  }
  copy(sk.Secret[0:6], parts[2])
  for i := 0; i < 4; i++ {
    if len(parts[3+i]) != 5 {
      return fmt.Errorf("%w: %v", ErrBadSecretKeyFormat, fmt.Errorf("invalid secret part %q", parts[3+i]))
    }
    copy(sk.Secret[6+(i*5):11+(i*5)], parts[3+i])
  }
  return nil
}

func (sk *SecretKey) String() string {
  skBytes, _ := sk.MarshalText()
  return string(skBytes)
}

var (
  SECRET_KEY_OBFUSCATION_KEY []byte = []byte("This is an obfuscation key used to mask the secret key in the local database and nothing more. If this seems interesting to you, come work with us :)")
)

func (sk *SecretKey) Obfuscate() string {
  obfuscated := make([]byte, len(sk.String()))
  for i, b := range []byte(sk.String()) {
    obfuscated[i] = b ^ SECRET_KEY_OBFUSCATION_KEY[i%len(SECRET_KEY_OBFUSCATION_KEY)]
  }
  obfuscatedHex := hex.EncodeToString(obfuscated)
  obfuscatedHex += "obfus"
  return obfuscatedHex
}

func deobfuscateSecretKey(obfuscatedSecretKey string) (string, error) {
  obfuscatedHex := strings.TrimSuffix(obfuscatedSecretKey, "obfus")
  obfuscatedBytes, err := hex.DecodeString(obfuscatedHex)
  if err != nil {
    return "", err
  }
  secretKey := make([]byte, len(obfuscatedBytes))
  for i, b := range []byte(obfuscatedBytes) {
    secretKey[i] = b ^ SECRET_KEY_OBFUSCATION_KEY[i%len(SECRET_KEY_OBFUSCATION_KEY)]
  }
  return string(secretKey), nil
}

func (sk *SecretKey) MarshalJSON() ([]byte, error) {
  return json.Marshal(sk.Obfuscate())
}

func (sk *SecretKey) UnmarshalJSON(data []byte) error {
  var keyStr string
  if err := json.Unmarshal(data, &keyStr); err != nil {
    return fmt.Errorf("failed to unmarshal secret key: %w", err)
  }
  // Check if the key is obfuscated
  if strings.HasSuffix(keyStr, "obfus") {
    deobfuscated, err := deobfuscateSecretKey(keyStr)
    if err != nil {
      return fmt.Errorf("failed to deobfuscate secret key: %w", err)
    }
    keyStr = deobfuscated
  }
  return sk.UnmarshalText([]byte(keyStr))
}

package cryptolib

import (
  "bytes"
  "crypto/rand"
  "testing"
)

// helper to make deterministic secret key for reproducible test vector
func fixedSecretKey() *SecretKey {
  sk := &SecretKey{}
  sk.Version = LatestSecretKeyVersion
  copy(sk.AccountID[:], []byte("ABCDEF"))
  copy(sk.Secret[:], []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZA")) // 27, truncate to 26
  return sk
}

func randomSalt16() Salt {
  var s Salt
  rand.Read(s[:])
  return s
}

// TestDeriveDeterministic tests that the derive function produces
// deterministic outputs for equivalent inputs.
func TestDeriveDeterministic(t *testing.T) {
  secret := fixedSecretKey()
  salt := randomSalt16()
  rounds := 1000
  auk1, err := DeriveAUK(&AUKParams{
    Email:    "User@example.com",
    Password: "  passWORD",
    Salt:     &salt,
    Secret:   secret,
    Rounds:   rounds,
  })
  if err != nil {
    t.Fatalf("derive: %v", err)
  }
  if len(auk1.Key.([]byte)) != 32 {
    t.Fatalf("expected 32 byte auk1 key")
  }
  // Same inputs produce same result (id fixed to AccountUnlockKeyID)
  auk2, err := DeriveAUK(&AUKParams{
    Email:    "user@example.com",
    Password: "passWORD",
    Salt:     &salt,
    Secret:   secret,
    Rounds:   rounds,
  })
  if err != nil {
    t.Fatalf("derive2: %v", err)
  }
  if len(auk2.Key.([]byte)) != 32 {
    t.Fatalf("expected 32 byte auk2 key")
  }
  if !bytes.Equal(auk1.Key.([]byte), auk2.Key.([]byte)) {
    t.Fatalf("AUK mismatch for equivalent normalized inputs")
  }
}

// TestDeriveNormalization tests that the derive function properly
// normalizes email and password inputs.
func TestDeriveNormalization(t *testing.T) {
  secret := fixedSecretKey()
  salt := randomSalt16()
  rounds := 1000
  auk1, err := DeriveAUK(&AUKParams{
    Email:    "User@example.com",
    Password: " passWORD",
    Salt:     &salt,
    Secret:   secret,
    Rounds:   rounds,
  })
  if err != nil {
    t.Fatalf("derive: %v", err)
  }
  auk2, err := DeriveAUK(&AUKParams{
    Email:    "user@example.com",
    Password: "  passWORD",
    Salt:     &salt,
    Secret:   secret,
    Rounds:   rounds,
  })
  if err != nil {
    t.Fatalf("derive2: %v", err)
  }
  if !bytes.Equal(auk1.Key.([]byte), auk2.Key.([]byte)) {
    t.Fatalf("AUK should match for normalized inputs")
  }
}

// TestDeriveDifferentInputs tests that the derive function produces
// different outputs for different inputs.
func TestDeriveDifferentInputs(t *testing.T) {
  secret := fixedSecretKey()
  salt := randomSalt16()
  rounds := 1000
  aukBase, err := DeriveAUK(&AUKParams{
    Email:    "user@example.com",
    Password: "passWORD",
    Salt:     &salt,
    Secret:   secret,
    Rounds:   rounds,
  })
  if err != nil {
    t.Fatalf("derive: %v", err)
  }

  tests := []struct {
    name   string
    modify func(params *AUKParams)
  }{
    {
      name: "different email",
      modify: func(params *AUKParams) {
        params.Email = "different@example.com"
      },
    },
    {
      name: "different password",
      modify: func(params *AUKParams) {
        params.Password = "differentPASSWORD"
      },
    },
    {
      name: "different salt",
      modify: func(params *AUKParams) {
        newSalt := randomSalt16()
        params.Salt = &newSalt
      },
    },
    {
      name: "different rounds",
      modify: func(params *AUKParams) {
        params.Rounds = 2000
      },
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      params := &AUKParams{
        Email:    "user@example.com",
        Password: "passWORD",
        Salt:     &salt,
        Secret:   secret,
        Rounds:   rounds,
      }
      tt.modify(params)
      auk, err := DeriveAUK(params)
      if err != nil {
        t.Fatalf("derive: %v", err)
      }
      if bytes.Equal(auk.Key.([]byte), aukBase.Key.([]byte)) {
        t.Fatalf("AUK should differ for %q", tt.name)
      }
    })
  }
}

// TestDeriveValidation tests that the derive function properly
// validates its inputs.
func TestDeriveValidation(t *testing.T) {
  secret := fixedSecretKey()
  salt := randomSalt16()
  rounds := 1000

  tests := []struct {
    name      string
    modify    func(params *AUKParams)
    expectErr bool
  }{
    {
      name: "zero rounds",
      modify: func(params *AUKParams) {
        params.Rounds = 0
      },
      expectErr: true,
    },
    {
      name: "invalid email",
      modify: func(params *AUKParams) {
        params.Email = "invalid-email"
      },
      expectErr: true,
    },
    {
      name: "empty password",
      modify: func(params *AUKParams) {
        params.Password = ""
      },
      expectErr: true,
    },
    {
      name: "nil salt",
      modify: func(params *AUKParams) {
        params.Salt = nil
      },
      expectErr: true,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      params := &AUKParams{
        Email:    "user@example.com",
        Password: "passWORD",
        Salt:     &salt,
        Secret:   secret,
        Rounds:   rounds,
      }
      tt.modify(params)
      _, err := DeriveAUK(params)
      if tt.expectErr {
        if err == nil {
          t.Fatalf("expected error but got none")
        }
      } else {
        if err != nil {
          t.Fatalf("unexpected error: %v", err)
        }
      }
    })
  }
}

package vault_test

import (
	"fmt"
	"secret/vault"
	"testing"
)

var v vault.Vault
var vf *vault.Vault

func init() {
	v = vault.MemoryVault("my-fake-key")
	vf = vault.FileVault("my-fake-key", ".secrets")
}
func TestSetKeyInVault(t *testing.T) {
	err := v.SetKeyInMemoryVault("demo-key", "some crazy value")
	if err != nil {
		t.Error(err)
	}
}

func TestGetKeyInVault(t *testing.T) {
	// set some stub value
	err := v.SetKeyInMemoryVault("demo-key", "some crazy value")
	if err != nil {
		t.Error(err)
	}

	plain, err := v.GetKeyInMemoryVault("demo-key")
	result := "some crazy value"
	if err != nil {
		t.Error(err)
	}
	if result != plain {
		t.Errorf("Result produced by GetKeyInVault is not correct")
	}
}

func TestSetKeyInFileVault(t *testing.T) {
	err := vf.SetKeyInFileVault("demo-key", "some crazy value")
	if err != nil {
		t.Error(err)
	}
}

func TestGetKeyInFileVault(t *testing.T) {
	// set some stub value
	err := vf.SetKeyInFileVault("demo-key", "some crazy value")
	if err != nil {
		t.Error(err)
	}

	plain, err := vf.GetKeyInFileVault("demo-key")
	fmt.Printf("❎ plain : %s\n", plain)
	result := "some crazy value"
	if err != nil {
		t.Error(err)
	}
	if result != plain {
		t.Errorf("Result produced by GetKeyInVault is not correct")
	}
}

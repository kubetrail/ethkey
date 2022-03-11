package run

import (
	"fmt"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestDerivedKeys(t *testing.T) {
	if err := testDerivedKeys(t); err != nil {
		t.Fatal(err)
	}
}

func testDerivedKeys(t *testing.T) error {
	mnemonic := "bag slim bind media where castle spend invite change rookie tumble honey sample kitten anxiety draft giraffe chief hurt music olympic initial fit wink"
	passphrase := "test"
	derivationPath := "m/44'/60'/0/0"

	seed := bip39.NewSeed(mnemonic, passphrase)

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return fmt.Errorf("failed to generate root key: %w", err)
	}

	path, err := hdwallet.ParseDerivationPath(derivationPath)
	if err != nil {
		return fmt.Errorf("failed to parse derivation path: %w", err)
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		return fmt.Errorf("failed to derive account from path: %w", err)
	}
}

package run

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func Validate(cmd *cobra.Command, args []string) error {
	prompt, err := getPromptStatus()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if prompt {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter prv or pub key: "); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	inputReader := bufio.NewReader(cmd.InOrStdin())
	key, err := inputReader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read mnemonic from input: %w", err)
	}
	key = strings.Trim(key, "\n")

	if len(key) < 2 {
		return fmt.Errorf("invalid key length")
	}

	if key[:2] == "0x" {
		address := common.Address{}
		if err := address.UnmarshalText([]byte(key)); err != nil {
			if err != nil {
				return fmt.Errorf("public key is invalid: %w", err)
			}
		}
		if prompt {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), "public key is valid"); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}

		return nil
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid private key, failed to derive public key of type *ecdsa.PublicKey, instead got: %T", publicKey)
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	if prompt {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(),
			"private key is valid with public address: %s\n", address.Hex()); err != nil {
			return fmt.Errorf("failed to write to output")
		}
	}

	return nil
}

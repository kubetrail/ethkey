package run

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kubetrail/ethkey/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Verify(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag(flags.Hash, cmd.Flags().Lookup(flags.Hash))
	_ = viper.BindPFlag(flags.Sign, cmd.Flags().Lookup(flags.Sign))
	_ = viper.BindPFlag(flags.PubKey, cmd.Flags().Lookup(flags.PubKey))

	hash := viper.GetString(flags.Hash)
	sign := viper.GetString(flags.Sign)
	key := viper.GetString(flags.PubKey)

	printOk := false
	if len(hash) == 0 ||
		len(sign) == 0 ||
		len(key) == 0 {
		printOk = true
	}

	inputReader := bufio.NewReader(cmd.InOrStdin())
	prompt, err := getPromptStatus()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if len(key) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter pub key: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		key, err = inputReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read pub key from input: %w", err)
		}
		key = strings.Trim(key, "\n")
	}

	if len(hash) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter hash: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		hash, err = inputReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read hash from input: %w", err)
		}
		hash = strings.Trim(hash, "\n")
	}

	if len(sign) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter sign: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		sign, err = inputReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read sign from input: %w", err)
		}
		sign = strings.Trim(sign, "\n")
	}

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	signBytes, err := hex.DecodeString(sign)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	pubKey, err := crypto.SigToPub(hashBytes, signBytes)
	if err != nil {
		return fmt.Errorf("failed to derive public key from hash and sign: %w", err)
	}

	address := crypto.PubkeyToAddress(*pubKey)
	if address.Hex() != key {
		return fmt.Errorf("public key of signature for given hash does not match input public key")
	}

	if printOk {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "signature is valid for given hash and public key"); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	return nil
}

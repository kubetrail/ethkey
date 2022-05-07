package run

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/ethkey/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Verify(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Hash, cmd.Flags().Lookup(flags.Hash))
	_ = viper.BindPFlag(flags.Sign, cmd.Flags().Lookup(flags.Sign))
	_ = viper.BindPFlag(flags.Key, cmd.Flags().Lookup(flags.Key))

	hash := viper.GetString(flags.Hash)
	sign := viper.GetString(flags.Sign)
	key := viper.GetString(flags.Key)

	var verified bool

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if len(key) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter pub key: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		key, err = keys.Read(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read pub key from input: %w", err)
		}
	}

	if len(hash) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter hash: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		hash, err = keys.Read(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read hash from input: %w", err)
		}
	}

	if len(sign) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter sign: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		sign, err = keys.Read(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read signature from input: %w", err)
		}
	}

	if !keys.IsValidBase58String(hash) {
		return fmt.Errorf("hash is not a valid base58 string")
	}

	if !keys.IsValidBase58String(sign) {
		return fmt.Errorf("signature is not a valid base58 string")
	}

	pubKey, err := crypto.SigToPub(base58.Decode(hash), base58.Decode(sign))
	if err != nil {
		return fmt.Errorf("failed to derive public key from hash and sign: %w", err)
	}

	address := crypto.PubkeyToAddress(*pubKey)
	if address.Hex() == key {
		verified = true
	}

	type output struct {
		Verified bool `json:"verified" yaml:"verified"`
	}

	out := &output{Verified: verified}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), verified); err != nil {
			return fmt.Errorf("failed to write signature verification to output: %w", err)
		}
	case flags.OutputFormatYaml:
		jb, err := yaml.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write signature verification to output: %w", err)
		}
	case flags.OutputFormatJson:
		jb, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write signature verification to output: %w", err)
		}
	default:
		return fmt.Errorf("failed to format in requested format, %s is not supported", persistentFlags.OutputFormat)
	}

	return nil
}

package run

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/ethkey/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Validate(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Key, cmd.Flags().Lookup(flags.Key))
	key := viper.GetString(flags.Key)

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	var valid bool

	if len(key) == 0 {
		if len(args) == 0 {
			if prompt {
				if err := keys.Prompt(cmd.OutOrStdout()); err != nil {
					return fmt.Errorf("failed to prompt for key: %w", err)
				}
			}

			key, err = keys.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read key from input: %w", err)
			}
		} else {
			key = args[0]
		}
	}

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
		valid = true
	} else {
		privateKey, err := crypto.HexToECDSA(key)
		if err != nil {
			return fmt.Errorf("invalid private key: %w", err)
		}

		publicKey := privateKey.Public()
		if _, ok := publicKey.(*ecdsa.PublicKey); ok {
			valid = true
		}
	}

	type output struct {
		Valid bool `json:"valid" yaml:"valid"`
	}

	out := &output{Valid: valid}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), valid); err != nil {
			return fmt.Errorf("failed to write key validity to output: %w", err)
		}
	case flags.OutputFormatYaml:
		jb, err := yaml.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write key validity to output: %w", err)
		}
	case flags.OutputFormatJson:
		jb, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write key validity to output: %w", err)
		}
	default:
		return fmt.Errorf("failed to format in requested format, %s is not supported", persistentFlags.OutputFormat)
	}

	return nil
}

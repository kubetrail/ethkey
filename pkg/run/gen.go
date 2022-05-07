package run

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubetrail/bip39/pkg/mnemonics"
	"github.com/kubetrail/bip39/pkg/passphrases"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/bip39/pkg/seeds"
	"github.com/kubetrail/ethkey/pkg/flags"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Gen(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.UsePassphrase, cmd.Flags().Lookup(flags.UsePassphrase))
	_ = viper.BindPFlag(flags.SkipMnemonicValidation, cmd.Flags().Lookup(flags.SkipMnemonicValidation))
	_ = viper.BindPFlag(flags.DerivationPath, cmd.Flags().Lookup(flags.DerivationPath))
	_ = viper.BindPFlag(flags.InputHexSeed, cmd.Flags().Lookup(flags.InputHexSeed))
	_ = viper.BindPFlag(flags.MnemonicLanguage, cmd.Flag(flags.MnemonicLanguage))

	usePassphrase := viper.GetBool(flags.UsePassphrase)
	skipMnemonicValidation := viper.GetBool(flags.SkipMnemonicValidation)
	derivationPath := viper.GetString(flags.DerivationPath)
	inputHexSeed := viper.GetBool(flags.InputHexSeed)
	language := viper.GetString(flags.MnemonicLanguage)

	derivationPath = strings.ToLower(derivationPath)
	derivationPath = strings.ReplaceAll(derivationPath, "h", "'")

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	var passphrase string
	var seed []byte

	if inputHexSeed && usePassphrase {
		return fmt.Errorf("cannot use passphrase when entering seed")
	}

	if inputHexSeed && skipMnemonicValidation {
		return fmt.Errorf("dont use --skip-mnemonic-validation when entering seed")
	}

	if !inputHexSeed {
		var mnemonic string
		if len(args) == 0 {
			if prompt {
				if err := mnemonics.Prompt(cmd.OutOrStdout()); err != nil {
					return fmt.Errorf("failed to write to output: %w", err)
				}
			}

			mnemonic, err = mnemonics.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read mnemonic from input: %w", err)
			}
		} else {
			mnemonic = mnemonics.NewFromFields(args)
		}

		if !skipMnemonicValidation {
			if mnemonic, err = mnemonics.Translate(mnemonic, language, mnemonics.LanguageEnglish); err != nil {
				return fmt.Errorf("failed to translate mnemonic to English, alternatively try --skip-mnemonic-validation flag: %w", err)
			}
		} else {
			mnemonic = mnemonics.Tidy(mnemonic)
		}

		if usePassphrase {
			passphrase, err = passphrases.New(cmd.OutOrStdout())
			if err != nil {
				return fmt.Errorf("failed to get passphrase: %w", err)
			}
		}

		seed = seeds.New(mnemonic, passphrase)
	} else {
		if len(args) == 0 {
			if prompt {
				if err := seeds.Prompt(cmd.OutOrStdout()); err != nil {
					return fmt.Errorf("failed to prompt for seed: %w", err)
				}
			}

			seed, err = seeds.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("invalid seed: %w", err)
			}
		} else {
			seed, err = hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode seed: %w", err)
			}
		}
	}

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

	prvHex, err := wallet.PrivateKeyHex(account)
	if err != nil {
		return fmt.Errorf("failed to generate private key for account: %w", err)
	}

	outPrv := fmt.Sprintf("%s", prvHex)
	outPub := fmt.Sprintf("%s", account.Address.Hex())

	type output struct {
		Prv string `json:"prv,omitempty" yaml:"prv,omitempty"`
		Pub string `json:"pub,omitempty" yaml:"pub,omitempty"`
	}

	out := &output{
		Prv: outPrv,
		Pub: outPub,
	}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative, flags.OutputFormatYaml:
		jb, err := yaml.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write key to output: %w", err)
		}
	case flags.OutputFormatJson:
		jb, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to serialize output to json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write key to output: %w", err)
		}
	default:
		return fmt.Errorf("failed to format in requested format, %s is not supported", persistentFlags.OutputFormat)
	}

	return nil
}

package run

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"syscall"

	"github.com/kubetrail/ethkey/pkg/flags"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/term"
)

func Gen(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag(flags.UsePassphrase, cmd.Flags().Lookup(flags.UsePassphrase))
	_ = viper.BindPFlag(flags.SkipMnemonicValidation, cmd.Flags().Lookup(flags.SkipMnemonicValidation))
	_ = viper.BindPFlag(flags.DerivationPath, cmd.Flags().Lookup(flags.DerivationPath))
	_ = viper.BindPFlag(flags.InputHexSeed, cmd.Flags().Lookup(flags.InputHexSeed))

	usePassphrase := viper.GetBool(flags.UsePassphrase)
	skipMnemonicValidation := viper.GetBool(flags.SkipMnemonicValidation)
	derivationPath := viper.GetString(flags.DerivationPath)
	inputHexSeed := viper.GetBool(flags.InputHexSeed)

	prompt, err := getPromptStatus()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	var passphrase []byte
	var seed []byte

	if inputHexSeed && usePassphrase {
		return fmt.Errorf("cannot use passphrase when entering seed")
	}

	if inputHexSeed && skipMnemonicValidation {
		return fmt.Errorf("dont use --skip-mnemonic-validation when entering seed")
	}

	if !inputHexSeed {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter mnemonic: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}

		inputReader := bufio.NewReader(cmd.InOrStdin())
		mnemonic, err := inputReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read mnemonic from input: %w", err)
		}
		mnemonic = strings.Trim(mnemonic, "\n")

		if !skipMnemonicValidation && !bip39.IsMnemonicValid(mnemonic) {
			return fmt.Errorf("mnemonic is invalid or please use --skip-mnemonic-validation flag")
		}

		if usePassphrase {
			if prompt {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter secret passphrase: "); err != nil {
					return fmt.Errorf("failed to write to output: %w", err)
				}
			}

			passphrase, err = term.ReadPassword(syscall.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read secret passphrase from input: %w", err)
			}
			if _, err := fmt.Fprintln(cmd.OutOrStdout()); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter secret passphrase again: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}

			passphraseConfirm, err := term.ReadPassword(syscall.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read secret passphrase from input: %w", err)
			}
			if _, err := fmt.Fprintln(cmd.OutOrStdout()); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}

			if !bytes.Equal(passphrase, passphraseConfirm) {
				return fmt.Errorf("passphrases do not match")
			}
		}
		seed = bip39.NewSeed(mnemonic, string(passphrase))
	} else {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter seed in hex: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}

		inputReader := bufio.NewReader(cmd.InOrStdin())
		hexSeed, err := inputReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read mnemonic from input: %w", err)
		}
		hexSeed = strings.Trim(hexSeed, "\n")

		seed, err = hex.DecodeString(hexSeed)
		if err != nil {
			return fmt.Errorf("invalid seed: %w", err)
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

	if prompt {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "pub:", outPub); err != nil {
			return fmt.Errorf("failed to write key to output: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "prv:", outPrv); err != nil {
			return fmt.Errorf("failed to write key to output: %w", err)
		}

		return nil
	}

	jb, err := json.Marshal(
		struct {
			Prv string `json:"prv,omitempty"`
			Pub string `json:"pub,omitempty"`
		}{
			Prv: outPrv,
			Pub: outPub,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to serialize output: %w", err)
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
		return fmt.Errorf("failed to write key to output: %w", err)
	}

	return nil
}

package run

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kubetrail/ethkey/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Sign(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag(flags.Filename, cmd.Flags().Lookup(flags.Filename))
	fileName := viper.GetString(flags.Filename)

	prompt, err := getPromptStatus()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if prompt {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter prv key: "); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	inputReader := bufio.NewReader(cmd.InOrStdin())
	key, err := inputReader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read mnemonic from input: %w", err)
	}
	key = strings.Trim(key, "\n")

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	var b []byte
	if len(fileName) > 0 {
		if fileName == "-" {
			if b, err = io.ReadAll(cmd.InOrStdin()); err != nil {
				return fmt.Errorf("failed to read stdin input: %w", err)
			}
		} else {
			if b, err = os.ReadFile(fileName); err != nil {
				return fmt.Errorf("failed to read input file %s: %w", fileName, err)
			}
		}
	} else {
		if len(args) == 0 {
			return fmt.Errorf("no input file or args, pl. provide input to sign")
		}
		b = []byte(strings.Join(args, " "))
	}

	hash := crypto.Keccak256(b)

	sign, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign input data hash: %w", err)
	}

	hashHex := hex.EncodeToString(hash)
	signHex := hex.EncodeToString(sign)

	if prompt {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "hash: ", hashHex); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "sign: ", signHex); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}

		return nil
	}

	jb, err := json.Marshal(
		struct {
			Hash string `json:"hash,omitempty"`
			Sign string `json:"sign,omitempty"`
		}{
			Hash: hashHex,
			Sign: signHex,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to serialize output: %w", err)
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
		return fmt.Errorf("failed to write to output: %w", err)
	}

	return nil
}

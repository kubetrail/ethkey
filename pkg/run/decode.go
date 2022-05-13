package run

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/ethkey/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Decode(cmd *cobra.Command, args []string) error {
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Key, cmd.Flag(flags.Key))

	key := viper.GetString(flags.Key)

	var publicKey *ecdsa.PublicKey
	var privateKey *ecdsa.PrivateKey

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	if len(key) == 0 {
		if len(args) == 0 {
			if prompt {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter pub key: "); err != nil {
					return fmt.Errorf("failed to write to output: %w", err)
				}
			}
			key, err = keys.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read pub key from input: %w", err)
			}
		} else {
			key = args[0]
		}
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return fmt.Errorf("failed to decode input key as hex string: %w", err)
	}

	switch len(keyBytes) {
	case 32:
		privateKey, err = crypto.HexToECDSA(key)
		if err != nil {
			return fmt.Errorf("failed to decode private key: %w", err)
		}
		if privateKey.Curve == nil ||
			privateKey.D == nil ||
			privateKey.PublicKey.Curve == nil ||
			privateKey.PublicKey.X == nil ||
			privateKey.PublicKey.Y == nil {
			return fmt.Errorf("failed to decode as ecdsa private key, invalid bytes")
		}
	case 64:
		keyBytes, _ = hex.DecodeString(fmt.Sprintf("%s%s", keyPrefix, key))
		x, y := elliptic.Unmarshal(crypto.S256(), keyBytes)
		publicKey = &ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     x,
			Y:     y,
		}
		if publicKey.Curve == nil ||
			publicKey.X == nil ||
			publicKey.Y == nil {
			return fmt.Errorf("failed to decode as ecdsa public key, invalid bytes")
		}
	default:
		return fmt.Errorf("invalid key len, needs to be either 32 for private key or 65 for public key, received %d", len(keyBytes))
	}

	if publicKey == nil {
		var ok bool
		publicKey, ok = privateKey.Public().(*ecdsa.PublicKey)
		if !ok {
			return fmt.Errorf("failed to type assert public key from private key as ecdsa public key")
		}
	}

	out := &output{KeyType: keyType}

	if privateKey != nil {
		out.PrvHex = hex.EncodeToString(crypto.FromECDSA(privateKey))
	}

	if publicKey != nil {
		b := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
		pubHex := hex.EncodeToString(b)
		out.PubHex = pubHex[2:]
		out.Addr = crypto.PubkeyToAddress(*publicKey).String()
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

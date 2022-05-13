/*
Copyright © 2022 kubetrail.io authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/kubetrail/bip39/pkg/mnemonics"
	"github.com/kubetrail/ethkey/pkg/flags"
	"github.com/kubetrail/ethkey/pkg/run"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate key from mnemonic",
	Long: `This command generates private/public keys
from mnemonic and optional secret passphrase per BIP-32 spec.

Alternatively, a seed in hex format can be provided bypassing
all mnemonic related computation and be directly used for
key generation

The keys are generated based on a chain derivation path
Path     |   Remark
---------|--------------------------------------------------------------
m        |   Master key (aka root key)
m/0      |   First child of master key
m/0'     |   First hardened child of master key
m/0/0    |   First child of first child of master key
m/0'/0   |   First child of first hardened child of master key
m/0/0'   |   First hardened child of first child of master key
m/0'/0'  |   First hardened child of first hardened child of master key

Per BIP-44 spec, the derivation path for ETH takes the form:
m/44'/60'/0'/0/0, where 44 refers to the purpose being satisfying
BIP-44 spec, 60 being the ETH blockchain identifier and rest of the
items being specific to account.

Mnemonic language can be specified from the following list:
1. English (default)
2. Japanese
3. ChineseSimplified
4. ChineseTraditional
5. Czech
6. French
7. Italian
8. Korean
9. Spanish

BIP-39 proposal: https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki

Please note that same keys will be generated for mnemonics from different languages
if the underlying entropy is the same. In other words, keys are always
generated after translating input mnemonic to English.

To generate keys, first optionally create a new mnemonic:
$ bip39 gen --length=12
cushion cover cupboard brother quiz board busy loyal kidney slogan catch pencil

Then use the mnemonic to create the ETH keys:
$ ethkey gen
Enter mnemonic: cushion cover cupboard brother quiz board busy loyal kidney slogan catch pencil
seed: c80211bbdee485ca4e7c9039847be4ee004d93af33a1589a250ac8bdaaa597305fbc2b8bdc89d6d6dbf9fd5c285c9bd540d67aa89882b1633b5f6abaf4abb898
prvHex: ed234d0929176fc58f699be15c7f606f745223d93ceb3b4042e55e825484c043
pubHex: d1b975977b838babda7675859fa1958c94a0a9c0615e0cdcb10ac270ce61f6147c3e204443048026d63b84fdcf4398cebf077773815f50727fa46459c011d6a5
addr: 0xAE600D1F94680Ef43Ab12F8d618F8aAfC208FE25
keyType: ecdsa

An ecdsa private/public key pair is generated and corresponding address is derived.
Please note that public key hex string is missing prefix 04 in case you are trying to
unmarshal it separately to raw Public key data structure.
`,
	RunE: run.Gen,
}

func init() {
	rootCmd.AddCommand(genCmd)
	f := genCmd.Flags()

	f.String(flags.DerivationPath, "m/44'/60'/0'/0/0", "Chain Derivation path")
	f.Bool(flags.UsePassphrase, false, "Prompt for secret passphrase")
	f.Bool(flags.InputHexSeed, false, "Treat input as hex seed instead of mnemonic")
	f.Bool(flags.SkipMnemonicValidation, false, "Skip mnemonic validation")
	f.String(flags.MnemonicLanguage, mnemonics.LanguageEnglish, "Mnemonic language")
}

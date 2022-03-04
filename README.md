# ethkey
Generate Ethereum keys using mnemonics

## disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

Download this repo to a folder and cd to it. Make sure `go` toolchain
is installed
```bash
go install
```
Add autocompletion for `bash` to your `.bashrc`
```bash
source <(ethkey completion bash)
```

## generate keys
Ethereum keys can be generated using mnemonic. [bip39](https://github.com/kubetrail/bip39)
can be used for generating new mnemonics:
```bash
bip39 gen
patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
```
```bash
ethkey gen
Enter mnemonic: patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
pub: 0x33009656efD3e4eA8763B72B020C2327Dee0B2db
prv: 12bea0236685f934e39e27bef1793966800f882912c099b4c584a8fbfd28b6e1
```

The default chain derivation path is `m/44'/60'/0'/0/0`, were `60'` stands for
hardened Ethereum chain address and `44'` stands for hardened purpose, implying
`BIP-44` spec.

Keys can be additionally protected using a passphrase:
```bash
ethkey gen --use-passphrase 
Enter mnemonic: patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
Enter secret passphrase: 
Enter secret passphrase again: 
pub: 0x6545273938F901dE49bc58987c5f3fBD01093DF0
prv: 68458a246782ad5cdff0dcd3b2e29e5006a0aa010b06fa4d97f30a7022e2414c
```

The chain derivation path can be changed
```bash
ethkey gen --derivation-path="m/44'/501'/0'/0/0"
Enter mnemonic: patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
pub: 0x810d963eC3C2DfD299E089C42cc7d07089FA5eD2
prv: dd0808b874693a136384b0f4cf923b6caf59a5068bb8bfeb1e38420b86f22f6f
```

Mnemonic is validated and expected to comply to `BIP-39` standard, however, an 
arbitrary mnemonic can be used by switching off validation

```bash
ethkey gen --skip-mnemonic-validation 
Enter mnemonic: this is an invalid mnemonic
pub: 0xc1541003D25C206873F5c28Ea684E1072026FC9A
prv: fc30faae343aba4a5151a56881e9c4ff332563fa8eedda661a7681db4ea604bb
```

## validate keys
Key validation checks for key string format, length and other characteristics.
For instance, if a private key is entered, it also checks if a public key
can be derived from it.

```bash
ethkey gen
Enter mnemonic: few happy dragon spray much obvious total drive hat brain impose bright test there outside peasant share kitchen prefer inmate moment cactus forward crisp
pub: 0x0dD73dE2AC23E6f4928D582aa6510144790DA88e
prv: 70f18df4c72c80a4c0bd47c6ec1b61dd4251c9a77104150c891f8d27a96beb73
```

These keys can be validated:
```bash
ethkey validate 
Enter prv or pub key: 0x0dD73dE2AC23E6f4928D582aa6510144790DA88e
public key is valid
```

```bash
ethkey validate 
Enter prv or pub key: 70f18df4c72c80a4c0bd47c6ec1b61dd4251c9a77104150c891f8d27a96beb73
private key is valid with public address: 0x0dD73dE2AC23E6f4928D582aa6510144790DA88e
```

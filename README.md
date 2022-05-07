# ethkey
Generate Ethereum keys using mnemonics or hex seed

## disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

```bash
go install github.com/kubetrail/ethkey@latest
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
```
```text
patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
```

```bash
ethkey gen
```
```text
Enter mnemonic: patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
pub: 0x33009656efD3e4eA8763B72B020C2327Dee0B2db
prv: 12bea0236685f934e39e27bef1793966800f882912c099b4c584a8fbfd28b6e1
```

Alternatively, pass mnemonic as CLI args
```bash
ethkey gen patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
```

The default chain derivation path is `m/44'/60'/0'/0/0`, were `60'` stands for
hardened Ethereum chain address and `44'` stands for hardened purpose, implying
`BIP-44` spec.

Keys can be additionally protected using a passphrase:
```bash
ethkey gen --use-passphrase
```
```text
Enter mnemonic: patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
Enter secret passphrase: 
Enter secret passphrase again: 
pub: 0x6545273938F901dE49bc58987c5f3fBD01093DF0
prv: 68458a246782ad5cdff0dcd3b2e29e5006a0aa010b06fa4d97f30a7022e2414c
```

The chain derivation path can be changed
```bash
ethkey gen --derivation-path="m/44'/501'/0'/0/0"
```
```text
Enter mnemonic: patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
pub: 0x810d963eC3C2DfD299E089C42cc7d07089FA5eD2
prv: dd0808b874693a136384b0f4cf923b6caf59a5068bb8bfeb1e38420b86f22f6f
```

Mnemonic is validated and expected to comply to `BIP-39` standard.
Furthermore, a mnemonic in a language different from English is first
translated to English such that the underlying entropy is preserved.

```bash
bip39 translate --to-language=Japanese patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
```
```text
てぶくろ うりきれ てすり あいこくしん ねむい ひりつ こんしゅう うやまう しねん ほうそう らいう そんぞく あつい いわば むかい おおよそ おいこす しちりん でんち はんい さとる やめる けらい みのがす
```

Now using the Japenese mnemonic will result in same keys as those generated using
it's English mnemonic equivalent:
```bash
ethkey gen --mnemonic-language=Japanese てぶくろ うりきれ てすり あいこくしん ねむい ひりつ こんしゅう うやまう しねん ほうそう らいう そんぞく あつい いわば むかい おおよそ おいこす しちりん でんち はんい さとる やめる けらい みのがす
```
```yaml
prv: 12bea0236685f934e39e27bef1793966800f882912c099b4c584a8fbfd28b6e1
pub: 0x33009656efD3e4eA8763B72B020C2327Dee0B2db
```

However, an arbitrary mnemonic can be used by switching off validation

```bash
ethkey gen --skip-mnemonic-validation
```
```text
Enter mnemonic: this is an invalid mnemonic
pub: 0xc1541003D25C206873F5c28Ea684E1072026FC9A
prv: fc30faae343aba4a5151a56881e9c4ff332563fa8eedda661a7681db4ea604bb
```

> It is a good practice to use valid mnemonics and also enter them
> via STDIN to avoid getting them captured in command history

## validate keys
Key validation checks for key string format, length and other characteristics.
For instance, if a private key is entered, it also checks if a public key
can be derived from it.

```bash
ethkey gen
```
```text
Enter mnemonic: few happy dragon spray much obvious total drive hat brain impose bright test there outside peasant share kitchen prefer inmate moment cactus forward crisp
pub: 0x0dD73dE2AC23E6f4928D582aa6510144790DA88e
prv: 70f18df4c72c80a4c0bd47c6ec1b61dd4251c9a77104150c891f8d27a96beb73
```

These keys can be validated:
```bash
ethkey validate \
  --key=0x0dD73dE2AC23E6f4928D582aa6510144790DA88e \
  --output-format=yaml
```
```yaml
valid: true
```

```bash
ethkey validate \
  --key=70f18df4c72c80a4c0bd47c6ec1b61dd4251c9a77104150c891f8d27a96beb73 \
  --output-format=yaml
```
```yaml
valid: true
```

## generate hash
Hash can be generated for an input
```bash
ethkey hash this arbitrary input \
  --output-format=yaml
```
```yaml
hash: 9PW5sgZmMnaBYgJxUQASyDQoeKoxPcgBLvCJEHVEFqb5
```

## sign hash
Hash generated in previous step can be signed using private key
```bash
ethkey sign \
  --key=70f18df4c72c80a4c0bd47c6ec1b61dd4251c9a77104150c891f8d27a96beb73 \
  --hash=9PW5sgZmMnaBYgJxUQASyDQoeKoxPcgBLvCJEHVEFqb5 \
  --output-format=yaml
```
```yaml
sign: GTpQzCsxorDxhWykp4CJPdimV5VKQ9Jqm1xsJsusYGFCuDKHVjVG5dCeVZG8Qu27AfTYUShewRiT5tERohQUMhbTu
```

## verify signature
```bash
ethkey verify \
  --key=0x0dD73dE2AC23E6f4928D582aa6510144790DA88e \
  --hash=9PW5sgZmMnaBYgJxUQASyDQoeKoxPcgBLvCJEHVEFqb5 \
  --sign=GTpQzCsxorDxhWykp4CJPdimV5VKQ9Jqm1xsJsusYGFCuDKHVjVG5dCeVZG8Qu27AfTYUShewRiT5tERohQUMhbTu \
  --output-format=yaml
```
```yaml
verified: true
```

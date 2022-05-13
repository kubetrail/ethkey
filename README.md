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
seed: 83302d7a03b461c84f9c8b80ddc2b5c9f46323b567b8aa9b92c82d2b710abd4a6662b8140ec61b4a993bf7e9d673606efa13c30427f1f746326991b57c53d287
prvHex: 12bea0236685f934e39e27bef1793966800f882912c099b4c584a8fbfd28b6e1
pubHex: 1cdc6ce2a8113342dae3b3c76f298b5aeeb901f3fa9043332ff4e79dd55b893ca1ffa1982934eb537ad49ff40fa9b644b0f7093785cffd7439c5afb2163bc540
addr: 0x33009656efD3e4eA8763B72B020C2327Dee0B2db
keyType: ecdsa
```

Alternatively, pass mnemonic as CLI args
> Please note that passing mnemonic as CLI arg poses a risk of it getting captured
> in command history and thereby stored on the disk, which is a security risk.
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
```yaml
Enter mnemonic: patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
Enter secret passphrase: 
Enter secret passphrase again: 
seed: e393f67697aecc395a187f44439113a7ad8368407c6268b9d1f7cef15f8b01d45d1ba2973515482d3851fe671c8d43c8a3fc2aac7e8320b74ac0edc80d615ca0
prvHex: 2978a24a1ad0e918aa55d9ef6925b7066dd9ceb2b0f843c25bea10c83df97ba5
pubHex: 3bbd51b50b45322bad358a381acfb8db8511525e5384f688130f51c1568c40d5135806a697733438051c3bc26788658a0f6c061cceb5b41bf63af6a8ac6d8412
addr: 0xfa6AEd1c3F820ed5820a75Dd05527d2BDEC5C93B
keyType: ecdsa
```

The chain derivation path can be changed
```bash
ethkey gen --derivation-path="m/44'/501'/0'/0/0" patient board palm abandon right sort find blood grace sweet vote load action bag trash calm burden glow phrase shoot frog vacant elegant tourist
```
```yaml
seed: 83302d7a03b461c84f9c8b80ddc2b5c9f46323b567b8aa9b92c82d2b710abd4a6662b8140ec61b4a993bf7e9d673606efa13c30427f1f746326991b57c53d287
prvHex: dd0808b874693a136384b0f4cf923b6caf59a5068bb8bfeb1e38420b86f22f6f
pubHex: 08a21dcae87a5d574aa5bd2e94968de960aa92190f2c09de7b73c00b1363b653414ee299f9e396b74aae8e964ce3cf313de6d5734236788aec927329f096070c
addr: 0x810d963eC3C2DfD299E089C42cc7d07089FA5eD2
keyType: ecdsa
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
seed: 83302d7a03b461c84f9c8b80ddc2b5c9f46323b567b8aa9b92c82d2b710abd4a6662b8140ec61b4a993bf7e9d673606efa13c30427f1f746326991b57c53d287
prvHex: 12bea0236685f934e39e27bef1793966800f882912c099b4c584a8fbfd28b6e1
pubHex: 1cdc6ce2a8113342dae3b3c76f298b5aeeb901f3fa9043332ff4e79dd55b893ca1ffa1982934eb537ad49ff40fa9b644b0f7093785cffd7439c5afb2163bc540
addr: 0x33009656efD3e4eA8763B72B020C2327Dee0B2db
keyType: ecdsa
```

However, an arbitrary mnemonic can be used by switching off validation

```bash
ethkey gen --skip-mnemonic-validation this is an invalid mnemonic
```
```yaml
seed: bb06e6570ed0b71ac71e4feefeb3a7e2e4cf04ba80a065408150800f86583add8d7ba2ed117444a00f95ca8966ea2e7ff5c8a84b0f5b35a43388d76f0eca043f
prvHex: fc30faae343aba4a5151a56881e9c4ff332563fa8eedda661a7681db4ea604bb
pubHex: ba4db1e7af2dab716f15fc2bab039826d82ff17d1e40bbf8ba292fda8658e913ab72173f8e5a8d22d032bf77db7a86f0e43b6af2ddc93aa976882cab64051120
addr: 0xc1541003D25C206873F5c28Ea684E1072026FC9A
keyType: ecdsa
```

> It is a good practice to use valid mnemonics and also enter them
> via STDIN to avoid getting them captured in command history

## decode keys
Given a private key or public key hex string other info can be derived

```bash
ethkey gen jungle pudding situate maple pattern demand unaware float welcome patrol birth chief
```
```yaml
seed: f85f7dbed5f799fedb12410db708221becd930c21b9364db3426e4fdeb0b537fb188e94e7c12fa3e21c2d796e3a318c435927b013af1483b7ce06a6d282bdbd8
prvHex: 60164407f30a87fae4da0f538b8573b3cfbaa7e96878d848c30df5b5aad617a3
pubHex: 97ad62d8f9f1364d1ab9fb6507124e13ca60ae20ae6704af065eb60c8373355c93639280e5bf99346f31381541299ae255531c8e714e1b36ce07aa4c34e3d94a
addr: 0x10DC542608b5AC7b906800F5773489aEfEeAf9Fe
keyType: ecdsa
```

These keys can be validated:
```bash
ethkey decode 60164407f30a87fae4da0f538b8573b3cfbaa7e96878d848c30df5b5aad617a3
```
```yaml
prvHex: 60164407f30a87fae4da0f538b8573b3cfbaa7e96878d848c30df5b5aad617a3
pubHex: 97ad62d8f9f1364d1ab9fb6507124e13ca60ae20ae6704af065eb60c8373355c93639280e5bf99346f31381541299ae255531c8e714e1b36ce07aa4c34e3d94a
addr: 0x10DC542608b5AC7b906800F5773489aEfEeAf9Fe
keyType: ecdsa
```

```bash
ethkey validate 97ad62d8f9f1364d1ab9fb6507124e13ca60ae20ae6704af065eb60c8373355c93639280e5bf99346f31381541299ae255531c8e714e1b36ce07aa4c34e3d94a
```
```yaml
pubHex: 97ad62d8f9f1364d1ab9fb6507124e13ca60ae20ae6704af065eb60c8373355c93639280e5bf99346f31381541299ae255531c8e714e1b36ce07aa4c34e3d94a
addr: 0x10DC542608b5AC7b906800F5773489aEfEeAf9Fe
keyType: ecdsa
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

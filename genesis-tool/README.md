## Genesis Tool
It provides genesis migration tool to update code before launching localterra.

### How to use

#### Install
```sh
$ make install
```
#### Migrate Codes
```shell
$ LocalTerra migrate-code [code-id]=[path-to-wasm] [code-id]=[path-to-wasm] ... --genesis [path-to-genesis] > migrated_genesis.json
```

#### Migrate state for LocalTerra
```sh
$ LocalTerra migrate-code [code-id]=[path-to-wasm] [code-id]=[path-to-wasm] ... --genesis [path-to-genesis] > migrated_genesis.json
```

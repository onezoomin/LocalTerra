package cmd

/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	wasmvm "github.com/CosmWasm/wasmvm"
)

// migrateCodeCmd represents the migrateCode command
var migrateCodeCmd = &cobra.Command{
	Use:   "migrate-code [codeID=path-to-wasm-file] [codeID=path-to-wasm-file]...",
	Args:  cobra.ArbitraryArgs,
	Short: "overwrite code info of given genesis with new wasm file",
	Long: `The command overwrite the code info of the given code ID from the genesis. 

$ LocalTerra migrate-code [codeID=path-to-wasm-file] [codeID=path-to-wasm-file]...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wasmVM, err := wasmvm.NewVM(
			"build",
			"stargate,staking,terra",
			32,
			false,
			500,
		)
		if err != nil {
			return errors.Wrap(err, "failed to init wasmVM")
		}

		genesisPath, err := cmd.Flags().GetString(flagGenesis)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(genesisPath)
		if err != nil {
			return errors.Wrap(err, "failed to read genesis file")
		}

		var genesis map[string]interface{}
		err = json.Unmarshal(data, &genesis)
		if err != nil {
			return errors.Wrap(err, "failed to parse genesis")
		}

		appState := genesis["app_state"].(map[string]interface{})
		wasmState := appState["wasm"].(map[string]interface{})

		codes := wasmState["codes"].([]interface{})
		for _, val := range args {
			strs := strings.Split(val, "=")
			codeID, err := strconv.ParseInt(strs[0], 10, 32)
			if err != nil {
				return errors.Wrap(err, "failed to parse arguments")
			}

			wasmPath := strs[1]
			wasmFile, err := os.ReadFile(wasmPath)
			if err != nil {
				return errors.Wrapf(err, "failed to read wasm file %s", wasmPath)
			}

			codeHash, err := wasmVM.Create(wasmFile)
			if err != nil {
				return errors.Wrap(err, "failed to compile code")
			}

			code := codes[codeID-1].(map[string]interface{})
			code["code_bytes"] = wasmFile

			codeInfo := code["code_info"].(map[string]interface{})
			codeInfo["code_hash"] = codeHash
		}

		indent, err := cmd.Flags().GetBool(flagIndent)
		if err != nil {
			return err
		}

		var bz []byte
		if indent {
			bz, err = json.MarshalIndent(genesis, "", "\t")
		} else {
			bz, err = json.Marshal(genesis)
		}
		if err != nil {
			return errors.Wrap(err, "failed to marshal genesis")
		}

		fmt.Print(string(bz))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCodeCmd)
}

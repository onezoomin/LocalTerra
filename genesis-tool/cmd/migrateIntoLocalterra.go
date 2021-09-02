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

	"LocalTerra/types"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const validatorNodeID = `DA0317A8E3251C9AEAA38C34820568DCD030CF3F`
const validatorPubKey = `/zGmkgCWRFsJLETAzlzYsbu7EHS5HWpaSyR22rlFM68=`
const validatorAccAddr = `terra1dcegyrekltswvyy0xy69ydgxn9x8x32zdtapd8`
const validatorValAddr = `terravaloper1dcegyrekltswvyy0xy69ydgxn9x8x32zdy3ua5`
const validatorConsAddr = `terravalcons1mgp3028ry5wf464r3s6gyptgmngrpnelhkuyvm`

var localterraAddrs = []string{
	"terra1dcegyrekltswvyy0xy69ydgxn9x8x32zdtapd8",
	"terra1x46rqay4d3cssq8gxxvqz8xt6nwlz4td20k38v",
	"terra17lmam6zguazs5q5u6z5mmx76uj63gldnse2pdp",
	"terra1757tkx08n0cqrw7p86ny9lnxsqeth0wgp0em95",
	"terra199vw7724lzkwz6lf2hsx04lrxfkz09tg8dlp6r",
	"terra18wlvftxzj6zt0xugy2lr9nxzu402690ltaf4ss",
	"terra1e8ryd9ezefuucd4mje33zdms9m2s90m57878v9",
	"terra17tv2hvwpg0ukqgd2y5ct2w54fyan7z0zxrm2f9",
	"terra1lkccuqgj6sjwjn8gsa9xlklqv4pmrqg9dx2fxc",
	"terra1333veey879eeqcff8j3gfcgwt8cfrg9mq20v6f",
	"terra1fmcjjt6yc9wqup2r06urnrd928jhrde6gcld6n",
}

var initialBalance = types.Coins{
	{
		Denom:  "ueur",
		Amount: "10000000000000000",
	},
	{
		Denom:  "ukrw",
		Amount: "1000000000000000000",
	},
	{
		Denom:  "uluna",
		Amount: "1000000000000000",
	},
	{
		Denom:  "usdr",
		Amount: "10000000000000000",
	},
	{
		Denom:  "uusd",
		Amount: "10000000000000000",
	},
}

// migrateIntoLocalterraCmd represents the migrateIntoLocalterra command
var migrateIntoLocalterraCmd = &cobra.Command{
	Use:   "migrate-into-localterra [voting-power]",
	Args:  cobra.ExactArgs(1),
	Short: "Migrate a genesis into localterra",
	Long: `Migrate a genesis into localterra 
	Append a localterra validator and accounts and 
	update chainID to localterra.

$ LocalTerra append-account-set 1000000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		genesisPath, err := cmd.Flags().GetString(flagGenesis)
		if err != nil {
			return err
		}

		votingPower, err := strconv.ParseUint(args[0], 10, 64)
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

		// when the initial height is non-zero,
		// initial_height will be overwritten
		initialHeight, err := cmd.Flags().GetString(flagInitialHeight)
		if err != nil {
			return err
		}

		if initialHeight == "0" {
			initialHeight = genesis["initial_height"].(string)
		} else {
			genesis["initial_height"] = initialHeight
		}

		// make appstate as localterra
		appState := genesis["app_state"].(map[string]interface{})

		// append localterra addresses
		authState := appState["auth"].(map[string]interface{})
		accounts := authState["accounts"].([]interface{})
		bankState := appState["bank"].(map[string]interface{})
		balances := bankState["balances"].([]interface{})

		var newBalances []interface{}
		for _, balance := range balances {
			balance := balance.(map[string]interface{})

			// bonded_token module account balance
			if "terra1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3nln0mh" == balance["address"] {
				coins := balance["coins"].([]interface{})
				for _, coin := range coins {
					coin := coin.(map[string]interface{})
					if coin["denom"] == "uluna" {
						b, err := strconv.ParseUint(coin["amount"].(string), 10, 64)
						if err != nil {
							panic(err)
						}

						coin["amount"] = strconv.FormatUint(b+votingPower*1_000_000, 10)
					}
				}

				balance["coins"] = coins
			}

			for i, addr := range localterraAddrs {
				if addr == balance["address"] {
					balance["coins"] = initialBalance.Add(types.MustParseToCoins(balance["coins"].([]interface{})))

					localterraAddrs[i] = localterraAddrs[len(localterraAddrs)-1]
					localterraAddrs = localterraAddrs[:len(localterraAddrs)-1]
					break
				}
			}

			newBalances = append(newBalances, balance)
		}

		numAccount := len(accounts)
		for i, addr := range localterraAddrs {
			accounts = append(accounts, types.Account{
				Type:          "/cosmos.auth.v1beta1.BaseAccount",
				Address:       addr,
				PubKey:        nil,
				AccountNumber: strconv.FormatInt(int64(numAccount+i), 10),
				Sequence:      "0",
			})

			newBalances = append(newBalances, types.Balance{
				Address: addr,
				Coins:   initialBalance,
			})
		}

		authState["accounts"] = accounts
		bankState["balances"] = newBalances

		supplyState := bankState["supply"].([]interface{})
		incrementSupply := types.Coins{
			{
				Denom:  "uluna",
				Amount: strconv.FormatUint(uint64(11000000000000000)+votingPower*1000000, 10),
			},
			{
				Denom:  "uusd",
				Amount: "110000000000000000",
			},
			{
				Denom:  "ukrw",
				Amount: "11000000000000000000",
			},
			{
				Denom:  "ueur",
				Amount: "110000000000000000",
			},
			{
				Denom:  "usdr",
				Amount: "110000000000000000",
			},
		}

		bankState["supply"] = incrementSupply.Add(types.MustParseToCoins(supplyState))

		// append validator
		validators := genesis["validators"].([]interface{})
		validators = append(validators, types.Validator{
			Address: validatorNodeID,
			Name:    "localterra",
			Power:   strconv.FormatUint(votingPower, 10),
			PubKey: types.ValidatorPubKey{
				Type:  "tendermint/PubKeyEd25519",
				Value: validatorPubKey,
			},
		})

		genesis["validators"] = validators

		// append delegations
		stakingState := appState["staking"].(map[string]interface{})
		delegations := stakingState["delegations"].([]interface{})
		delegations = append(delegations, types.Delegation{
			DelegatorAddress: validatorAccAddr,
			Shares:           strconv.FormatUint(votingPower*1_000_000, 10) + ".000000000000000000",
			ValidatorAddress: validatorValAddr,
		})

		stakingState["delegations"] = delegations

		// update total power
		lastTotalPower, err := strconv.ParseUint(stakingState["last_total_power"].(string), 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse total_power")
		}

		stakingState["last_total_power"] = strconv.FormatUint(lastTotalPower+votingPower, 10)

		// update last validator powers
		lastValidatorPowers := stakingState["last_validator_powers"].([]interface{})
		lastValidatorPowers = append(lastValidatorPowers, types.ValidatorPower{
			Address: validatorValAddr,
			Power:   strconv.FormatUint(votingPower, 10),
		})

		stakingState["last_validator_powers"] = lastValidatorPowers

		// append staking validator
		stakingValidators := stakingState["validators"].([]interface{})
		stakingValidators = append(stakingValidators, types.StakingValidator{
			Commission: types.StakingCommission{
				CommissionRates: types.StakingCommissionRates{
					MaxChangeRate: "0.010000000000000000",
					MaxRate:       "0.200000000000000000",
					Rate:          "0.100000000000000000",
				},
				UpdateTime: "2020-08-24T08:43:02.336889Z",
			},
			ConsensusPubkey: types.StakingConsensusPubkey{
				Type: "/cosmos.crypto.ed25519.PubKey",
				Key:  validatorPubKey,
			},
			DelegatorShares: strconv.FormatUint(votingPower*1_000_000, 10) + ".000000000000000000",
			Description: types.StakingDescription{
				Details:         "",
				Identity:        "",
				Moniker:         "localterra",
				SecurityContact: "",
				Website:         "https://github.com/terra-project/LocalTerra",
			},
			Jailed:            false,
			MinSelfDelegation: "1",
			OperatorAddress:   validatorValAddr,
			Status:            "BOND_STATUS_BONDED",
			Tokens:            strconv.FormatUint(votingPower*1_000_000, 10),
			UnbondingHeight:   "0",
			UnbondingTime:     "1970-01-01T00:00:00Z",
		})
		stakingState["validators"] = stakingValidators

		// register slashing infos
		slashingState := appState["slashing"].(map[string]interface{})
		missedBlocks := slashingState["missed_blocks"].([]interface{})
		missedBlocks = append(missedBlocks, types.MissedBlockInfo{
			Address:      validatorConsAddr,
			MissedBlocks: make([]interface{}, 0),
		})
		slashingState["missed_blocks"] = missedBlocks

		signingInfos := slashingState["signing_infos"].([]interface{})
		signingInfos = append(signingInfos, types.SigningInfo{
			Address: validatorConsAddr,
			ValidatorSigingInfo: types.ValidatorSigningInfo{
				Address:             validatorConsAddr,
				IndexOffset:         "0",
				JailedUntil:         "1970-01-01T00:00:00Z",
				MissedBlocksCounter: "0",
				StartHeight:         "0",
				Tombstoned:          false,
			},
		})
		slashingState["signing_infos"] = signingInfos

		// register distribution infos
		distributionState := appState["distribution"].(map[string]interface{})
		delegatorStartingInfos := distributionState["delegator_starting_infos"].([]interface{})
		delegatorStartingInfos = append(delegatorStartingInfos, types.DelegatorStartingInfo{
			DelegatorAddress: validatorAccAddr,
			ValidatorAddress: validatorValAddr,
			StartingInfo: types.StartingInfo{
				Height:         initialHeight,
				PreviousPeriod: "1",
				Stake:          strconv.FormatUint(votingPower*1_000_000, 10) + ".000000000000000000",
			},
		})
		distributionState["delegator_starting_infos"] = delegatorStartingInfos

		outstandingRewards := distributionState["outstanding_rewards"].([]interface{})
		outstandingRewards = append(outstandingRewards, types.OutstandingRewards{
			ValidatorAddress:   validatorValAddr,
			OutstandingRewards: make([]interface{}, 0),
		})
		distributionState["outstanding_rewards"] = outstandingRewards

		validatorAccumulatedCommissions := distributionState["validator_accumulated_commissions"].([]interface{})
		validatorAccumulatedCommissions = append(validatorAccumulatedCommissions, types.ValidatorAccumulatedCommission{
			ValidatorAddress: validatorValAddr,
			Accumulated: types.AccumulatedCommission{
				Commission: make([]interface{}, 0),
			},
		})
		distributionState["validator_accumulated_commissions"] = validatorAccumulatedCommissions

		validatorCurrentRewards := distributionState["validator_current_rewards"].([]interface{})
		validatorCurrentRewards = append(validatorCurrentRewards, types.ValidatorCurrentReward{
			Rewards: types.RewardInfo{
				Period:  "2",
				Rewards: make([]interface{}, 0),
			},
			ValidatorAddress: validatorValAddr,
		})
		distributionState["validator_current_rewards"] = validatorCurrentRewards

		validatorHistoricalRewards := distributionState["validator_historical_rewards"].([]interface{})
		validatorHistoricalRewards = append(validatorHistoricalRewards, types.ValidatorHistoricalReward{
			Period:           "1",
			ValidatorAddress: validatorValAddr,
			Rewards: types.HistoricalRewardInfo{
				CumulativeRewardRatio: make([]interface{}, 0),
				ReferenceCount:        2,
			},
		})
		distributionState["validator_historical_rewards"] = validatorHistoricalRewards

		genesis["chain_id"] = "localterra"

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

const flagInitialHeight = "initial-height"

func init() {
	rootCmd.AddCommand(migrateIntoLocalterraCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	migrateIntoLocalterraCmd.Flags().String(flagInitialHeight, "0", "non-zero optional initial height will overwrite the initial height")
}

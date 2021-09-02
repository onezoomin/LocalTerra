package types

type Account struct {
	Type          string      `json:"@type"`
	Address       string      `json:"address"`
	PubKey        interface{} `json:"pub_key"`
	AccountNumber string      `json:"account_number"`
	Sequence      string      `json:"sequence"`
}

type Balance struct {
	Address string `json:"address"`
	Coins   []Coin `json:"coins"`
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Validator struct {
	Address string          `json:"address"`
	Name    string          `json:"name"`
	Power   string          `json:"power"`
	PubKey  ValidatorPubKey `json:"pub_key"`
}

type ValidatorPubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	Shares           string `json:"shares"`
	ValidatorAddress string `json:"validator_address"`
}

type ValidatorPower struct {
	Address string `json:"address"`
	Power   string `json:"power"`
}

type StakingValidator struct {
	Commission        StakingCommission      `json:"commission"`
	ConsensusPubkey   StakingConsensusPubkey `json:"consensus_pubkey"`
	DelegatorShares   string                 `json:"delegator_shares"`
	Description       StakingDescription     `json:"description"`
	Jailed            bool                   `json:"jailed"`
	MinSelfDelegation string                 `json:"min_self_delegation"`
	OperatorAddress   string                 `json:"operator_address"`
	Status            string                 `json:"status"`
	Tokens            string                 `json:"tokens"`
	UnbondingHeight   string                 `json:"unbonding_height"`
	UnbondingTime     string                 `json:"unbonding_time"`
}

type StakingConsensusPubkey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type StakingDescription struct {
	Details         string `json:"details"`
	Identity        string `json:"identity"`
	Moniker         string `json:"moniker"`
	SecurityContact string `json:"security_contact"`
	Website         string `json:"website"`
}

type StakingCommission struct {
	CommissionRates StakingCommissionRates `json:"commission_rates"`
	UpdateTime      string                 `json:"update_time"`
}

type StakingCommissionRates struct {
	MaxChangeRate string `json:"max_change_rate"`
	MaxRate       string `json:"max_rate"`
	Rate          string `json:"rate"`
}

type SigningInfo struct {
	Address             string               `json:"address"`
	ValidatorSigingInfo ValidatorSigningInfo `json:"validator_signing_info"`
}

type ValidatorSigningInfo struct {
	Address             string `json:"address"`
	IndexOffset         string `json:"index_offset"`
	JailedUntil         string `json:"jailed_until"`
	MissedBlocksCounter string `json:"missed_blocks_counter"`
	StartHeight         string `json:"start_height"`
	Tombstoned          bool   `json:"tombstoned"`
}

type MissedBlockInfo struct {
	Address      string        `json:"address"`
	MissedBlocks []interface{} `json:"missed_blocks"`
}

type DelegatorStartingInfo struct {
	DelegatorAddress string       `json:"delegator_address"`
	ValidatorAddress string       `json:"validator_address"`
	StartingInfo     StartingInfo `json:"starting_info"`
}

type StartingInfo struct {
	Height         string `json:"height"`
	PreviousPeriod string `json:"previous_period"`
	Stake          string `json:"stake"`
}

type OutstandingRewards struct {
	ValidatorAddress   string        `json:"validator_address"`
	OutstandingRewards []interface{} `json:"outstanding_rewards"`
}

type ValidatorAccumulatedCommission struct {
	ValidatorAddress string                `json:"validator_address"`
	Accumulated      AccumulatedCommission `json:"accumulated"`
}

type AccumulatedCommission struct {
	Commission []interface{} `json:"commission"`
}

type ValidatorCurrentReward struct {
	Rewards          RewardInfo `json:"rewards"`
	ValidatorAddress string     `json:"validator_address"`
}

type RewardInfo struct {
	Period  string        `json:"period"`
	Rewards []interface{} `json:"rewards"`
}

type ValidatorHistoricalReward struct {
	Period           string               `json:"period"`
	ValidatorAddress string               `json:"validator_address"`
	Rewards          HistoricalRewardInfo `json:"rewards"`
}

type HistoricalRewardInfo struct {
	CumulativeRewardRatio []interface{} `json:"cumulative_reward_ratio"`
	ReferenceCount        int           `json:"reference_count"`
}

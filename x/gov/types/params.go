package types

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	params "github.com/cosmos/cosmos-sdk/x/params/types"
)

// parameter store keys
var (
	ParamStoreKeyDepositParams = []byte("depositparams")
	ParamStoreKeyVotingParams  = []byte("votingparams")
	ParamStoreKeyTallyParams   = []byte("tallyparams")
	CertVotesKeyPrefix         = []byte("certvote")
)

// ParamKeyTable is the key declaration for parameters.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamStoreKeyDepositParams, DepositParams{}, validateDepositParams),
		params.NewParamSetPair(ParamStoreKeyVotingParams, govTypes.VotingParams{}, validateVotingParams),
		params.NewParamSetPair(ParamStoreKeyTallyParams, TallyParams{}, validateTally),
	)
}

// NewDepositParams creates a new DepositParams object
func NewDepositParams(minInitialDeposit, minDeposit sdk.Coins, maxDepositPeriod time.Duration) DepositParams {
	return DepositParams{
		MinInitialDeposit: minInitialDeposit,
		MinDeposit:        minDeposit,
		MaxDepositPeriod:  maxDepositPeriod,
	}
}

func (dp DepositParams) String() string {
	return fmt.Sprintf(`Deposit Params:
  Min Initial Deposit: %s
  Min Deposit:         %s
  Max Deposit Period:  %s`, dp.MinInitialDeposit, dp.MinDeposit, dp.MaxDepositPeriod)
}

// Equal checks equality of DepositParams
func (dp DepositParams) Equal(dp2 DepositParams) bool {
	return dp.MinInitialDeposit.IsEqual(dp2.MinInitialDeposit) &&
		dp.MinDeposit.IsEqual(dp2.MinDeposit) &&
		dp.MaxDepositPeriod == dp2.MaxDepositPeriod
}

func validateDepositParams(i interface{}) error {
	v, ok := i.(DepositParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.MinDeposit.IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", v.MinDeposit)
	}
	if v.MaxDepositPeriod <= 0 {
		return fmt.Errorf("maximum deposit period must be positive: %d", v.MaxDepositPeriod)
	}

	return nil
}

// Params returns all the governance params
type Params struct {
	VotingParams  govTypes.VotingParams `json:"voting_params" yaml:"voting_params"`
	TallyParams   TallyParams           `json:"tally_params" yaml:"tally_params"`
	DepositParams DepositParams         `json:"deposit_params" yaml:"deposit_parmas"`
}

func (gp Params) String() string {
	return gp.VotingParams.String() + "\n" +
		gp.TallyParams.String() + "\n" + gp.DepositParams.String()
}

// NewParams returns a Params structs including voting, deposit and tally params
func NewParams(vp govTypes.VotingParams, tp TallyParams, dp DepositParams) Params {
	return Params{
		VotingParams:  vp,
		DepositParams: dp,
		TallyParams:   tp,
	}
}

func (tp TallyParams) String() string {
	b, err := json.MarshalIndent(tp, "", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func validateTally(i interface{}) error {
	v, ok := i.(TallyParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := validateTallyParams(*v.CertifierUpdateSecurityVoteTally); err != nil {
		return err
	}
	if err := validateTallyParams(*v.CertifierUpdateStakeVoteTally); err != nil {
		return err
	}
	if err := validateTallyParams(*v.DefaultTally); err != nil {
		return err
	}
	return nil
}

func validateTallyParams(tallyParams govTypes.TallyParams) error {
	if tallyParams.Quorum.IsNegative() {
		return fmt.Errorf("quorom cannot be negative: %s", tallyParams.Quorum)
	}
	if tallyParams.Quorum.GT(sdk.OneDec()) {
		return fmt.Errorf("quorom too large: %s", tallyParams)
	}
	if !tallyParams.Threshold.IsPositive() {
		return fmt.Errorf("vote threshold must be positive: %s", tallyParams.Threshold)
	}
	if tallyParams.Threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold too large: %s", tallyParams)
	}
	if !tallyParams.VetoThreshold.IsPositive() {
		return fmt.Errorf("veto threshold must be positive: %s", tallyParams.Threshold)
	}
	if tallyParams.VetoThreshold.GT(sdk.OneDec()) {
		return fmt.Errorf("veto threshold too large: %s", tallyParams)
	}

	return nil
}

func validateVotingParams(i interface{}) error {
	v, ok := i.(govTypes.VotingParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.VotingPeriod <= 0 {
		return fmt.Errorf("voting period must be positive: %s", v.VotingPeriod)
	}

	return nil
}

// CertVotesKey gets the first part of the cert votes key based on the proposalID
func CertVotesKey(proposalID uint64) []byte {
	return append(CertVotesKeyPrefix, govTypes.GetProposalIDBytes(proposalID)...)
}

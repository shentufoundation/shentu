package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	params "github.com/cosmos/cosmos-sdk/x/params/types"
)

const ParamCustom = "custom"

// parameter store keys
var (
	ParamStoreKeyCustomParams = []byte("customparams")
	CertVotesKeyPrefix        = []byte("certvote")
)

// ParamKeyTable is the key declaration for parameters.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(govTypes.ParamStoreKeyDepositParams, govTypes.DepositParams{}, validateDepositParams),
		params.NewParamSetPair(govTypes.ParamStoreKeyVotingParams, govTypes.VotingParams{}, validateVotingParams),
		params.NewParamSetPair(govTypes.ParamStoreKeyTallyParams, govTypes.TallyParams{}, validateTally),
		params.NewParamSetPair(ParamStoreKeyCustomParams, CustomParams{}, validateCustomParams),
	)
}

func validateDepositParams(i interface{}) error {
	v, ok := i.(govTypes.DepositParams)
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
	VotingParams  govTypes.VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams   govTypes.TallyParams   `json:"tally_params" yaml:"tally_params"`
	DepositParams govTypes.DepositParams `json:"deposit_params" yaml:"deposit_parmas"`
	CustomParams  CustomParams           `json:"custom_params" yaml:"custom_params"`
}

func (gp Params) String() string {
	return gp.VotingParams.String() + "\n" +
		gp.TallyParams.String() + "\n" + gp.DepositParams.String() + "\n" +
		gp.CustomParams.String()
}

// NewParams returns a Params structs including voting, deposit and tally params
func NewParams(vp govTypes.VotingParams, tp govTypes.TallyParams, dp govTypes.DepositParams, cp CustomParams) Params {
	return Params{
		VotingParams:  vp,
		DepositParams: dp,
		TallyParams:   tp,
		CustomParams:  cp,
	}
}

func validateTally(i interface{}) error {
	v, ok := i.(govTypes.TallyParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := validateTallyParams(v); err != nil {
		return err
	}
	return nil
}

// String implements stringer insterface
func (cp CustomParams) String() string {
	out, _ := yaml.Marshal(cp)
	return string(out)
}

func validateCustomParams(i interface{}) error {
	v, ok := i.(CustomParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := validateTallyParams(*v.CertifierUpdateSecurityVoteTally); err != nil {
		return err
	}
	if err := validateTallyParams(*v.CertifierUpdateStakeVoteTally); err != nil {
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

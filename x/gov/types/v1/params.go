package v1

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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
		params.NewParamSetPair(govtypesv1.ParamStoreKeyDepositParams, govtypesv1.DepositParams{}, validateDepositParams),
		params.NewParamSetPair(govtypesv1.ParamStoreKeyVotingParams, govtypesv1.VotingParams{}, validateVotingParams),
		params.NewParamSetPair(govtypesv1.ParamStoreKeyTallyParams, govtypesv1.TallyParams{}, validateTally),
		params.NewParamSetPair(ParamStoreKeyCustomParams, CustomParams{}, validateCustomParams),
	)
}

func validateDepositParams(i interface{}) error {
	v, ok := i.(*govtypesv1.DepositParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !sdk.Coins(v.MinDeposit).IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", v.MinDeposit)
	}

	if v.MaxDepositPeriod == nil || v.MaxDepositPeriod.Seconds() <= 0 {
		return fmt.Errorf("maximum deposit period must be positive: %d", v.MaxDepositPeriod)
	}

	return nil
}

// Params returns all the governance params
type Params struct {
	VotingParams  govtypesv1.VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams   govtypesv1.TallyParams   `json:"tally_params" yaml:"tally_params"`
	DepositParams govtypesv1.DepositParams `json:"deposit_params" yaml:"deposit_parmas"`
	CustomParams  CustomParams             `json:"custom_params" yaml:"custom_params"`
}

func (gp Params) String() string {
	return gp.VotingParams.String() + "\n" +
		gp.TallyParams.String() + "\n" + gp.DepositParams.String() + "\n" +
		gp.CustomParams.String()
}

// NewParams returns a Params structs including voting, deposit and tally params
func NewParams(vp govtypesv1.VotingParams, tp govtypesv1.TallyParams, dp govtypesv1.DepositParams, cp CustomParams) Params {
	return Params{
		VotingParams:  vp,
		DepositParams: dp,
		TallyParams:   tp,
		CustomParams:  cp,
	}
}

func validateTally(i interface{}) error {
	v, ok := i.(*govtypesv1.TallyParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := validateTallyParams(v); err != nil {
		return err
	}
	return nil
}

func validateCustomParams(i interface{}) error {
	v, ok := i.(CustomParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := validateTallyParams(v.CertifierUpdateSecurityVoteTally); err != nil {
		return err
	}
	if err := validateTallyParams(v.CertifierUpdateStakeVoteTally); err != nil {
		return err
	}

	return nil
}

func validateTallyParams(i interface{}) error {
	v, ok := i.(*govtypesv1.TallyParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	quorum, err := sdk.NewDecFromStr(v.Quorum)
	if err != nil {
		return fmt.Errorf("invalid quorum string: %w", err)
	}
	if quorum.IsNegative() {
		return fmt.Errorf("quorom cannot be negative: %s", quorum)
	}
	if quorum.GT(sdk.OneDec()) {
		return fmt.Errorf("quorom too large: %s", v)
	}

	threshold, err := sdk.NewDecFromStr(v.Threshold)
	if err != nil {
		return fmt.Errorf("invalid threshold string: %w", err)
	}
	if !threshold.IsPositive() {
		return fmt.Errorf("vote threshold must be positive: %s", threshold)
	}
	if threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold too large: %s", v)
	}

	vetoThreshold, err := sdk.NewDecFromStr(v.VetoThreshold)
	if err != nil {
		return fmt.Errorf("invalid vetoThreshold string: %w", err)
	}
	if !vetoThreshold.IsPositive() {
		return fmt.Errorf("veto threshold must be positive: %s", vetoThreshold)
	}
	if vetoThreshold.GT(sdk.OneDec()) {
		return fmt.Errorf("veto threshold too large: %s", v)
	}

	return nil
}

func validateVotingParams(i interface{}) error {
	v, ok := i.(*govtypesv1.VotingParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.VotingPeriod == nil || v.VotingPeriod.Seconds() <= 0 {
		return fmt.Errorf("voting period must be positive: %s", v.VotingPeriod)
	}

	return nil
}

// CertVotesKey gets the first part of the cert votes key based on the proposalID
func CertVotesKey(proposalID uint64) []byte {
	return append(CertVotesKeyPrefix, govtypes.GetProposalIDBytes(proposalID)...)
}

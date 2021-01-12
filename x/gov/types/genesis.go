package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/common"
)

// GenesisState defines the governance genesis state.
type GenesisState struct {
	StartingProposalID uint64                `json:"starting_proposal_id" yaml:"starting_proposal_id"`
	Deposits           Deposits              `json:"deposits" yaml:"deposits"`
	Votes              Votes                 `json:"votes" yaml:"votes"`
	Proposals          Proposals             `json:"proposals" yaml:"proposals"`
	DepositParams      DepositParams         `json:"deposit_params" yaml:"deposit_params"`
	VotingParams       govTypes.VotingParams `json:"voting_params" yaml:"voting_params"`
	TallyParams        TallyParams           `json:"tally_params" yaml:"tally_params"`
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() GenesisState {
	minInitialDepositTokens := sdk.TokensFromConsensusPower(0)
	minDepositTokens := sdk.TokensFromConsensusPower(512)
	return GenesisState{
		StartingProposalID: govTypes.DefaultStartingProposalID,
		DepositParams: DepositParams{
			MinInitialDeposit: sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, minInitialDepositTokens)},
			MinDeposit:        sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, minDepositTokens)},
			MaxDepositPeriod:  govTypes.DefaultPeriod,
		},
		VotingParams: govTypes.DefaultVotingParams(),
		TallyParams: TallyParams{
			DefaultTally: govTypes.TallyParams{
				Quorum:    sdk.NewDecWithPrec(334, 3),
				Threshold: sdk.NewDecWithPrec(5, 1),
				Veto:      sdk.NewDecWithPrec(334, 3),
			},
			CertifierUpdateSecurityVoteTally: govTypes.TallyParams{
				Quorum:    sdk.NewDecWithPrec(334, 3),
				Threshold: sdk.NewDecWithPrec(667, 3),
				Veto:      sdk.NewDecWithPrec(334, 3),
			},
			CertifierUpdateStakeVoteTally: govTypes.TallyParams{
				Quorum:    sdk.NewDecWithPrec(334, 3),
				Threshold: sdk.NewDecWithPrec(9, 1),
				Veto:      sdk.NewDecWithPrec(334, 3),
			},
		},
	}
}

// ValidateGenesis validates crisis genesis data.
func ValidateGenesis(data GenesisState) error {
	for _, tp := range []govTypes.TallyParams{
		data.TallyParams.DefaultTally,
		data.TallyParams.CertifierUpdateStakeVoteTally,
		data.TallyParams.CertifierUpdateSecurityVoteTally,
	} {
		threshold := tp.Threshold
		if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
			return fmt.Errorf("governance vote threshold should be positive and less or equal to one, is %s",
				threshold.String())
		}

		veto := tp.Veto
		if veto.IsNegative() || veto.GT(sdk.OneDec()) {
			return fmt.Errorf("governance vote veto threshold should be positive and less or equal to one, is %s",
				veto.String())
		}
	}

	if !data.DepositParams.MinDeposit.IsValid() {
		return fmt.Errorf("governance deposit amount must be a valid sdk.Coins amount, is %s",
			data.DepositParams.MinDeposit.String())
	}

	return nil
}

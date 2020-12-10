package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/common"
)

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	minInitialDepositTokens := sdk.TokensFromConsensusPower(0)
	minDepositTokens := sdk.TokensFromConsensusPower(512)

	// quorum, threshold, and veto threshold params
	defaultTally := govTypes.NewTallyParams(sdk.NewDecWithPrec(334, 3), sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(334, 3))
	certifierUpdateSecurityVoteTally := govTypes.NewTallyParams(sdk.NewDecWithPrec(334, 3), sdk.NewDecWithPrec(667, 3), sdk.NewDecWithPrec(334, 3))
	certifierUpdateStakeVoteTally := govTypes.NewTallyParams(sdk.NewDecWithPrec(334, 3), sdk.NewDecWithPrec(9, 1), sdk.NewDecWithPrec(334, 3))
	
	return &GenesisState{
		StartingProposalId: govTypes.DefaultStartingProposalID,
		DepositParams: DepositParams{
			MinInitialDeposit: sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, minInitialDepositTokens)},
			MinDeposit:        sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, minDepositTokens)},
			MaxDepositPeriod:  govTypes.DefaultPeriod,
		},
		VotingParams: govTypes.DefaultVotingParams(),
		TallyParams: TallyParams{
			DefaultTally: &defaultTally,
			CertifierUpdateSecurityVoteTally: &certifierUpdateSecurityVoteTally,
			CertifierUpdateStakeVoteTally: &certifierUpdateStakeVoteTally,
		},
	}
}

// ValidateGenesis validates gov genesis data.
func ValidateGenesis(data *GenesisState) error {
	err := validateTallyParams(*data.TallyParams.DefaultTally)
	if err != nil {
		return err
	}
	err = validateTallyParams(*data.TallyParams.CertifierUpdateStakeVoteTally)
	if err != nil {
		return err
	}
	err = validateTallyParams(*data.TallyParams.CertifierUpdateSecurityVoteTally)
	if err != nil {
		return err
	}
	
	if !data.DepositParams.MinDeposit.IsValid() {
		return fmt.Errorf("governance deposit amount must be a valid sdk.Coins amount, is %s",
			data.DepositParams.MinDeposit.String())
	}

	return nil
}

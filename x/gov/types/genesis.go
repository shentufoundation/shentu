package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
			MinInitialDeposit: sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, minInitialDepositTokens)},
			MinDeposit:        sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, minDepositTokens)},
			MaxDepositPeriod:  govTypes.DefaultPeriod,
		},
		VotingParams: govTypes.DefaultVotingParams(),
		TallyParams: TallyParams{
			DefaultTally:                     &defaultTally,
			CertifierUpdateSecurityVoteTally: &certifierUpdateSecurityVoteTally,
			CertifierUpdateStakeVoteTally:    &certifierUpdateStakeVoteTally,
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

var _ types.UnpackInterfacesMessage = GenesisState{}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (data GenesisState) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	for _, p := range data.Proposals {
		err := p.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}

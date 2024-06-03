package v1

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/common"
)

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	minDepositTokens := sdk.TokensFromConsensusPower(512, sdk.DefaultPowerReduction)

	// quorum, threshold, and veto threshold params
	defaultTally := govtypesv1.NewTallyParams(sdk.NewDecWithPrec(334, 3), sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(334, 3))
	certifierUpdateSecurityVoteTally := govtypesv1.NewTallyParams(sdk.NewDecWithPrec(334, 3), sdk.NewDecWithPrec(667, 3), sdk.NewDecWithPrec(334, 3))
	certifierUpdateStakeVoteTally := govtypesv1.NewTallyParams(sdk.NewDecWithPrec(334, 3), sdk.NewDecWithPrec(9, 1), sdk.NewDecWithPrec(334, 3))

	defaultPeriod := govtypesv1.DefaultPeriod
	votingParams := govtypesv1.DefaultVotingParams()

	return &GenesisState{
		StartingProposalId: govtypesv1.DefaultStartingProposalID,
		DepositParams: &govtypesv1.DepositParams{
			MinDeposit:       sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, minDepositTokens)),
			MaxDepositPeriod: &defaultPeriod,
		},
		VotingParams: &votingParams,
		TallyParams:  &defaultTally,
		CustomParams: &CustomParams{
			CertifierUpdateSecurityVoteTally: &certifierUpdateSecurityVoteTally,
			CertifierUpdateStakeVoteTally:    &certifierUpdateStakeVoteTally,
		},
	}
}

// ValidateGenesis validates gov genesis data.
func ValidateGenesis(data *GenesisState) error {
	err := validateTallyParams(data.TallyParams)
	if err != nil {
		return err
	}
	err = validateTallyParams(data.CustomParams.CertifierUpdateStakeVoteTally)
	if err != nil {
		return err
	}
	err = validateTallyParams(data.CustomParams.CertifierUpdateSecurityVoteTally)
	if err != nil {
		return err
	}
	err = validateDepositParams(data.DepositParams)
	if err != nil {
		return fmt.Errorf("governance deposit amount must be a valid sdk.Coins amount, is %s",
			data.DepositParams.MinDeposit)
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

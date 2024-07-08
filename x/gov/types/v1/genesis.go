package v1

import (
	"errors"

	"github.com/shentufoundation/shentu/v2/common"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	minDepositTokens := sdk.TokensFromConsensusPower(512, sdk.DefaultPowerReduction)

	// quorum, threshold, and veto threshold params
	certifierUpdateSecurityVoteTally := govtypesv1.NewTallyParams(sdk.NewDecWithPrec(334, 3).String(), sdk.NewDecWithPrec(667, 3).String(), sdk.NewDecWithPrec(334, 3).String())
	certifierUpdateStakeVoteTally := govtypesv1.NewTallyParams(sdk.NewDecWithPrec(334, 3).String(), sdk.NewDecWithPrec(9, 1).String(), sdk.NewDecWithPrec(334, 3).String())

	param := govtypesv1.DefaultParams()
	param.MinDeposit = sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, minDepositTokens))
	return &GenesisState{
		StartingProposalId: govtypesv1.DefaultStartingProposalID,
		Params:             &param,
		CustomParams: &CustomParams{
			CertifierUpdateSecurityVoteTally: &certifierUpdateSecurityVoteTally,
			CertifierUpdateStakeVoteTally:    &certifierUpdateStakeVoteTally,
		},
	}
}

// ValidateGenesis validates gov genesis data.
func ValidateGenesis(data *GenesisState) error {
	if data.StartingProposalId == 0 {
		return errors.New("starting proposal id must be greater than 0")
	}

	//var errGroup errgroup.Group
	//
	//// weed out duplicate proposals
	//proposalIds := make(map[uint64]struct{})
	//for _, p := range data.Proposals {
	//	if _, ok := proposalIds[p.Id]; ok {
	//		return fmt.Errorf("duplicate proposal id: %d", p.Id)
	//	}
	//
	//	proposalIds[p.Id] = struct{}{}
	//}
	//
	//err := validateTallyParams(data.TallyParams)
	//if err != nil {
	//	return err
	//}
	//err = validateTallyParams(data.CustomParams.CertifierUpdateStakeVoteTally)
	//if err != nil {
	//	return err
	//}
	//err = validateTallyParams(data.CustomParams.CertifierUpdateSecurityVoteTally)
	//if err != nil {
	//	return err
	//}
	//err = validateDepositParams(data.DepositParams)
	//if err != nil {
	//	return fmt.Errorf("governance deposit amount must be a valid sdk.Coins amount, is %s",
	//		data.DepositParams.MinDeposit)
	//}

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

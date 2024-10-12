package v1

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	"golang.org/x/sync/errgroup"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/common"
)

// NewGenesisState creates a new genesis state for the governance module
func NewGenesisState(startingProposalID uint64, params govtypesv1.Params, customParams CustomParams) *GenesisState {
	return &GenesisState{
		StartingProposalId: startingProposalID,
		Params:             &params,
		CustomParams:       &customParams,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	minDepositTokens := sdk.TokensFromConsensusPower(512, sdk.DefaultPowerReduction)

	// quorum, threshold, and veto threshold params
	certifierUpdateSecurityVoteTally := govtypesv1.NewTallyParams(math.LegacyNewDecWithPrec(334, 3).String(), math.LegacyNewDecWithPrec(667, 3).String(), math.LegacyNewDecWithPrec(334, 3).String())
	certifierUpdateStakeVoteTally := govtypesv1.NewTallyParams(math.LegacyNewDecWithPrec(334, 3).String(), math.LegacyNewDecWithPrec(9, 1).String(), math.LegacyNewDecWithPrec(334, 3).String())

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

	var errGroup errgroup.Group

	// weed out duplicate proposals
	proposalIds := make(map[uint64]struct{})
	for _, p := range data.Proposals {
		if _, ok := proposalIds[p.Id]; ok {
			return fmt.Errorf("duplicate proposal id: %d", p.Id)
		}

		proposalIds[p.Id] = struct{}{}
	}

	// weed out duplicate deposits
	errGroup.Go(func() error {
		type depositKey struct {
			ProposalId uint64 //nolint:revive // staying consistent with main and v0.47
			Depositor  string
		}
		depositIds := make(map[depositKey]struct{})
		for _, d := range data.Deposits {
			if _, ok := proposalIds[d.ProposalId]; !ok {
				return fmt.Errorf("deposit %v has non-existent proposal id: %d", d, d.ProposalId)
			}

			dk := depositKey{d.ProposalId, d.Depositor}
			if _, ok := depositIds[dk]; ok {
				return fmt.Errorf("duplicate deposit: %v", d)
			}

			depositIds[dk] = struct{}{}
		}

		return nil
	})

	// weed out duplicate votes
	errGroup.Go(func() error {
		type voteKey struct {
			ProposalId uint64 //nolint:revive // staying consistent with main and v0.47
			Voter      string
		}
		voteIds := make(map[voteKey]struct{})
		for _, v := range data.Votes {
			if _, ok := proposalIds[v.ProposalId]; !ok {
				return fmt.Errorf("vote %v has non-existent proposal id: %d", v, v.ProposalId)
			}

			vk := voteKey{v.ProposalId, v.Voter}
			if _, ok := voteIds[vk]; ok {
				return fmt.Errorf("duplicate vote: %v", v)
			}

			voteIds[vk] = struct{}{}
		}

		return nil
	})

	// verify params
	errGroup.Go(func() error {
		return data.Params.ValidateBasic()
	})

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

	return errGroup.Wait()
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

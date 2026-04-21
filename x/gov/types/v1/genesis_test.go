package v1_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// ValidateGenesis is the schema-layer gate: shentud validate-genesis
// must reject a bundled cert-update proposal instead of letting the
// node start and panic at InitGenesis.
func TestValidateGenesis_RejectsBundledCertUpdateProposal(t *testing.T) {
	authority := sdk.AccAddress("auth________________")
	recipient := sdk.AccAddress("recipient___________")
	certifier := sdk.AccAddress("certifier___________")

	certMsg := certtypes.NewMsgUpdateCertifier(authority, certifier, "add", certtypes.Add)
	sendMsg := banktypes.NewMsgSend(authority, recipient, sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))))

	certAny, err := types.NewAnyWithValue(certMsg)
	require.NoError(t, err)
	sendAny, err := types.NewAnyWithValue(sendMsg)
	require.NoError(t, err)

	now := time.Unix(100, 0)
	end := now.Add(48 * time.Hour)
	params := govtypesv1.DefaultParams()
	gs := &typesv1.GenesisState{
		StartingProposalId: 1,
		Params:             &params,
		Proposals: []*govtypesv1.Proposal{
			{
				Id:             42,
				Messages:       []*types.Any{certAny, sendAny},
				Status:         govtypesv1.StatusDepositPeriod,
				SubmitTime:     &now,
				DepositEndTime: &end,
				TotalDeposit:   sdk.NewCoins(),
				Title:          "bundle",
				Summary:        "summary",
				Proposer:       authority.String(),
			},
		},
	}

	err = typesv1.ValidateGenesis(gs)
	require.Error(t, err)
	require.Contains(t, err.Error(), "proposal 42")
	require.Contains(t, err.Error(), "exactly one message")
}

// A solo cert-update proposal in genesis must still validate; the
// solo-message rule applies only to bundles.
func TestValidateGenesis_AcceptsSoloCertUpdateProposal(t *testing.T) {
	authority := sdk.AccAddress("auth________________")
	certifier := sdk.AccAddress("certifier___________")

	certMsg := certtypes.NewMsgUpdateCertifier(authority, certifier, "add", certtypes.Add)
	certAny, err := types.NewAnyWithValue(certMsg)
	require.NoError(t, err)

	now := time.Unix(100, 0)
	end := now.Add(48 * time.Hour)
	params := govtypesv1.DefaultParams()
	gs := &typesv1.GenesisState{
		StartingProposalId: 1,
		Params:             &params,
		Proposals: []*govtypesv1.Proposal{
			{
				Id:             42,
				Messages:       []*types.Any{certAny},
				Status:         govtypesv1.StatusDepositPeriod,
				SubmitTime:     &now,
				DepositEndTime: &end,
				TotalDeposit:   sdk.NewCoins(),
				Title:          "solo",
				Summary:        "summary",
				Proposer:       authority.String(),
			},
		},
	}

	require.NoError(t, typesv1.ValidateGenesis(gs))
}

package bounty_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func TestExportGenesis(t *testing.T) {
	params := types.DefaultParams()
	acc1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	acc2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	dataGS := &types.GenesisState{
		Programs: []*types.Program{
			{
				ProgramId:    "1",
				Name:         "1",
				AdminAddress: acc1.String(),
			},
			{
				ProgramId:    "2",
				Name:         "2",
				AdminAddress: acc1.String(),
			},
			{
				ProgramId:    "3",
				Name:         "3",
				AdminAddress: acc1.String(),
			},
		},
		Findings: []*types.Finding{
			{
				FindingId:        "1",
				ProgramId:        "1",
				SubmitterAddress: acc1.String(),
			},
			{
				FindingId:        "4",
				ProgramId:        "3",
				SubmitterAddress: acc1.String(),
			},
		},
		StartingTheoremId: 1,
		Theorems: []*types.Theorem{
			{
				Id:          1,
				Title:       "Test Theorem 1",
				Description: "This is a test theorem",
				Proposer:    acc1.String(),
				Status:      types.TheoremStatus_THEOREM_STATUS_PASSED,
			},
			{
				Id:          2,
				Title:       "Test Theorem 2",
				Description: "This is another test theorem",
				Proposer:    acc1.String(),
				Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
				EndTime:     &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
			},
		},
		Proofs: []*types.Proof{
			{
				Id:         "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
				TheoremId:  1,
				Prover:     acc1.String(),
				Status:     types.ProofStatus_PROOF_STATUS_PASSED,
				SubmitTime: &[]time.Time{time.Now()}[0],
			},
			{
				Id:         "60303ae22b998861bce3b28f33eec1be758a213c86c93c076dbe9f558c11c752",
				TheoremId:  1,
				Prover:     acc1.String(),
				Status:     types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
				SubmitTime: &[]time.Time{time.Now()}[0],
				EndTime:    &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
			},
		},
		Grants: []*types.Grant{
			{
				TheoremId: 1,
				Grantor:   acc1.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(1000))),
			},
		},
		Deposits: []*types.Deposit{
			{
				ProofId:   "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
				Depositor: acc1.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(500))),
			},
		},
		Rewards: []*types.Reward{
			{
				Address: acc1.String(),
				Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(100))),
			},
		},
		ImportedRewards: []*types.Reward{
			{
				Address: acc2.String(),
				Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(50))),
			},
		},
		Params: &params,
	}

	app1 := shentuapp.Setup(t, false)
	ctx1 := app1.BaseApp.NewContext(false)
	k1 := app1.BountyKeeper

	err := bounty.InitGenesis(ctx1, app1.AccountKeeper, k1, dataGS)
	require.NoError(t, err)
	exported1 := bounty.ExportGenesis(ctx1, k1)

	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	k2 := app2.BountyKeeper

	exported2 := bounty.ExportGenesis(ctx2, k2)
	require.False(t, reflect.DeepEqual(exported1, exported2))

	err = bounty.InitGenesis(ctx2, app2.AccountKeeper, k2, exported1)
	require.NoError(t, err)
	exported3 := bounty.ExportGenesis(ctx2, k2)
	require.True(t, reflect.DeepEqual(exported1, exported3))

	// Verify specific fields are exported correctly
	require.Len(t, exported1.Programs, 3)
	require.Len(t, exported1.Findings, 2)
	require.Len(t, exported1.Theorems, 2)
	require.Len(t, exported1.Proofs, 2)
	require.Len(t, exported1.Grants, 1)
	require.Len(t, exported1.Deposits, 1)
	require.Len(t, exported1.Rewards, 1)
	require.Len(t, exported1.ImportedRewards, 1)
	require.Equal(t, uint64(1), exported1.StartingTheoremId)
	require.NotNil(t, exported1.Params)
}

func TestValidateGenesis(t *testing.T) {
	acc1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	acc2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	params := types.DefaultParams()

	testCases := []struct {
		name     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			name:     "nil genesis state",
			genState: nil,
			valid:    false,
		},
		{
			name:     "valid minimal genesis state",
			genState: types.DefaultGenesisState(),
			valid:    true,
		},
		{
			name: "invalid params",
			genState: &types.GenesisState{
				StartingTheoremId: 1,
				Params:            nil,
			},
			valid: false,
		},
		{
			name: "invalid starting theorem id",
			genState: &types.GenesisState{
				StartingTheoremId: 0,
				Params:            &params,
			},
			valid: false,
		},
		{
			name: "valid genesis with all fields",
			genState: &types.GenesisState{
				Programs: []*types.Program{
					{
						ProgramId:    "1",
						Name:         "test program",
						AdminAddress: acc1.String(),
					},
				},
				Findings: []*types.Finding{
					{
						FindingId:        "1",
						ProgramId:        "1",
						SubmitterAddress: acc1.String(),
					},
				},
				StartingTheoremId: 1,
				Theorems: []*types.Theorem{
					{
						Id:          1,
						Title:       "Test Theorem",
						Description: "This is a test theorem",
						Proposer:    acc1.String(),
						Status:      types.TheoremStatus_THEOREM_STATUS_PASSED,
					},
				},
				Proofs: []*types.Proof{
					{
						Id:        "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
						TheoremId: 1,
						Prover:    acc1.String(),
						Status:    types.ProofStatus_PROOF_STATUS_PASSED,
					},
				},
				Grants: []*types.Grant{
					{
						TheoremId: 1,
						Grantor:   acc1.String(),
						Amount:    sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(1000))),
					},
				},
				Deposits: []*types.Deposit{
					{
						ProofId:   "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
						Depositor: acc1.String(),
						Amount:    sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(500))),
					},
				},
				Rewards: []*types.Reward{
					{
						Address: acc1.String(),
						Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(100))),
					},
				},
				ImportedRewards: []*types.Reward{
					{
						Address: acc2.String(),
						Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(50))),
					},
				},
				Params: &params,
			},
			valid: true,
		},
		{
			name: "invalid program - duplicate program id",
			genState: &types.GenesisState{
				Programs: []*types.Program{
					{
						ProgramId:    "1",
						Name:         "test program 1",
						AdminAddress: acc1.String(),
					},
					{
						ProgramId:    "1",
						Name:         "test program 2",
						AdminAddress: acc1.String(),
					},
				},
				StartingTheoremId: 1,
				Params:            &params,
			},
			valid: false,
		},
		{
			name: "invalid finding - program does not exist",
			genState: &types.GenesisState{
				Programs: []*types.Program{
					{
						ProgramId:    "1",
						Name:         "test program",
						AdminAddress: acc1.String(),
					},
				},
				Findings: []*types.Finding{
					{
						FindingId:        "1",
						ProgramId:        "2", // non-existent program
						SubmitterAddress: acc1.String(),
					},
				},
				StartingTheoremId: 1,
				Params:            &params,
			},
			valid: false,
		},
		{
			name: "invalid theorem - duplicate theorem id",
			genState: &types.GenesisState{
				StartingTheoremId: 1,
				Theorems: []*types.Theorem{
					{
						Id:          1,
						Title:       "Test Theorem 1",
						Description: "This is a test theorem",
						Proposer:    acc1.String(),
						Status:      types.TheoremStatus_THEOREM_STATUS_PASSED,
					},
					{
						Id:          1, // duplicate id
						Title:       "Test Theorem 2",
						Description: "This is another test theorem",
						Proposer:    acc1.String(),
						Status:      types.TheoremStatus_THEOREM_STATUS_PASSED,
					},
				},
				Params: &params,
			},
			valid: false,
		},
		{
			name: "invalid proof - theorem does not exist",
			genState: &types.GenesisState{
				StartingTheoremId: 1,
				Theorems: []*types.Theorem{
					{
						Id:          1,
						Title:       "Test Theorem",
						Description: "This is a test theorem",
						Proposer:    acc1.String(),
						Status:      types.TheoremStatus_THEOREM_STATUS_PASSED,
					},
				},
				Proofs: []*types.Proof{
					{
						Id:        "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
						TheoremId: 2, // non-existent theorem
						Prover:    acc1.String(),
						Status:    types.ProofStatus_PROOF_STATUS_PASSED,
					},
				},
				Params: &params,
			},
			valid: false,
		},
		{
			name: "invalid grant - theorem does not exist",
			genState: &types.GenesisState{
				StartingTheoremId: 1,
				Theorems: []*types.Theorem{
					{
						Id:          1,
						Title:       "Test Theorem",
						Description: "This is a test theorem",
						Proposer:    acc1.String(),
						Status:      types.TheoremStatus_THEOREM_STATUS_PASSED,
					},
				},
				Grants: []*types.Grant{
					{
						TheoremId: 2, // non-existent theorem
						Grantor:   acc1.String(),
						Amount:    sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(1000))),
					},
				},
				Params: &params,
			},
			valid: false,
		},
		{
			name: "invalid deposit - proof does not exist",
			genState: &types.GenesisState{
				StartingTheoremId: 1,
				Theorems: []*types.Theorem{
					{
						Id:          1,
						Title:       "Test Theorem",
						Description: "This is a test theorem",
						Proposer:    acc1.String(),
						Status:      types.TheoremStatus_THEOREM_STATUS_PASSED,
					},
				},
				Proofs: []*types.Proof{
					{
						Id:        "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
						TheoremId: 1,
						Prover:    acc1.String(),
						Status:    types.ProofStatus_PROOF_STATUS_PASSED,
					},
				},
				Deposits: []*types.Deposit{
					{
						ProofId:   "60303ae22b998861bce3b28f33eec1be758a213c86c93c076dbe9f558c11c752", // non-existent proof
						Depositor: acc1.String(),
						Amount:    sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(500))),
					},
				},
				Params: &params,
			},
			valid: false,
		},
		{
			name: "invalid reward - empty address",
			genState: &types.GenesisState{
				StartingTheoremId: 1,
				Rewards: []*types.Reward{
					{
						Address: "", // empty address
						Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(100))),
					},
				},
				Params: &params,
			},
			valid: false,
		},
		{
			name: "invalid imported reward - empty address",
			genState: &types.GenesisState{
				StartingTheoremId: 1,
				ImportedRewards: []*types.Reward{
					{
						Address: "", // empty address
						Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(50))),
					},
				},
				Params: &params,
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateGenesis(tc.genState)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestImportedRewardsGenesis(t *testing.T) {
	acc1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	acc2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	acc3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	params := types.DefaultParams()

	// Test genesis state with imported rewards
	dataGS := &types.GenesisState{
		StartingTheoremId: 1,
		ImportedRewards: []*types.Reward{
			{
				Address: acc1.String(),
				Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(100))),
			},
			{
				Address: acc2.String(),
				Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(200))),
			},
		},
		Params: &params,
	}

	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	k := app.BountyKeeper

	// Test InitGenesis with imported rewards
	err := bounty.InitGenesis(ctx, app.AccountKeeper, k, dataGS)
	require.NoError(t, err)

	// Test ExportGenesis with imported rewards
	exported := bounty.ExportGenesis(ctx, k)
	require.Len(t, exported.ImportedRewards, 2)

	// Verify imported rewards are correctly stored and exported
	importedReward1, err := k.ImportedRewards.Get(ctx, acc1)
	require.NoError(t, err)
	require.Equal(t, acc1.String(), importedReward1.Address)
	require.Equal(t, sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(100))), importedReward1.Reward)

	importedReward2, err := k.ImportedRewards.Get(ctx, acc2)
	require.NoError(t, err)
	require.Equal(t, acc2.String(), importedReward2.Address)
	require.Equal(t, sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(200))), importedReward2.Reward)

	// Test that non-existent imported reward returns error
	_, err = k.ImportedRewards.Get(ctx, acc3)
	require.Error(t, err)
}

func TestGenesisRoundTrip(t *testing.T) {
	acc1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	acc2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	params := types.DefaultParams()

	// Create a comprehensive genesis state
	originalGS := &types.GenesisState{
		Programs: []*types.Program{
			{
				ProgramId:    "program1",
				Name:         "Test Program",
				AdminAddress: acc1.String(),
			},
		},
		Findings: []*types.Finding{
			{
				FindingId:        "finding1",
				ProgramId:        "program1",
				SubmitterAddress: acc1.String(),
			},
		},
		StartingTheoremId: 1,
		Theorems: []*types.Theorem{
			{
				Id:          1,
				Title:       "Test Theorem",
				Description: "This is a test theorem",
				Proposer:    acc1.String(),
				Status:      types.TheoremStatus_THEOREM_STATUS_PASSED,
			},
		},
		Proofs: []*types.Proof{
			{
				Id:        "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
				TheoremId: 1,
				Prover:    acc1.String(),
				Status:    types.ProofStatus_PROOF_STATUS_PASSED,
			},
		},
		Grants: []*types.Grant{
			{
				TheoremId: 1,
				Grantor:   acc1.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(1000))),
			},
		},
		Deposits: []*types.Deposit{
			{
				ProofId:   "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
				Depositor: acc1.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin("uctk", sdkmath.NewInt(500))),
			},
		},
		Rewards: []*types.Reward{
			{
				Address: acc1.String(),
				Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(100))),
			},
		},
		ImportedRewards: []*types.Reward{
			{
				Address: acc2.String(),
				Reward:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("uctk", sdkmath.LegacyNewDec(50))),
			},
		},
		Params: &params,
	}

	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	k := app.BountyKeeper

	// Test InitGenesis -> ExportGenesis -> InitGenesis round trip
	err := bounty.InitGenesis(ctx, app.AccountKeeper, k, originalGS)
	require.NoError(t, err)

	exported1 := bounty.ExportGenesis(ctx, k)
	require.True(t, reflect.DeepEqual(originalGS, exported1))

	// Test another round trip
	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	k2 := app2.BountyKeeper

	err = bounty.InitGenesis(ctx2, app2.AccountKeeper, k2, exported1)
	require.NoError(t, err)

	exported2 := bounty.ExportGenesis(ctx2, k2)
	require.True(t, reflect.DeepEqual(exported1, exported2))
}

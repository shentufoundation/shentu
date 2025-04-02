package bounty_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func TestExportGenesis(t *testing.T) {
	dataGS := &types.GenesisState{
		Programs: []*types.Program{
			{
				ProgramId: "1",
			},
			{
				ProgramId: "2",
			},
			{
				ProgramId: "3",
			},
		},
		Findings: []*types.Finding{
			{
				FindingId: "1",
				ProgramId: "1",
			},
			{
				FindingId: "4",
				ProgramId: "3",
			},
		},
	}

	app1 := shentuapp.Setup(t, false)
	ctx1 := app1.BaseApp.NewContext(false)
	k1 := app1.BountyKeeper

	bounty.InitGenesis(ctx1, app1.AccountKeeper, k1, dataGS)
	exported1 := bounty.ExportGenesis(ctx1, k1)

	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	k2 := app2.BountyKeeper

	exported2 := bounty.ExportGenesis(ctx2, k2)
	require.False(t, reflect.DeepEqual(exported1, exported2))

	bounty.InitGenesis(ctx2, app2.AccountKeeper, k2, exported1)
	exported3 := bounty.ExportGenesis(ctx2, k2)
	require.True(t, reflect.DeepEqual(exported1, exported3))
}

func TestValidateGenesis(t *testing.T) {
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

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func ValidateGenesis(data *GenesisState) error {
func TestValidateGenesis(t *testing.T) {
	type args struct {
		dataGS GenesisState
	}

	type errArgs struct {
		shouldPass bool
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"genesis(1)  -> success",
			args{
				dataGS: *DefaultGenesisState(),
			},
			errArgs{
				shouldPass: true,
			},
		},
		{
			"genesis(2)  -> fail",
			args{
				dataGS: GenesisState{
					Programs: []*Program{
						{
							ProgramId: "100",
						},
					},
					Findings: []*Finding{},
				},
			},
			errArgs{
				shouldPass: false,
			},
		},

		{"genesis(3)  -> findingId error",
			args{
				dataGS: GenesisState{
					Programs: []*Program{
						{
							ProgramId: "1",
						},
					},
					Findings: []*Finding{
						{
							FindingId: "100",
						},
					},
				},
			},
			errArgs{
				shouldPass: false,
			},
		},
		{"genesis(4)  -> invalid programId error",
			args{
				dataGS: GenesisState{
					Programs: []*Program{
						{
							ProgramId: "1",
						},
					},
					Findings: []*Finding{
						{
							ProgramId: "10",
							FindingId: "1",
						},
					},
				},
			},
			errArgs{
				shouldPass: false,
			},
		},
	}

	for _, tc := range tests {
		err := ValidateGenesis(&tc.args.dataGS)
		if tc.errArgs.shouldPass {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

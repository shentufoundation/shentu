package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/shentufoundation/shentu/v2/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	minGrant := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdkmath.NewInt(1000000)))
	minDeposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdkmath.NewInt(1000000)))
	maxProofPeriod := time.Minute * 5
	proofHashLockPeriod := time.Minute * 5
	proofDetailLockPeriod := time.Minute * 5
	checkerResponsePeriod := time.Minute * 5

	return &GenesisState{
		Programs:          []Program{},
		Findings:          []Finding{},
		StartingTheoremId: 1,
		Params: &Params{
			MinGrant:              minGrant,
			MinDeposit:            minDeposit,
			TheoremMaxProofPeriod: &maxProofPeriod,
			ProofHashLockPeriod:   &proofHashLockPeriod,
			ProofDetailLockPeriod: &proofDetailLockPeriod,
			CheckerResponsePeriod: &checkerResponsePeriod,
		},
	}
}

// ValidateGenesis - validate bounty genesis data
func ValidateGenesis(data *GenesisState) error {
	programs := make(map[string]int)
	for i, program := range data.Programs {
		programIndex, ok := programs[program.ProgramId]
		if ok {
			//repeat programId
			return errorsmod.Wrapf(ErrProgramID, "already program[%s], this program[%s]",
				data.Programs[programIndex].String(), program.String())
		}

		if err := program.ValidateBasic(); err != nil {
			return err
		}
		programs[program.ProgramId] = i
	}

	for _, finding := range data.Findings {
		//Check if it is a valid programID
		_, ok := programs[finding.ProgramId]
		if !ok {
			return ErrProgramID
		}

		if err := finding.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

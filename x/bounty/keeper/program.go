package keeper

import (
	"crypto/sha256"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) ActivateProgram(ctx sdk.Context, pid string, caller sdk.AccAddress) error {
	program, err := k.Programs.Get(ctx, pid)
	if err != nil {
		return err
	}
	// Check if the program is already active
	if program.Status == types.ProgramStatusActive {
		return types.ErrProgramAlreadyActive
	}

	// Check the permissions. Only the bounty cert address can operate.
	if !k.certKeeper.IsBountyAdmin(ctx, caller) {
		return types.ErrProgramOperatorNotAllowed
	}

	program.Status = types.ProgramStatusActive
	return k.Programs.Set(ctx, program.ProgramId, program)
}

func (k Keeper) CloseProgram(ctx sdk.Context, pid string, caller sdk.AccAddress) error {
	program, err := k.Programs.Get(ctx, pid)
	if err != nil {
		return err
	}

	// Check if the program is already closed
	if program.Status == types.ProgramStatusClosed {
		return types.ErrProgramAlreadyClosed
	}

	// The program cannot be closed
	// There are 3 finding states: FindingStatusSubmitted FindingStatusActive FindingStatusConfirmed
	fidsList, err := k.GetProgramFindings(ctx, pid)
	if err != nil {
		return err
	}
	for _, fid := range fidsList {
		finding, err := k.Findings.Get(ctx, fid)
		if err != nil {
			return err
		}
		if finding.Status == types.FindingStatusSubmitted ||
			finding.Status == types.FindingStatusActive ||
			finding.Status == types.FindingStatusConfirmed {
			return types.ErrProgramCloseNotAllowed
		}
	}

	// Check the permissions. Only the admin of the program or bounty cert address can operate.
	if program.AdminAddress != caller.String() && !k.certKeeper.IsBountyAdmin(ctx, caller) {
		return types.ErrProgramOperatorNotAllowed
	}

	// Close the program and update its status in the store
	program.Status = types.ProgramStatusClosed
	return k.Programs.Set(ctx, program.ProgramId, program)
}

func (k Keeper) GetProgramFingerprintHash(program *types.Program) string {
	programFingerprint := &types.ProgramFingerprint{
		ProgramId:    program.ProgramId,
		Name:         program.Name,
		Detail:       program.Detail,
		AdminAddress: program.AdminAddress,
		Status:       program.Status,
	}

	bz := k.cdc.MustMarshal(programFingerprint)
	hash := sha256.Sum256(bz)
	return hex.EncodeToString(hash[:])
}

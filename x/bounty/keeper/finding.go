package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) ConfirmFinding(ctx sdk.Context, msg *types.MsgConfirmFinding) (types.Finding, error) {
	finding, err := k.Findings.Get(ctx, msg.FindingId)
	if err != nil {
		return finding, err
	}
	// only StatusActive can be confirmed
	if finding.Status != types.FindingStatusActive {
		return finding, types.ErrFindingStatusInvalid
	}

	// get program
	program, err := k.Programs.Get(ctx, finding.ProgramId)
	if err != nil {
		return finding, err
	}
	if program.Status != types.ProgramStatusActive {
		return finding, types.ErrProgramNotActive
	}

	// only program admin can confirm finding
	if program.AdminAddress != msg.OperatorAddress {
		return finding, types.ErrProgramOperatorNotAllowed
	}

	// fingerprint comparison
	fingerprintHash := k.GetFindingFingerprintHash(&finding)
	if msg.Fingerprint != fingerprintHash {
		return finding, types.ErrFindingHashInvalid
	}
	return finding, nil
}

func (k Keeper) GetFindingFingerprintHash(finding *types.Finding) string {
	findingFingerprint := &types.FindingFingerprint{
		ProgramId:   finding.ProgramId,
		FindingId:   finding.FindingId,
		FindingHash: finding.FindingHash,
		Status:      finding.Status,
		PaymentHash: finding.PaymentHash,
	}

	bz := k.cdc.MustMarshal(findingFingerprint)
	hash := sha256.Sum256(bz)
	return hex.EncodeToString(hash[:])
}

func (k Keeper) GetProgramFindings(ctx context.Context, programID string) ([]string, error) {
	var findingIDs []string

	rng := collections.NewPrefixedPairRange[string, string](programID)
	err := k.ProgramFindings.Walk(ctx, rng, func(key collections.Pair[string, string]) (stop bool, err error) {
		if key.K1() == programID {
			findingIDs = append(findingIDs, key.K2())
		}
		return false, nil
	})
	if err != nil {
		return findingIDs, err
	}

	return findingIDs, nil
}

//func (k Keeper) AppendFidToFidList(ctx sdk.Context, pid, fid string) error {
//	fids, err := k.GetPidFindingIDList(ctx, pid)
//	if err != nil {
//		return err
//	}
//
//	fids = append(fids, fid)
//	err = k.SetPidFindingIDList(ctx, pid, fids)
//	return err
//}

//func StringsToBytes(list []string) ([]byte, error) {
//	marshal, err := json.Marshal(list)
//	if err != nil {
//		return nil, err
//	}
//	return marshal, nil
//}
//
//func BytesToStrings(list []byte) ([]string, error) {
//	var fids []string
//	err := json.Unmarshal(list, &fids)
//	if err != nil {
//		return nil, err
//	}
//
//	return fids, nil
//}

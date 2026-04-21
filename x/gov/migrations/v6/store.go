// Package v6 migrates the gov module off the two-round security
// proposal model. Under that model, a CertifierUpdate proposal that
// passed the certifier round was marked cert_voted=true (CertVotesKey
// prefix) and advanced to a validator stake round. The new model
// removes the stake round entirely: cert-update proposals are decided
// by the cert round alone, so cert_voted=true entries are never
// written again.
//
// MigrateStore cleans up orphaned cert_voted keys and panics if any
// v6-era state would be misinterpreted by v7 — namely an in-flight
// cert_voted=true proposal (stake votes would re-tally as head-counts),
// a bundled cert-update proposal (the non-cert messages would pass
// cert-only under v7), or a legacy weighted certifier ballot (the old
// tally counted it as 1 head; the new tally drops it).
package v6

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// CertVotesKeyPrefix is the KV store prefix used by the old two-round
// security-proposal model to record cert-update proposals that had
// cleared the certifier round. It is not written anywhere in the
// current binary; this migration is the only remaining reader and
// deletes every entry under this prefix. Kept package-local so no
// other code can re-adopt it.
var CertVotesKeyPrefix = []byte("certvote")

// IsCertUpdateMsgFn reports whether a proposal message is a
// cert-update (either the Msg or the legacy content proposal). Passed
// in from the keeper layer so the migration doesn't import keeper.
type IsCertUpdateMsgFn func(sdk.Msg) (bool, error)

// MigrateStore performs the v6→v7 in-place migration for the gov
// module. It scans the legacy cert_voted entries plus the live
// Proposals/Votes collections for state that v7 would misinterpret,
// panics with a human-readable blocker list if any is found, and
// otherwise deletes every cert_voted entry.
//
// The caller passes in the gov keeper's Proposals and Votes
// collections plus a cert-update-message matcher so this package
// doesn't need to import x/gov/keeper.
func MigrateStore(
	ctx sdk.Context,
	storeService corestoretypes.KVStoreService,
	proposals collections.Map[uint64, govtypesv1.Proposal],
	votes collections.Map[collections.Pair[uint64, sdk.AccAddress], govtypesv1.Vote],
	isCertUpdateMsg IsCertUpdateMsgFn,
) error {
	kv := storeService.OpenKVStore(ctx)
	store := runtime.KVStoreAdapter(kv)
	certVotes := prefix.NewStore(store, CertVotesKeyPrefix)

	var blockers []string
	keys := make([][]byte, 0)

	it := certVotes.Iterator(nil, nil)
	for ; it.Valid(); it.Next() {
		v := it.Value()
		if len(v) != 8 {
			// Defensive: entries have always been 8 bytes (BigEndian
			// uint64) but we don't want to panic on an unexpected
			// value. Skip and delete.
			continue
		}
		proposalID := binary.BigEndian.Uint64(v)
		active, err := isActiveCertUpdate(ctx, proposals, isCertUpdateMsg, proposalID)
		if err != nil {
			_ = it.Close()
			return err
		}
		if active {
			blockers = append(blockers, fmt.Sprintf(
				"proposal %d: cert_voted=true still set on an active CertifierUpdate proposal (two-round legacy)",
				proposalID,
			))
		}
		keys = append(keys, append([]byte(nil), it.Key()...))
	}
	if err := it.Close(); err != nil {
		return err
	}

	extra, err := scanLegacyCertState(ctx, proposals, votes, isCertUpdateMsg)
	if err != nil {
		return err
	}
	blockers = append(blockers, extra...)

	if len(blockers) > 0 {
		panic(fmt.Sprintf(
			"gov v6→v7 migration blocked by legacy two-round state:\n  - %s\n"+
				"Let those proposals complete their voting period on the previous binary, then retry the upgrade.",
			strings.Join(blockers, "\n  - "),
		))
	}

	for _, k := range keys {
		certVotes.Delete(k)
	}
	return nil
}

// isActiveCertUpdate reports whether the proposal is a CertifierUpdate
// proposal still in its voting period — the condition under which a
// stale cert_voted=true entry would cause v7 to re-tally stake votes
// as head-counts.
func isActiveCertUpdate(
	ctx sdk.Context,
	proposals collections.Map[uint64, govtypesv1.Proposal],
	isCertUpdateMsg IsCertUpdateMsgFn,
	proposalID uint64,
) (bool, error) {
	proposal, err := proposals.Get(ctx, proposalID)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	if proposal.Status != govtypesv1.StatusVotingPeriod {
		return false, nil
	}
	msgs, err := proposal.GetMsgs()
	if err != nil {
		return false, err
	}
	for _, msg := range msgs {
		isCertUpdate, err := isCertUpdateMsg(msg)
		if err != nil {
			return false, err
		}
		if isCertUpdate {
			return true, nil
		}
	}
	return false, nil
}

// scanLegacyCertState walks the Proposals collection and reports any
// v6-era state that v7 would silently misinterpret:
//
//   - Bundled cert-update proposals — a proposal whose message list
//     contains a MsgUpdateCertifier alongside other messages. v6
//     allowed these and routed them through both the certifier round
//     and the validator stake round. v7's cert-only path would decide
//     the whole bundle from the certifier round, letting the non-cert
//     messages execute without stake approval.
//
//   - Legacy weighted/multi-option certifier ballots — v6 counted any
//     stored MsgVoteWeighted with a single option as one head
//     regardless of its weight string. v7's SecurityTally requires
//     weight=1 and len(options)=1, so these ballots would be dropped
//     and could silently flip an outcome.
//
// Both DepositPeriod and VotingPeriod are scanned for bundles: a
// deposit-period bundle that reaches min deposit post-upgrade would be
// admitted by ActivateVotingPeriod and then silently rerouted to plain
// stake voting because CertifierVoteIsRequired returns false for
// bundles. Vote invariants are only checked in VotingPeriod since no
// ballots can exist before.
func scanLegacyCertState(
	ctx sdk.Context,
	proposals collections.Map[uint64, govtypesv1.Proposal],
	votes collections.Map[collections.Pair[uint64, sdk.AccAddress], govtypesv1.Vote],
	isCertUpdateMsg IsCertUpdateMsgFn,
) ([]string, error) {
	var blockers []string
	one := math.LegacyOneDec()

	err := proposals.Walk(ctx, nil, func(id uint64, proposal govtypesv1.Proposal) (bool, error) {
		status := proposal.Status
		if status != govtypesv1.StatusVotingPeriod && status != govtypesv1.StatusDepositPeriod {
			return false, nil
		}
		msgs, err := proposal.GetMsgs()
		if err != nil {
			return false, err
		}
		hasCert := false
		for _, msg := range msgs {
			isCert, err := isCertUpdateMsg(msg)
			if err != nil {
				return false, err
			}
			if isCert {
				hasCert = true
				break
			}
		}
		if !hasCert {
			return false, nil
		}
		if len(msgs) > 1 {
			phase := "voting"
			if status == govtypesv1.StatusDepositPeriod {
				phase = "deposit"
			}
			blockers = append(blockers, fmt.Sprintf(
				"proposal %d (%s period): bundles cert-update with %d other messages (v7 forbids bundled cert-update proposals)",
				id, phase, len(msgs)-1,
			))
			return false, nil
		}

		if status != govtypesv1.StatusVotingPeriod {
			return false, nil
		}
		rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](id)
		walkErr := votes.Walk(ctx, rng, func(_ collections.Pair[uint64, sdk.AccAddress], vote govtypesv1.Vote) (bool, error) {
			if len(vote.Options) != 1 {
				blockers = append(blockers, fmt.Sprintf(
					"proposal %d: certifier %s cast a %d-option ballot (v7 requires single-option)",
					id, vote.Voter, len(vote.Options),
				))
				return false, nil
			}
			weight, err := math.LegacyNewDecFromStr(vote.Options[0].Weight)
			if err != nil || !weight.Equal(one) {
				blockers = append(blockers, fmt.Sprintf(
					"proposal %d: certifier %s cast a weighted ballot (weight=%q; v7 requires weight=1)",
					id, vote.Voter, vote.Options[0].Weight,
				))
			}
			return false, nil
		})
		return false, walkErr
	})
	return blockers, err
}

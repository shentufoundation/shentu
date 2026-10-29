// Package ante provides the x/blacklist ante decorator that rejects
// transactions signed (or authorized) by blacklisted addresses.
package ante

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/shentufoundation/shentu/v2/x/blacklist/keeper"
	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

// BlacklistDecorator rejects transactions that are signed by, fee-granted by,
// or authz-executed on behalf of a blacklisted address.
type BlacklistDecorator struct {
	k keeper.Keeper
}

// NewBlacklistDecorator constructs a new BlacklistDecorator.
func NewBlacklistDecorator(k keeper.Keeper) BlacklistDecorator {
	return BlacklistDecorator{k: k}
}

// AnteHandle implements sdk.AnteDecorator. Rejected transactions fail before
// fee deduction and signature verification.
func (d BlacklistDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// 1. Every transaction signer.
	if sigTx, ok := tx.(authsigning.SigVerifiableTx); ok {
		signers, err := sigTx.GetSigners()
		if err != nil {
			return ctx, errorsmod.Wrap(err, "failed to get signers")
		}
		for _, signer := range signers {
			if err := d.checkAddress(ctx, signer, "signer"); err != nil {
				return ctx, err
			}
		}
	}

	// 2. The fee granter: prevents a blacklisted address's balance from being
	// drained as gas payment for transactions signed by other accounts.
	if feeTx, ok := tx.(sdk.FeeTx); ok && len(feeTx.FeeGranter()) > 0 {
		if err := d.checkAddress(ctx, sdk.AccAddress(feeTx.FeeGranter()), "fee granter"); err != nil {
			return ctx, err
		}
	}

	// 3. Authz execution: a tx signed by a clean grantee can carry inner
	// messages whose actual signer is a blacklisted granter. Recursively
	// check every (possibly nested) inner message's signers.
	if err := d.checkMsgs(ctx, tx.GetMsgs()); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

// checkAddress rejects the transaction if the address is currently
// blacklisted.
func (d BlacklistDecorator) checkAddress(ctx sdk.Context, addr sdk.AccAddress, role string) error {
	entry, blacklisted, err := d.k.IsBlacklisted(ctx, addr.String())
	if err != nil {
		return err
	}
	if !blacklisted {
		return nil
	}
	return errorsmod.Wrapf(
		types.ErrAddressBlacklisted,
		"%s %s is blacklisted (category: %s)",
		role, addr.String(), entry.Category.String(),
	)
}

// checkMsgs walks the message tree, checking the signers of every message and
// recursively descending into authz MsgExec wrappers.
func (d BlacklistDecorator) checkMsgs(ctx sdk.Context, msgs []sdk.Msg) error {
	pending := make([]sdk.Msg, 0, len(msgs))
	pending = append(pending, msgs...)

	for len(pending) > 0 {
		msg := pending[0]
		pending = pending[1:]

		// For authz MsgExec this returns the grantee; inner messages are
		// unpacked below so their real signers are checked as well.
		signers, _, err := d.k.Cdc().GetMsgV1Signers(msg)
		if err != nil {
			return errorsmod.Wrap(err, "failed to get message signers")
		}
		for _, signer := range signers {
			if err := d.checkAddress(ctx, sdk.AccAddress(signer), "message signer"); err != nil {
				return err
			}
		}

		if exec, ok := msg.(*authz.MsgExec); ok {
			inner, err := d.unpackExec(exec)
			if err != nil {
				return err
			}
			pending = append(pending, inner...)
		}
	}

	return nil
}

// unpackExec decodes the inner messages of an authz MsgExec.
func (d BlacklistDecorator) unpackExec(exec *authz.MsgExec) ([]sdk.Msg, error) {
	msgs := make([]sdk.Msg, 0, len(exec.Msgs))
	for _, anyMsg := range exec.Msgs {
		var inner sdk.Msg
		if err := d.k.Cdc().UnpackAny(anyMsg, &inner); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to unpack authz inner message")
		}
		msgs = append(msgs, inner)
	}
	return msgs, nil
}

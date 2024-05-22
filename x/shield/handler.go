package shield

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/shentufoundation/shentu/v2/x/shield/keeper"
	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

func NewShieldClaimProposalHandler(k keeper.Keeper) govtypesv1beta1.Handler {
	return func(ctx sdk.Context, content govtypesv1beta1.Content) error {
		switch c := content.(type) {
		case *types.ShieldClaimProposal:
			return handleShieldClaimProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized shield proposal content type: %T", c)
		}
	}
}

func handleShieldClaimProposal(ctx sdk.Context, k keeper.Keeper, p *types.ShieldClaimProposal) error {
	proposerAddr, err := sdk.AccAddressFromBech32(p.Proposer)
	if err != nil {
		panic(err)
	}
	if err := k.CreateReimbursement(ctx, p.ProposalId, p.Loss, proposerAddr); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateReimbursement,
			sdk.NewAttribute(types.AttributeKeyPurchaseID, strconv.FormatUint(p.PurchaseId, 10)),
			sdk.NewAttribute(types.AttributeKeyCompensationAmount, p.Loss.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, p.Proposer),
		),
	})
	return nil
}

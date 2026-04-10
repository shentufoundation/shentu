package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

const (
	OpWeightMsgUpdateCertifier      = "op_weight_msg_update_certifier"
	DefaultWeightMsgUpdateCertifier = 5
)

// ProposalMsgs defines the module weighted proposals' contents.
func ProposalMsgs(k keeper.Keeper) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OpWeightMsgUpdateCertifier,
			DefaultWeightMsgUpdateCertifier,
			SimulateMsgUpdateCertifier(k),
		),
	}
}

// SimulateMsgUpdateCertifier returns a function that generates a random MsgUpdateCertifier.
func SimulateMsgUpdateCertifier(k keeper.Keeper) simtypes.MsgSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		authority := authtypes.NewModuleAddress(govtypes.ModuleName)
		certifiers := k.GetAllCertifiers(ctx)

		var operation types.AddOrRemove
		var certifier sdk.AccAddress
		var proposer sdk.AccAddress

		// Pick a random proposer from existing certifiers.
		if len(certifiers) > 0 {
			idx := r.Intn(len(certifiers))
			proposerAddr, err := sdk.AccAddressFromBech32(certifiers[idx].Address)
			if err == nil {
				proposer = proposerAddr
			}
		}
		if proposer == nil {
			proposer = accs[r.Intn(len(accs))].Address
		}

		switch r.Intn(2) {
		case 0: // add a new certifier
			operation = types.Add
			// Find an account that is not already a certifier.
			for _, acc := range accs {
				isCertifier, err := k.IsCertifier(ctx, acc.Address)
				if err == nil && !isCertifier {
					certifier = acc.Address
					break
				}
			}
			if certifier == nil {
				// All accounts are certifiers; pick a random one anyway.
				certifier = accs[r.Intn(len(accs))].Address
			}

		case 1: // remove an existing certifier
			if len(certifiers) <= 1 {
				// Cannot remove the only certifier; fall back to add.
				operation = types.Add
				certifier = accs[r.Intn(len(accs))].Address
			} else {
				operation = types.Remove
				idx := r.Intn(len(certifiers))
				addr, err := sdk.AccAddressFromBech32(certifiers[idx].Address)
				if err != nil {
					// Fallback to add.
					operation = types.Add
					certifier = accs[r.Intn(len(accs))].Address
				} else {
					certifier = addr
				}
			}
		}

		description := simtypes.RandStringOfLength(r, 10)
		return types.NewMsgUpdateCertifier(authority, certifier, description, operation, proposer)
	}
}

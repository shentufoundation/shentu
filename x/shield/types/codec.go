package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the account types and interface.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreatePool{}, "shield/MsgCreatePool", nil)
	cdc.RegisterConcrete(MsgUpdatePool{}, "shield/MsgUpdatePool", nil)
	cdc.RegisterConcrete(MsgPausePool{}, "shield/MsgPausePool", nil)
	cdc.RegisterConcrete(MsgResumePool{}, "shield/MsgResumePool", nil)
	cdc.RegisterConcrete(MsgDepositCollateral{}, "shield/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawCollateral{}, "shield/MsgWithdrawCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawRewards{}, "shield/MsgWithdrawRewards", nil)
	cdc.RegisterConcrete(MsgWithdrawForeignRewards{}, "shield/MsgWithdrawForeignRewards", nil)
	cdc.RegisterConcrete(MsgClearPayouts{}, "shield/MsgClearPayouts", nil)
	cdc.RegisterConcrete(ShieldClaimProposal{}, "shield/ShieldClaimProposal", nil)
	cdc.RegisterConcrete(MsgPurchaseShield{}, "shield/MsgPurchaseShield", nil)
	cdc.RegisterConcrete(MsgWithdrawReimbursement{}, "shield/MsgWithdrawReimbursement", nil)
}

// ModuleCdc is the generic sealed codec to be used throughout module.
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}

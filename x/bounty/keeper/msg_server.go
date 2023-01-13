package keeper

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto/ecies"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateProgram(goCtx context.Context, msg *types.MsgCreateProgram) (*types.MsgCreateProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creatorAddr, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, msg.Deposit)
	if err != nil {
		return nil, err
	}

	nextID := k.GetNextProgramID(ctx)

	program := types.Program{
		ProgramId:         nextID,
		CreatorAddress:    msg.CreatorAddress,
		SubmissionEndTime: msg.SubmissionEndTime,
		Description:       msg.Description,
		EncryptionKey:     msg.EncryptionKey,
		Deposit:           msg.Deposit,
		CommissionRate:    msg.CommissionRate,
	}

	k.SetProgram(ctx, program)

	k.SetNextProgramID(ctx, nextID+1)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateProgram,
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(program.ProgramId, 10)),
			sdk.NewAttribute(types.AttributeKeyDeposit, sdk.NewCoins(msg.Deposit...).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CreatorAddress),
		),
	})

	return &types.MsgCreateProgramResponse{ProgramId: nextID}, nil
}

func (k msgServer) SubmitFinding(goCtx context.Context, msg *types.MsgSubmitFinding) (*types.MsgSubmitFindingResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	program, isExist := k.GetProgram(ctx, msg.ProgramId)
	if !isExist {
		return nil, fmt.Errorf("no program id:%d", msg.ProgramId)
	}

	if !program.Active {
		return nil, fmt.Errorf("program id:%d is closed", msg.ProgramId)
	}

	var eciesEncKey ecies.PublicKey
	err = k.cdc.UnpackAny(program.EncryptionKey, &eciesEncKey)
	if err != nil {
		return nil, fmt.Errorf("EncryptionKey error")
	}

	encryptedDesc, err := ecies.Encrypt(rand.Reader, &eciesEncKey, []byte(msg.Desc), nil, nil)
	if err != nil {
		return nil, err
	}

	encryptedPoc, err := ecies.Encrypt(rand.Reader, &eciesEncKey, []byte(msg.Poc), nil, nil)
	if err != nil {
		return nil, err
	}

	findingID := k.GetNextFindingID(ctx)

	var descAny *codectypes.Any
	var pocAny *codectypes.Any

	encDesc := types.EciesEncryptedDesc{
		Desc: encryptedDesc,
	}
	if descAny, err = codectypes.NewAnyWithValue(&encDesc); err != nil {
		return nil, err
	}

	encPoc := types.EciesEncryptedPoc{
		Poc: encryptedPoc,
	}
	if pocAny, err = codectypes.NewAnyWithValue(&encPoc); err != nil {
		return nil, err
	}

	finding := types.Finding{
		FindingId:        findingID,
		Title:            msg.Title,
		EncryptedDesc:    descAny,
		ProgramId:        msg.ProgramId,
		SeverityLevel:    msg.SeverityLevel,
		EncryptedPoc:     pocAny,
		SubmitterAddress: msg.SubmitterAddress,
		FindingStatus:    types.FindingStatusUnConfirmed,
	}

	err = k.AppendFidToFidList(ctx, msg.ProgramId, findingID)
	if err != nil {
		return nil, err
	}

	k.SetFinding(ctx, finding)
	k.SetNextFindingID(ctx, findingID+1)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSubmitFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, strconv.FormatUint(finding.FindingId, 10)),
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(finding.ProgramId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SubmitterAddress),
		),
	})

	return &types.MsgSubmitFindingResponse{
		FindingId: finding.FindingId,
	}, nil
}

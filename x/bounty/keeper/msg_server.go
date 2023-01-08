package keeper

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto/ecies"

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

	nextID := k.GetNextProgramID(ctx)

	program := types.Program{
		ProgramId:         nextID,
		CreatorAddress:    msg.CreatorAddress,
		SubmissionEndTime: msg.SubmissionEndTime,
		//	JudgingEndTime:    msg.JudgingEndTime,
		//	ClaimEndTime:      msg.ClaimEndTime,
		Description:    msg.Description,
		EncryptionKey:  msg.EncryptionKey,
		Deposit:        msg.Deposit,
		CommissionRate: msg.CommissionRate,
	}

	k.SetProgram(ctx, program)

	// increment before storing
	nextID++
	k.SetNextProgramID(ctx, nextID)

	creatorAddr, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, msg.Deposit)
	if err != nil {
		return nil, err
	}

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

	return &types.MsgCreateProgramResponse{}, nil
}

func (k msgServer) SubmitFinding(goCtx context.Context, msg *types.MsgSubmitFinding) (*types.MsgSubmitFindingResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	program, isExist := k.GetProgram(ctx, msg.Pid)
	if !isExist {
		return nil, fmt.Errorf("no program id:%d", msg.Pid)
	}

	encryptedDesc, err := ecies.Encrypt(rand.Reader, ecies.PublicKey(program.EncryptionKey), []byte(msg.Desc), nil, nil)
	if err != nil {
		return nil, err
	}

	encryptedPoc, err := ecies.Encrypt(rand.Reader, ecies.PublicKey(program.EncryptionKey), []byte(msg.Poc), nil, nil)
	if err != nil {
		return nil, err
	}

	nextID := k.GetNextFindingID(ctx)

	finding := types.Finding{
		FindingId:        nextID,
		Title:            msg.Title,
		EncryptedDesc:    encryptedDesc,
		Pid:              msg.Pid,
		SeverityLevel:    msg.SeverityLevel,
		EncryptedPoc:     encryptedPoc,
		SubmitterAddress: msg.SubmitterAddress,
	}

	k.SetFinding(ctx, finding)
	k.SetNextFindingID(ctx, nextID+1)

	err = k.AppendFidToFidList(ctx, msg.Pid, nextID)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSubmitFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, strconv.FormatUint(finding.FindingId, 10)),
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(finding.Pid, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SubmitterAddress),
		),
	})

	return &types.MsgSubmitFindingResponse{
		Fid: finding.FindingId,
	}, nil
}

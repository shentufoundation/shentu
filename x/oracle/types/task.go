package types

import (
	"encoding/json"
	"time"

	"github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTask returns a new task.
func NewTask(
	contract string,
	function string,
	beginBlock int64,
	bounty sdk.Coins,
	description string,
	expiration time.Time,
	creator sdk.AccAddress,
	expireHeight int64,
	waitingBlocks int64,
) Task {
	return Task{
		Contract:      contract,
		Function:      function,
		BeginBlock:    beginBlock,
		Bounty:        bounty,
		Description:   description,
		Expiration:    expiration,
		Creator:       creator.String(),
		ExpireHeight:  expireHeight,
		WaitingBlocks: waitingBlocks,
		Status:        TaskStatusPending,
	}
}

// NewResponse returns a new response.
func NewResponse(score sdk.Int, operator sdk.AccAddress) Response {
	return Response{
		Operator: operator.String(),
		Score:    score,
	}
}

type Responses []Response

// String implements the Stringer interface.
func (r Responses) String() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

func NewTxTask(
	txHash []byte,
	creator string,
	bounty sdk.Coins,
	validTime time.Time,
	status TaskStatus,
) *TxTask {
	return &TxTask{
		TxHash:    txHash,
		Creator:   creator,
		Bounty:    bounty,
		ValidTime: validTime,
		Status:    status,
	}
}

type TaskI interface {
	proto.Message

	GetID() []byte
	GetCreator() string
	GetResponses() Responses
	IsExpired(ctx sdk.Context) bool
	IsValid(ctx sdk.Context) bool
	GetValidTime() (int64, time.Time)
	GetBounty() sdk.Coins
	GetStatus() TaskStatus
	GetScore() int64
	AddResponse(response Response)
	SetStatus(status TaskStatus)
	SetScore(score int64)
	ShouldAgg(ctx sdk.Context) bool
}

func (t *Task) GetID() []byte {
	return append([]byte(t.Contract), []byte(t.Function)...)
}

func (t *Task) GetCreator() string {
	return t.Creator
}

func (t *Task) GetResponses() Responses {
	return t.Responses
}

func (t *Task) IsExpired(ctx sdk.Context) bool {
	return t.Expiration.Before(ctx.BlockTime())
}

func (t *Task) GetValidTime() (int64, time.Time) {
	return t.ExpireHeight, time.Time{}
}

func (t *Task) IsValid(ctx sdk.Context) bool {
	return t.Status != TaskStatusNil && t.ExpireHeight >= ctx.BlockHeight()
}

func (t *Task) GetBounty() sdk.Coins {
	return t.Bounty
}

func (t *Task) GetStatus() TaskStatus {
	return t.Status
}

func (t *Task) GetScore() int64 {
	return t.Result.Int64()
}

func (t *Task) AddResponse(response Response) {
	t.Responses = append(t.Responses, response)
}

func (t *Task) SetStatus(status TaskStatus) {
	t.Status = status
}

func (t *Task) SetScore(score int64) {
	t.Result = sdk.NewInt(score)
}

func (t *Task) ShouldAgg(ctx sdk.Context) bool {
	return t.ExpireHeight == ctx.BlockHeight()
}

func (t *TxTask) GetID() []byte {
	return t.TxHash
}

func (t *TxTask) GetCreator() string {
	return t.Creator
}

func (t *TxTask) GetResponses() Responses {
	return t.Responses
}

func (t *TxTask) IsExpired(ctx sdk.Context) bool {
	return t.Expiration.Before(ctx.BlockTime())
}

func (t *TxTask) GetValidTime() (int64, time.Time) {
	return -1, t.ValidTime
}

func (t *TxTask) IsValid(ctx sdk.Context) bool {
	return t.Status != TaskStatusNil && !t.ValidTime.Before(ctx.BlockTime())
}

func (t *TxTask) GetBounty() sdk.Coins {
	return t.Bounty
}

func (t *TxTask) GetStatus() TaskStatus {
	return t.Status
}

func (t *TxTask) GetScore() int64 {
	return t.Score
}

func (t *TxTask) AddResponse(response Response) {
	t.Responses = append(t.Responses, response)
}

func (t *TxTask) SetStatus(status TaskStatus) {
	t.Status = status
}

func (t *TxTask) SetScore(score int64) {
	t.Score = score
}

func (t *TxTask) ShouldAgg(ctx sdk.Context) bool {
	return !t.ValidTime.After(ctx.BlockTime())
}

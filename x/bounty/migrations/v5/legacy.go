package v5

import (
	"time"

	types1 "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gogoproto/proto"
)

// Params represents the old bounty module parameters with CheckerRate field.
// This is used only for migration purposes to decode the old state.
type Params struct {
	// Minimum grant for a theorem to enter the proof period.
	MinGrant []types1.Coin `protobuf:"bytes,1,rep,name=min_grant,json=minGrant,proto3" json:"min_grant"`
	// Minimum deposit for a proof to enter the proof_hash_lock period.
	MinDeposit []types1.Coin `protobuf:"bytes,2,rep,name=min_deposit,json=minDeposit,proto3" json:"min_deposit"`
	// Duration of the theorem proof period. Initial value: 2 weeks.
	TheoremMaxProofPeriod *time.Duration `protobuf:"bytes,3,opt,name=theorem_max_proof_period,json=theoremMaxProofPeriod,proto3,stdduration" json:"theorem_max_proof_period,omitempty"`
	// Duration of the proof max lock period. 10min
	ProofMaxLockPeriod *time.Duration `protobuf:"bytes,4,opt,name=proof_max_lock_period,json=proofMaxLockPeriod,proto3,stdduration" json:"proof_max_lock_period,omitempty"`
	// use complexity_fee_rocq / complexity_fee_lean instead.
	ComplexityFee types1.Coin `protobuf:"bytes,5,opt,name=complexity_fee,json=complexityFee,proto3" json:"complexity_fee"`
	// Maximum allowed complexity to prevent overflow and ensure reasonable reward calculations.
	MaxComplexity int64 `protobuf:"varint,6,opt,name=max_complexity,json=maxComplexity,proto3" json:"max_complexity,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}

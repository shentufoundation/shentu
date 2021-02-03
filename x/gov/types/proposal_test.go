package types

import (
	"fmt"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

var fakeProposerAddress = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

var (
	testCoin = sdk.NewCoins()
	times    = []time.Time{
		time.Date(2019, 1, 1, 1, 1, 1, 1, time.UTC),
		time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
		time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC),
		time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
	}
	proposals = []Proposal{
		{gov.TextProposal{"title0", "desc0"}, 0, StatusDepositPeriod, false, fakeProposerAddress,
			gov.EmptyTallyResult(), times[0], times[1],
			sdk.NewCoins(), time.Time{}, time.Time{}},
	}
	strs = []string{
		fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Status:             %s
  Is Council Member:  %t
  Proposer Address:   %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s
  Description:        %s`, 0, "title0", types.ProposalTypeText, StatusDepositPeriod, false, fakeProposerAddress, times[0], times[1],
			testCoin, time.Time{}, time.Time{}, "desc0"),
	}
)

func TestProposalStatus_Format(t *testing.T) {
	statusDepositPeriod, _ := ProposalStatusFromString("DepositPeriod")
	statusCertifierVotingPeriod, _ := ProposalStatusFromString("CertifierVotingPeriod")
	statusPassed, _ := ProposalStatusFromString("Passed")
	statusRejected, _ := ProposalStatusFromString("Rejected")
	statusFailed, _ := ProposalStatusFromString("Failed")
	statusNil, _ := ProposalStatusFromString("")
	statusValidatorVotingPeriod, _ := ProposalStatusFromString("ValidatorVotingPeriod")
	statusDefault, _ := ProposalStatusFromString("asdasd")

	tests := []struct {
		pt                   ProposalStatus
		sprintFArgs          string
		expectedStringOutput string
	}{
		{statusDepositPeriod, "%s", "DepositPeriod"},
		{statusCertifierVotingPeriod, "%s", "CertifierVotingPeriod"},
		{statusPassed, "%s", "Passed"},
		{statusRejected, "%s", "Rejected"},
		{statusFailed, "%s", "Failed"},
		{statusNil, "%s", ""},
		{statusValidatorVotingPeriod, "%s", "ValidatorVotingPeriod"},
		{statusDefault, "%s", ""},

		{statusNil, "%v", "0"},
		{statusDepositPeriod, "%v", "1"},
		{statusCertifierVotingPeriod, "%v", "2"},
		{statusValidatorVotingPeriod, "%v", "3"},
		{statusPassed, "%v", "4"},
		{statusRejected, "%v", "5"},
		{statusFailed, "%v", "6"},
	}
	for _, tt := range tests {
		got := fmt.Sprintf(tt.sprintFArgs, tt.pt)
		require.Equal(t, tt.expectedStringOutput, got)
	}
}

func TestNewProposal(t *testing.T) {
	type args struct {
		content                 types.Content
		id                      uint64
		isProposerCouncilMember bool
		proposerAddress         sdk.AccAddress
		submitTime              time.Time
		depositEndTime          time.Time
	}
	tests := []struct {
		name string
		args args
		want Proposal
	}{
		{"proposal 0", args{gov.TextProposal{"title0", "desc0"}, 0, false, fakeProposerAddress,
			times[0], times[1]}, proposals[0]},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewProposal(tt.args.content, tt.args.id, tt.args.proposerAddress,
				tt.args.isProposerCouncilMember, tt.args.submitTime, tt.args.depositEndTime)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestProposal_String(t *testing.T) {
	type fields struct {
		Content                 types.Content
		ProposalID              uint64
		Status                  ProposalStatus
		IsProposerCouncilMember bool
		ProposerAddress         sdk.AccAddress
		FinalTallyResult        gov.TallyResult
		SubmitTime              time.Time
		DepositEndTime          time.Time
		TotalDeposit            sdk.Coins
		VotingStartTime         time.Time
		VotingEndTime           time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"proposal0", fields{gov.TextProposal{"title0", "desc0"}, 0,
			StatusDepositPeriod, false, fakeProposerAddress,
			gov.EmptyTallyResult(), times[0],
			times[1], testCoin, time.Time{}, time.Time{}}, strs[0]},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := Proposal{
				Content:                 tt.fields.Content,
				ProposalID:              tt.fields.ProposalID,
				Status:                  tt.fields.Status,
				IsProposerCouncilMember: tt.fields.IsProposerCouncilMember,
				ProposerAddress:         tt.fields.ProposerAddress,
				FinalTallyResult:        tt.fields.FinalTallyResult,
				SubmitTime:              tt.fields.SubmitTime,
				DepositEndTime:          tt.fields.DepositEndTime,
				TotalDeposit:            tt.fields.TotalDeposit,
				VotingStartTime:         tt.fields.VotingStartTime,
				VotingEndTime:           tt.fields.VotingEndTime,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidProposalStatus(t *testing.T) {
	type args struct {
		status ProposalStatus
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"deposit", args{StatusDepositPeriod}, true},
		{"voting", args{StatusCertifierVotingPeriod}, true},
		{"pass", args{StatusPassed}, true},
		{"reject", args{StatusRejected}, true},
		{"fail", args{StatusFailed}, true},
		{"voting2", args{StatusValidatorVotingPeriod}, true},
		{"invalid", args{0x07}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidProposalStatus(tt.args.status); got != tt.want {
				t.Errorf("ValidProposalStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProposalStatus_Marshal(t *testing.T) {
	tests := []struct {
		want    []byte
		name    string
		status  ProposalStatus
		wantErr bool
	}{
		{[]byte{byte(StatusDepositPeriod)}, "deposit", StatusDepositPeriod, false},
		{[]byte{byte(StatusCertifierVotingPeriod)}, "voting", StatusCertifierVotingPeriod, false},
		{[]byte{byte(StatusPassed)}, "pass", StatusPassed, false},
		{[]byte{byte(StatusRejected)}, "reject", StatusRejected, false},
		{[]byte{byte(StatusFailed)}, "fail", StatusFailed, false},
		{[]byte{byte(StatusValidatorVotingPeriod)}, "voting2", StatusValidatorVotingPeriod, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.status.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestProposalStatus_Unmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		args    args
		name    string
		status  ProposalStatus
		wantErr bool
	}{
		{args{[]byte{byte(StatusDepositPeriod)}}, "deposit", StatusDepositPeriod, false},
		{args{[]byte{byte(StatusCertifierVotingPeriod)}}, "voting", StatusCertifierVotingPeriod, false},
		{args{[]byte{byte(StatusPassed)}}, "pass", StatusPassed, false},
		{args{[]byte{byte(StatusRejected)}}, "reject", StatusRejected, false},
		{args{[]byte{byte(StatusFailed)}}, "fail", StatusFailed, false},
		{args{[]byte{byte(StatusValidatorVotingPeriod)}}, "voting2", StatusValidatorVotingPeriod, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.status.Unmarshal(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProposalStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		status  ProposalStatus
		wantErr bool
	}{
		{"deposit", StatusDepositPeriod, false},
		{"voting", StatusCertifierVotingPeriod, false},
		{"pass", StatusPassed, false},
		{"reject", StatusRejected, false},
		{"fail", StatusFailed, false},
		{"voting2", StatusValidatorVotingPeriod, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.status.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestProposalStatus_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		args    args
		name    string
		status  ProposalStatus
		wantErr bool
	}{
		{args{[]byte{byte(99)}}, "error case", StatusDepositPeriod, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.status.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTallyResult(t *testing.T) {
	type args struct {
		yes        sdk.Int
		abstain    sdk.Int
		no         sdk.Int
		noWithVeto sdk.Int
	}
	tests := []struct {
		name string
		args args
		want gov.TallyResult
	}{
		{"10,10,10,10", args{sdk.NewInt(10), sdk.NewInt(20), sdk.NewInt(30),
			sdk.NewInt(40)}, gov.TallyResult{sdk.NewInt(10), sdk.NewInt(20),
			sdk.NewInt(30), sdk.NewInt(40)}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := gov.NewTallyResult(tt.args.yes, tt.args.abstain, tt.args.no, tt.args.noWithVeto)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTallyResult_Equals(t *testing.T) {
	type fields struct {
		Yes        sdk.Int
		Abstain    sdk.Int
		No         sdk.Int
		NoWithVeto sdk.Int
	}
	type args struct {
		comp gov.TallyResult
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"equal", fields{sdk.NewInt(10), sdk.NewInt(20), sdk.NewInt(30),
			sdk.NewInt(40)}, args{gov.TallyResult{sdk.NewInt(10), sdk.NewInt(20),
			sdk.NewInt(30), sdk.NewInt(40)}}, true},
		{"not equal", fields{sdk.NewInt(10), sdk.NewInt(20), sdk.NewInt(30),
			sdk.NewInt(40)}, args{gov.TallyResult{sdk.NewInt(96), sdk.NewInt(97),
			sdk.NewInt(98), sdk.NewInt(99)}}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := gov.TallyResult{
				Yes:        tt.fields.Yes,
				Abstain:    tt.fields.Abstain,
				No:         tt.fields.No,
				NoWithVeto: tt.fields.NoWithVeto,
			}
			got := tr.Equals(tt.args.comp)
			assert.Equal(t, got, tt.want)
		})
	}
}

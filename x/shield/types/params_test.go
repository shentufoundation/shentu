package types_test

import (
	"strings"
	"testing"
	"time"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type ParamTestSuite struct {
	suite.Suite
}

func (suite *ParamTestSuite) TestPoolParamsValidation() {
	type args struct {
		ProtectionPeriod time.Duration
		MinPoolLife      time.Duration
		ShieldFeesRate   sdk.Dec
		WithdrawPeriod   time.Duration
	}

	testCases := []struct {
		name        string
		args        args
		expectPass  bool
		expectedErr string
	}{
		{
			name: "default",
			args: args{
				ProtectionPeriod: types.DefaultProtectionPeriod,
				MinPoolLife:      types.DefaultMinPoolLife,
				ShieldFeesRate:   types.DefaultShieldFeesRate,
				WithdrawPeriod:   types.DefaultWithdrawPeriod,
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "MinPoolLife <= ProtectionPeriod",
			args: args{
				ProtectionPeriod: time.Hour * 24 * 14, // 14 days
				MinPoolLife:      time.Hour * 24 * 7,  // 7 days
				ShieldFeesRate:   sdk.NewDecWithPrec(1, 2),
				WithdrawPeriod:   time.Hour * 24 * 21,
			},
			expectPass:  false,
			expectedErr: "",
		},
	}
	for _, tc := range testCases {
		tc := tc // scopelint doesn't complain
		suite.Run(tc.name, func() {
			params := types.NewPoolParams(tc.args.ProtectionPeriod, tc.args.MinPoolLife, tc.args.WithdrawPeriod, tc.args.ShieldFeesRate)
			err := types.ValidatePoolParams(params)
			if tc.expectPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
				suite.Require().True(strings.Contains(err.Error(), tc.expectedErr))
			}
		})
	}
}

func (suite *ParamTestSuite) TestClaimProposalParamsValidation() {
	type args struct {
		ClaimPeriod  time.Duration
		PayoutPeriod time.Duration
		MinDeposit   sdk.Coins
		DepositRate  sdk.Dec
		FeesRate     sdk.Dec
	}

	testCases := []struct {
		name        string
		args        args
		expectPass  bool
		expectedErr string
	}{
		{
			name: "default",
			args: args{
				ClaimPeriod:  types.DefaultClaimPeriod,
				PayoutPeriod: types.DefaultPayoutPeriod,
				MinDeposit:   types.DefaultMinClaimProposalDeposit,
				DepositRate:  types.DefaultClaimProposalDepositRate,
				FeesRate:     types.DefaultClaimProposalFeesRate,
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "PayoutPeriod  <= ClaimPeriod",
			args: args{
				ClaimPeriod:  time.Hour * 24 * 21, // 21 days
				PayoutPeriod: time.Hour * 24 * 14, // 14 days
				MinDeposit:   sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(100000000))),
				DepositRate:  sdk.NewDecWithPrec(10, 2),
				FeesRate:     sdk.NewDecWithPrec(1, 2),
			},
			expectPass:  false,
			expectedErr: "",
		},
	}
	for _, tc := range testCases {
		tc := tc // scopelint doesn't complain
		suite.Run(tc.name, func() {
			params := types.NewClaimProposalParams(tc.args.ClaimPeriod, tc.args.PayoutPeriod, tc.args.MinDeposit, tc.args.DepositRate, tc.args.FeesRate)
			err := types.ValidateClaimProposalParams(params)
			if tc.expectPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
				suite.Require().True(strings.Contains(err.Error(), tc.expectedErr))
			}
		})
	}
}

func TestParamTestSuite(t *testing.T) {
	suite.Run(t, new(ParamTestSuite))
}

package types_test

import (
	"strings"
	"testing"
	"time"

	"github.com/certikfoundation/shentu/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type ParamTestSuite struct {
	suite.Suite
}

func (suite *ParamTestSuite) TestParamValidation() {
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
	}
	for _, tc := range testCases {
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

func TestParamTestSuite(t *testing.T) {
	suite.Run(t, new(ParamTestSuite))
}

package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/v2/x/bank/client/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestLockedSendTxCmd() {
	buf := new(bytes.Buffer)
	val := s.network.Validators[0]

	testCases := []struct {
		name         string
		args         []string
		respType     proto.Message
		expectedCode uint32
		expectErr    bool
	}{
		{
			"valid transaction (gen-only)",
			[]string{
				val.Address.String(),
				val.Address.String(),
				sdk.NewCoins(
					sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
					sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
				).String(),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
				fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
			},
			&sdk.TxResponse{},
			0,
			false,
		},
		{
			"valid transaction",
			[]string{
				val.ValAddress.String(),
				val.ValAddress.String(),
				sdk.NewCoins(
					sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
					sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
				).String(),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			&sdk.TxResponse{},
			0,
			false,
		},
		{
			"not enough fees",
			[]string{
				val.ValAddress.String(),
				val.ValAddress.String(),
				sdk.NewCoins(
					sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
					sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
				).String(),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String()),
			},
			&sdk.TxResponse{},
			sdkerrors.ErrInsufficientFee.ABCICode(),
			false,
		},
		{
			"not enough gas",
			[]string{
				val.ValAddress.String(),
				val.ValAddress.String(),
				sdk.NewCoins(
					sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
					sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
				).String(),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
				"--gas=10",
			},
			&sdk.TxResponse{},
			sdkerrors.ErrOutOfGas.ABCICode(),
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx.WithOutput(buf)

			buf.Reset()

			cmd := cli.LockedSendTxCmd()
			cmd.SetErr(buf)
			cmd.SetOut(buf)
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.expectErr {
				s.Require().Error(err)
				fmt.Println(err)
				fmt.Println(cmd)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(buf.Bytes(), tc.respType), buf.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}

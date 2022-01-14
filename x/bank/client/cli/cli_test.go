package cli_test

import (
	"fmt"

	"github.com/certikfoundation/shentu/v2/app"

	// "github.com/certikfoundation/shentu/v2/common"
	// "github.com/certikfoundation/shentu/v2/x/auth/types"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/certikfoundation/shentu/v2/x/bank/client/cli"

	// "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	// cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	s.cfg = app.DefaultConfig()

	s.cfg.NumValidators = 2
	s.cfg.BondDenom = "uctk"
	s.cfg.AccountTokens = sdk.NewInt(100_000_000_000)
	s.cfg.StakingTokens = sdk.NewInt(100_000_000_000)

	s.network = network.New(s.T(), s.cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestLockedSendTx() {
	val := s.network.Validators[0]
	val1 := s.network.Validators[1]
	from := val.Address
	to := val1.Address

	fmt.Println("payment", from.String(), to.String(), s.cfg.MinGasPrices)
	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"locked send",
			[]string{
				from.String(),
				to.String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String(),

				//fmt.Sprintf("--%s=%s", cli.FlagDuration, "24h"),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, from.String()),
				// common args
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
			},
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.LockedSendTxCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

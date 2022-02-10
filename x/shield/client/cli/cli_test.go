package cli_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"

	"github.com/certikfoundation/shentu/v2/app"
	"github.com/certikfoundation/shentu/v2/x/shield/client/cli"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

var (
	pubKey = secp256k1.GenPrivKey().PubKey()
	acc1   = sdk.AccAddress(pubKey.Address())
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	s.cfg = app.DefaultConfig()

	s.cfg.NumValidators = 1
	s.cfg.BondDenom = "uctk"
	s.cfg.AccountTokens = sdk.NewInt(100_000_000_000_000_000)
	s.cfg.StakingTokens = sdk.NewInt(100_000_000_000_000_000)

	kx := keyring.NewInMemory()

	shieldAdmin, _ := kx.NewAccount("acc1", "empower ridge mystery shrimp predict alarm swear brick across funny vendor essay antique vote place lava proof gaze crush head east arch twin lady", keyring.DefaultBIP39Passphrase, sdk.FullFundraiserPath, hd.Secp256k1)
	fmt.Println("shieldacc", shieldAdmin.GetAddress().String())
	defaultGenesis := types.DefaultGenesisState()
	defaultGenesis.ShieldAdmin = shieldAdmin.GetAddress().String()

	s.cfg.GenesisState[types.ModuleName] = s.cfg.Codec.MustMarshalJSON(defaultGenesis)

	s.network = network.New(s.T(), s.cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// create a key and add this to first validator's keyring to allow sending tx
	kb := s.network.Validators[0].ClientCtx.Keyring

	info1, _ := kb.NewAccount("acc1", "empower ridge mystery shrimp predict alarm swear brick across funny vendor essay antique vote place lava proof gaze crush head east arch twin lady", keyring.DefaultBIP39Passphrase, sdk.FullFundraiserPath, hd.Secp256k1)
	pubKey = info1.GetPubKey()
	acc1 = info1.GetAddress()

	val := s.network.Validators[0]

	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		acc1,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100_000_000_000_000))),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
	)
	s.Require().NoError(err)

	// delegate coins to the first validator
	extraArgs := []string{
		val.ValAddress.String(),
		sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000000)).String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(5000))).String()),
	}
	_, err = clitestutil.ExecTestCLICmd(val.ClientCtx, stakingcli.NewDelegateCmd(), extraArgs)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestCmdDepositCollateral() {
	val := s.network.Validators[0]

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"should fail on insufficient delegation",
			[]string{
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000000000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),

				// common args
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
			},
			false,
			&sdk.TxResponse{},
			0x71,
		},
		{
			"should pass on initial collateral deposit",
			[]string{
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1),

				// common args
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
			},
			false,
			&sdk.TxResponse{},
			0x0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdDepositCollateral()
			clientCtx := val.ClientCtx

			if tc.expectErr {
				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
				s.Require().Error(err)
			} else {
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

//func (s *IntegrationTestSuite) TestCmdWithdrawCollateral() {
//	val := s.network.Validators[0]
//	//val1 := s.network.Validators[1]
//
//	testCases := []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//
//		{
//			"initial collateral deposit",
//			[]string{
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String(),
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				//fmt.Sprintf("--%s=%s", cli.FlagUnlocker, from.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdDepositCollateral()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//
//	val = s.network.Validators[0]
//	//val1 := s.network.Validators[1]
//
//	testCases = []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//		{
//			"should fail on withdrawing more than deposit",
//			[]string{
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100000000))).String(),
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				//fmt.Sprintf("--%s=%s", cli.FlagUnlocker, from.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x85,
//		},
//		{
//			"withdraw collateral",
//			[]string{
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000))).String(),
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdWithdrawCollateral()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//}
//
//func (s *IntegrationTestSuite) TestCmdCreatePool() {
//	val := s.network.Validators[0]
//
//	// deposit collateral for this new address
//	testCases := []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//
//		{
//			"initial collateral deposit",
//			[]string{
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String(),
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
//
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdDepositCollateral()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//
//	// send pool creation from this address
//
//	testCases = []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//		{
//			"should fail on invalid admin",
//			[]string{
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String(),
//				"creatorX",
//				val.Address.String(),
//				//fmt.Sprintf("--%s=%s", "deposit", sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String()),
//				fmt.Sprintf("--%s=%s", "shield-rate", "1.0"),
//				fmt.Sprintf("--%s=%s", "description", "creating a pool"),
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x65,
//		},
//		{
//			"should pass on perfect pool creation",
//			[]string{
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String(),
//				"creatorX",
//				val.Address.String(),
//				//fmt.Sprintf("--%s=%s", "deposit", sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String()),
//				fmt.Sprintf("--%s=%s", "shield-rate", "1.0"),
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdCreatePool()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//}
//
//func (s *IntegrationTestSuite) TestCmdUpdatePool() {
//	val := s.network.Validators[0]
//
//	testCases := []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//		{
//			"should fail on different admin",
//			[]string{
//				sdk.NewInt(1).String(),
//				//fmt.Sprintf("--%s=%s", "shield", sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String()),
//				//fmt.Sprintf("--%s=%s", "deposit", sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String()),
//				fmt.Sprintf("--%s=%s", "shield-rate", "1.0"),
//
//				fmt.Sprintf("--%s=%s", "description", "updating a pool"),
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x65,
//		},
//		{
//			"should pass on valid update",
//			[]string{
//				sdk.NewInt(1).String(),
//				//fmt.Sprintf("--%s=%s", "shield", sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String()),
//				//fmt.Sprintf("--%s=%s", "deposit", sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000))).String()),
//				fmt.Sprintf("--%s=%s", "shield-rate", "1.0"),
//
//				fmt.Sprintf("--%s=%s", "description", "updating a pool"),
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdUpdatePool()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//}
//
//func (s *IntegrationTestSuite) TestCmdPurchaseShield() {
//	val := s.network.Validators[0]
//
//	testCases := []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//		{
//			"should fail on too small purchase",
//			[]string{
//				sdk.NewInt(1).String(),
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
//				"buying shield",
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x8d,
//		},
//		{
//			"should pass on other than admin purchase",
//			[]string{
//				sdk.NewInt(1).String(),
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1_000_000_000))).String(),
//				"buying shield",
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//		{
//			"should pass on an admin purchase",
//			[]string{
//				sdk.NewInt(1).String(),
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1_000_000_000))).String(),
//				"buying shield",
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdPurchaseShield()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//}
//
//func (s *IntegrationTestSuite) TestCmdUnstakeFromShield() {
//	val := s.network.Validators[0]
//
//	testCases := []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//		{
//			"should fail on insufficient delegation unstake",
//			[]string{
//				sdk.NewInt(1).String(),
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String(),
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x7f,
//		},
//		{
//			"should pass on valid unstake other than admin",
//			[]string{
//				sdk.NewInt(1).String(),
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String(),
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdUnstake()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//}
//
//func (s *IntegrationTestSuite) TestCmdWithdrawAwards() {
//	val := s.network.Validators[0]
//
//	testCases := []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//		{
//			"should fail on invalid withdraw rewards",
//			[]string{
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x80,
//		},
//		{
//			"should pass on valid withdraw rewards",
//			[]string{
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdWithdrawRewards()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//}
//
//func (s *IntegrationTestSuite) TestCmdWithdrawForeignAwards() {
//	val := s.network.Validators[0]
//
//	testCases := []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//		{
//			"should pass on admin withdrawing his rewards",
//			[]string{
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String(),
//				acc1.String(),
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//		{
//			"should pass on other withdrawing admin rewards",
//			[]string{
//				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String(),
//				acc1.String(),
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdWithdrawForeignRewards()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//}
//
//func (s *IntegrationTestSuite) TestCmdUpdateSponsor() {
//	val := s.network.Validators[0]
//
//	testCases := []struct {
//		name         string
//		args         []string
//		expectErr    bool
//		respType     proto.Message
//		expectedCode uint32
//	}{
//		{
//			"should fail on invalid sponsor update",
//			[]string{
//				sdk.NewInt(1).String(),
//				"newSponsor",
//				acc1.String(),
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x65,
//		},
//		{
//			"should pass on valid sponsor update",
//			[]string{
//				sdk.NewInt(1).String(),
//				"newSponsor",
//				val.Address.String(),
//
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, acc1.String()),
//				// common args
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, "400000uctk"),
//			},
//			false,
//			&sdk.TxResponse{},
//			0x0,
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//
//		s.Run(tc.name, func() {
//			cmd := cli.GetCmdUpdateSponsor()
//			clientCtx := val.ClientCtx
//
//			if tc.expectErr {
//				_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().Error(err)
//			} else {
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//				s.Require().NoError(err, out.String())
//				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
//
//				txResp := tc.respType.(*sdk.TxResponse)
//				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
//			}
//		})
//	}
//}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

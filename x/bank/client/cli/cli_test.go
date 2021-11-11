package cli_test

import (
	//"context"
	"fmt"

	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	//tmcli "github.com/tendermint/tendermint/libs/cli"

	//	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	//clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"github.com/cosmos/cosmos-sdk/types/query"
	//"github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	genesisState := cfg.GenesisState
	cfg.NumValidators = 1

	var bankGenesis types.GenesisState
	s.Require().NoError(cfg.Codec.UnmarshalJSON(genesisState[types.ModuleName], &bankGenesis))

	bankGenesis.DenomMetadata = []types.Metadata{
		{
			Description: "The native staking token of the Shentu Chain.",
			DenomUnits: []*types.DenomUnit{
				{
					Denom:    "uctk",
					Exponent: 0,
					Aliases:  []string{"microctk"},
				},
				{
					Denom:    "ctk",
					Exponent: 6,
					Aliases:  []string{"CTK"},
				},
			},
			Base:    "uctk",
			Display: "ctk",
		},
		{
			Description: "Ethereum mainnet token",
			DenomUnits: []*types.DenomUnit{
				{
					Denom:    "wei",
					Exponent: 0,
				},
				{
					Denom:    "eth",
					Exponent: 6,
					Aliases:  []string{"ETH"},
				},
			},
			Base:    "wei",
			Display: "eth",
		},
	}

	bankGenesisBz, err := cfg.Codec.MarshalJSON(&bankGenesis)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = bankGenesisBz
	cfg.GenesisState = genesisState

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestLockedSendTxCmd() {
	val := s.network.Validators[0]

	testCases := []struct {
		args         []string
		name         string
		from, to     sdk.AccAddress
		respType     proto.Message
		amount       sdk.Coins
		expectedCode uint32
		expectErr    bool
	}{
		{
			[]string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			"valid transaction",
			val.Address,
			val.Address,
			&sdk.TxResponse{},
			sdk.NewCoins(
				sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
			),
			0,
			false,
		},
		{
			[]string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1))).String()),
			},
			"not enough fees",
			val.Address,
			val.Address,
			&sdk.TxResponse{},
			sdk.NewCoins(
				sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
			),
			sdkerrors.ErrInsufficientFee.ABCICode(),
			false,
		},
		{
			[]string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
				"--gas=10",
			},
			"not enough gas",
			val.Address,
			val.Address,
			&sdk.TxResponse{},
			sdk.NewCoins(
				sdk.NewCoin(fmt.Sprintf("%stoken", val.Moniker), sdk.NewInt(10)),
				sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
			),
			sdkerrors.ErrOutOfGas.ABCICode(),
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx

			bz, err := banktestutil.MsgSendExec(clientCtx, tc.from, tc.to, tc.amount, tc.args...)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				s.Require().NoError(clientCtx.JSONCodec.UnmarshalJSON(bz.Bytes(), tc.respType), bz.String())
				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func NewCoin(denom string, amount sdk.Int) *sdk.Coin {
	coin := sdk.NewCoin(denom, amount)
	return &coin
}

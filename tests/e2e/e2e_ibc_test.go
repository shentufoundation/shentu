package e2e

import (
	"context"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testIBCTokanTransfer() {
	s.Run("send_uctk_to_chainB", func() {
		var (
			balances      sdk.Coins
			err           error
			beforeBalance int64
			ibcStakeDenom string
		)

		address, err := s.chainA.validators[0].keyInfo.GetAddress()
		s.Require().NoError(err)
		sender := address.String()

		address, err = s.chainB.validators[0].keyInfo.GetAddress()
		s.Require().NoError(err)
		recipient := address.String()

		chainBAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainB.id][0].GetHostPort("1317/tcp"))

		s.Require().Eventually(
			func() bool {
				balances, err = queryShentuAllBalances(chainBAPIEndpoint, recipient)
				s.Require().NoError(err)
				return balances.Len() != 0
			},
			time.Minute,
			5*time.Second,
		)
		for _, c := range balances {
			if strings.Contains(c.Denom, "ibc/") {
				beforeBalance = c.Amount.Int64()
				break
			}
		}

		token := sdk.NewInt64Coin(uctkDenom, 3300000000) // 3,300ctk
		s.sendIBC(s.chainA, 0, sender, recipient, "", token, feesAmountCoin)
		s.hermesClearPacket(hermesConfigWithGasPrices, s.chainA.id, transferPort, transferChannel)

		// require the recipient account receives the IBC tokens (IBC packets ACKd)
		s.Require().Eventually(
			func() bool {
				balances, err = queryShentuAllBalances(chainBAPIEndpoint, recipient)
				s.Require().NoError(err)
				return balances.Len() > 1
			},
			time.Minute,
			5*time.Second,
		)

		for _, c := range balances {
			if strings.Contains(c.Denom, "ibc/") {
				ibcStakeDenom = c.Denom
				s.Require().Equal(token.Amount.Int64()+beforeBalance, c.Amount.Int64())
				break
			}
		}

		s.Require().NotEmpty(ibcStakeDenom)
	})
}

func (s *IntegrationTestSuite) sendIBC(c *chain, valIdx int, sender, recipient, note string, token, fees sdk.Coin) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ibcCmd := []string{
		shentuBinary,
		txCommand,
		"ibc-transfer",
		"transfer",
		transferPort,
		transferChannel,
		recipient,
		token.String(),
		fmt.Sprintf("--from=%s", sender),
		fmt.Sprintf("--fees=%s", fees.String()),
		fmt.Sprintf("--chain-id=%s", c.id),
		fmt.Sprintf("--memo=%s", note),
		"--keyring-backend=test",
		"--broadcast-mode=sync",
		"--output=json",
		"-y",
	}
	s.T().Logf("sending %s from %s (%s) to %s (%s) with memo %s", token.String(), s.chainA.id, sender, s.chainB.id, recipient, note)
	s.execShentuTxCmd(ctx, c, ibcCmd, valIdx, s.execValidationDefault(c, valIdx))
	s.T().Log("successfully sent IBC tokens")
}

func (s *IntegrationTestSuite) hermesClearPacket(configPath, chainID, portID, channelID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	hermesCmd := []string{
		hermesBinary,
		"--json",
		fmt.Sprintf("--config=%s", configPath),
		"clear",
		"packets",
		fmt.Sprintf("--chain=%s", chainID),
		fmt.Sprintf("--channel=%s", channelID),
		fmt.Sprintf("--port=%s", portID),
	}

	s.executeHermesCommand(ctx, hermesCmd, nil)
	s.T().Log("successfully cleared IBC packets")
}

func (s *IntegrationTestSuite) createConnection() {
	s.T().Logf("creating connection between %s and %s", s.chainA.id, s.chainB.id)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	hermesCmd := []string{
		hermesBinary,
		"--json",
		"create",
		"connection",
		"--a-chain",
		s.chainA.id,
		"--b-chain",
		s.chainB.id,
	}
	s.executeHermesCommand(ctx, hermesCmd, nil)
	s.T().Logf("successfully created connection between %s and %s", s.chainA.id, s.chainB.id)
}

func (s *IntegrationTestSuite) createChannel() {
	s.T().Logf("creating channel between %s and %s", s.chainA.id, s.chainB.id)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	hermesCmd := []string{
		hermesBinary,
		"--json",
		"create",
		"channel",
		"--a-chain",
		s.chainA.id,
		"--a-connection",
		"connection-0",
		"--a-port",
		transferPort,
		"--b-port",
		transferPort,
		"--channel-version",
		"ics20-1",
		"--order",
		"unordered",
	}
	s.executeHermesCommand(ctx, hermesCmd, nil)
	s.T().Logf("successfully created channel between %s and %s", s.chainA.id, s.chainB.id)
}

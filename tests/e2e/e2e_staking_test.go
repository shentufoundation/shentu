package e2e

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (s *IntegrationTestSuite) executeDelegate(c *chain, valIdx int, amount, valOperAddress, delegatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx staking delegate %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		stakingtypes.ModuleName,
		"delegate",
		valOperAddress,
		amount,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, delegatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully delegated %s to %s", delegatorAddr, amount, valOperAddress)
}

func (s *IntegrationTestSuite) executeUnbond(c *chain, valIdx int, amount, valOperAddress, delegatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx staking unbond %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		stakingtypes.ModuleName,
		"unbond",
		valOperAddress,
		amount,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, delegatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully unbond %s from %s", delegatorAddr, amount, valOperAddress)
}

func queryDelegation(endpoint, validatorAddr, delegatorAddr string) (*stakingtypes.QueryDelegationResponse, error) {
	grpcReq := &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := stakingtypes.NewQueryClient(conn)

	grpcRsp, err := client.Delegation(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}

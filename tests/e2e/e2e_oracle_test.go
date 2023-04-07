package e2e

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func (s *IntegrationTestSuite) executeOracleCreateOperator(c *chain, valIdx int, operatorAddr, collateral, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	s.T().Logf("Executing shentu tx create operator %s", c.id)
	defer cancel()

	command := []string{
		shentuBinary,
		txCommand,
		types.ModuleName,
		"create-operator",
		operatorAddr,
		collateral,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, operatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}
	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("successfully add operator on %s", operatorAddr)
}

func (s *IntegrationTestSuite) executeOracleCreateTxTask(c *chain, valIdx int, txBytes, chainId, bounty, valTime, creatorAddr, fees string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	s.T().Logf("Executing shentu tx create tx-task %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		types.ModuleName,
		"create-txtask",
		txBytes,
		chainId,
		bounty,
		valTime,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, creatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	stdOut, _ := s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	txResp := sdk.TxResponse{}
	if err := cdc.UnmarshalJSON(stdOut, &txResp); err != nil {
		return "", err
	}

	s.T().Logf("%s successfully submit tx-task on %s", creatorAddr, txResp.TxHash)
	return txResp.TxHash, nil
}

func (s *IntegrationTestSuite) executeOracleRespondTxTask(c *chain, valIdx, score int, taskHash, operatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	s.T().Logf("Executing shentu tx respond tx-task %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		types.ModuleName,
		"respond-to-txtask",
		taskHash,
		fmt.Sprintf("%d", score),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, operatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}
	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("successfully respond by operator %s", operatorAddr)
}

func (s *IntegrationTestSuite) executeOracleCreateTask(c *chain, valIdx int, contract, function, bounty, creatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	s.T().Logf("Executing shentu tx create task %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		types.ModuleName,
		"create-task",
		contract,
		function,
		bounty,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, creatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}
	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("successfully create task by %s", creatorAddr)
}

func (s *IntegrationTestSuite) executeOracleRespondTask(c *chain, valIdx, score int, contract, function, operatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	s.T().Logf("Executing shentu tx respond task %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		types.ModuleName,
		"respond-to-task",
		contract,
		function,
		fmt.Sprintf("%d", score),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, operatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}
	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("successfully respond by operator %s", operatorAddr)
}

func (s *IntegrationTestSuite) executeOracleClaimReward(c *chain, valIdx int, operatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	s.T().Logf("Executing shentu tx claim reward %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		types.ModuleName,
		"claim-reward",
		operatorAddr,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, operatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}
	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("successfully claim reward by operator %s", operatorAddr)
}

func (s *IntegrationTestSuite) executeOracleRemoveOperator(c *chain, valIdx int, operatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	s.T().Logf("Executing shentu tx create operator %s", c.id)
	defer cancel()

	command := []string{
		shentuBinary,
		txCommand,
		types.ModuleName,
		"remove-operator",
		operatorAddr,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, operatorAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}
	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("successfully remove operator on %s", operatorAddr)
}

func queryOracleTaskHash(endpoint, txHash string) (string, error) {
	txRsp, err := getShentuTx(endpoint, txHash)
	if err != nil {
		return "", err
	}
	for _, log := range txRsp.Logs {
		for _, event := range log.Events {
			if event.Type == "create_tx_task" {
				for _, attribute := range event.Attributes {
					if attribute.Key == "atx_hash" {
						return attribute.Value, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("field not find")
}

func queryOracleOperator(endpoint, operatorAddr string) (*types.QueryOperatorResponse, error) {
	grpcReq := &types.QueryOperatorRequest{
		Address: operatorAddr,
	}
	conn, _ := connectGrpc(endpoint)
	defer conn.Close()
	client := types.NewQueryClient(conn)
	grpcRsp, err := client.Operator(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	return grpcRsp, nil
}

func queryOracleTxTask(endpoint, ataskHash string) (*types.QueryTxTaskResponse, error) {
	grpcReq := &types.QueryTxTaskRequest{
		AtxHash: ataskHash,
	}
	conn, _ := connectGrpc(endpoint)
	defer conn.Close()
	client := types.NewQueryClient(conn)
	grpcRsp, err := client.TxTask(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	return grpcRsp, nil
}

func queryOracleTask(endpoint, contract, function string) (*types.QueryTaskResponse, error) {
	grpcReq := &types.QueryTaskRequest{
		Contract: contract,
		Function: function,
	}
	conn, _ := connectGrpc(endpoint)
	defer conn.Close()
	client := types.NewQueryClient(conn)
	grpcRsp, err := client.Task(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	return grpcRsp, nil
}

func queryOracleOperators(endpoint string) (*types.QueryOperatorsResponse, error) {
	grpcReq := &types.QueryOperatorsRequest{}
	conn, _ := connectGrpc(endpoint)
	defer conn.Close()
	client := types.NewQueryClient(conn)
	grpcRsp, err := client.Operators(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	return grpcRsp, nil
}

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ory/dockertest/v3/docker"
)

func (s *IntegrationTestSuite) connectIBCChains() {
	s.T().Logf("connecting %s and %s chains via IBC", s.chainA.id, s.chainB.id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	exec, err := s.dkrPool.Client.CreateExec(docker.CreateExecOptions{
		Context:      ctx,
		AttachStdout: true,
		AttachStderr: true,
		Container:    s.hermesResource.Container.ID,
		User:         "root",
		Cmd: []string{
			"hermes",
			"create",
			"channel",
			s.chainA.id,
			s.chainB.id,
			"--port-a=transfer",
			"--port-b=transfer",
		},
	})
	s.Require().NoError(err)

	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	err = s.dkrPool.Client.StartExec(exec.ID, docker.StartExecOptions{
		Context:      ctx,
		Detach:       false,
		OutputStream: &outBuf,
		ErrorStream:  &errBuf,
	})
	s.Require().NoErrorf(
		err,
		"failed to connect chains; stdout: %s, stderr: %s", outBuf.String(), errBuf.String(),
	)

	s.T().Logf("connected %s and %s chains via IBC", s.chainA.id, s.chainB.id)
}

func (s *IntegrationTestSuite) sendIBC(srcChainID, dstChainID, recipient string, token sdk.Coin) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("sending %s from %s to %s (%s)", token, srcChainID, dstChainID, recipient)

	exec, err := s.dkrPool.Client.CreateExec(docker.CreateExecOptions{
		Context:      ctx,
		AttachStdout: true,
		AttachStderr: true,
		Container:    s.hermesResource.Container.ID,
		User:         "root",
		Cmd: []string{
			"hermes",
			"tx",
			"raw",
			"ft-transfer",
			dstChainID,
			srcChainID,
			"transfer",  // source chain port ID
			"channel-0", // since only one connection/channel exists, assume 0
			token.Amount.String(),
			fmt.Sprintf("--denom=%s", token.Denom),
			fmt.Sprintf("--receiver=%s", recipient),
			"--timeout-height-offset=1000",
		},
	})
	s.Require().NoError(err)

	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	err = s.dkrPool.Client.StartExec(exec.ID, docker.StartExecOptions{
		Context:      ctx,
		Detach:       false,
		OutputStream: &outBuf,
		ErrorStream:  &errBuf,
	})
	s.Require().NoErrorf(
		err,
		"failed to send IBC tokens; stdout: %s, stderr: %s", outBuf.String(), errBuf.String(),
	)

	s.T().Log("successfully sent IBC tokens")
}

func (s *IntegrationTestSuite) getLatestBlockHeight(endpoint string) (int64, error) {
	grpcReq := &tmservice.GetLatestBlockRequest{}

	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := tmservice.NewServiceClient(conn)

	grpcRsp, err := client.GetLatestBlock(context.Background(), grpcReq)
	if err != nil {
		return 0, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp.GetBlock().GetHeader().Height, nil
}

func (s *IntegrationTestSuite) execShentuTxCmd(ctx context.Context, c *chain, cmd []string, valIdx int, validation func([]byte, []byte) bool) {
	if validation == nil {
		validation = s.defaultExecValidation(s.chainA, 0)
	}
	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)
	exec, err := s.dkrPool.Client.CreateExec(docker.CreateExecOptions{
		Context:      ctx,
		AttachStdout: true,
		AttachStderr: true,
		Container:    s.valResources[c.id][valIdx].Container.ID,
		User:         "root",
		Cmd:          cmd,
	})
	s.Require().NoError(err)

	err = s.dkrPool.Client.StartExec(exec.ID, docker.StartExecOptions{
		Context:      ctx,
		Detach:       false,
		OutputStream: &outBuf,
		ErrorStream:  &errBuf,
	})
	s.Require().NoError(err)

	stdOut := outBuf.Bytes()
	stdErr := errBuf.Bytes()
	if !validation(stdOut, stdErr) {
		s.Require().FailNowf("tx validation failed", "stdout: %s, stderr: %s",
			string(stdOut), string(stdErr))
	}
}

func (s *IntegrationTestSuite) defaultExecValidation(chain *chain, valIdx int) func([]byte, []byte) bool {
	return func(stdOut []byte, stdErr []byte) bool {
		var txResp sdk.TxResponse
		if err := cdc.UnmarshalJSON(stdOut, &txResp); err != nil {
			return false
		}
		if strings.Contains(txResp.String(), "code: 0") || txResp.Code == 0 {
			endpoint := s.valResources[chain.id][valIdx].GetHostPort("9090/tcp")
			s.Require().Eventually(
				func() bool {
					err := queryShentuTx(endpoint, txResp.TxHash)
					return err == nil
				},
				time.Minute,
				5*time.Second,
				"stdOut: %s, stdErr: %s",
				string(stdOut), string(stdErr),
			)
			return true
		}
		return false
	}
}

func connectGrpc(endpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect %s: %v", endpoint, err)
	}
	return conn, nil
}

func queryShentuTx(endpoint, txHash string) error {
	grpcReq := &tx.GetTxRequest{
		Hash: txHash,
	}

	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := tx.NewServiceClient(conn)

	grpcRsp, err := client.GetTx(context.Background(), grpcReq)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	if grpcRsp.GetTxResponse().Code != 0 {
		return fmt.Errorf("tx %s failed with status code %v", txHash, grpcRsp.GetTxResponse().Code)
	}
	return nil
}

func queryShentuAllBalances(endpoint, addr string) (sdk.Coins, error) {
	grpcReq := &banktypes.QueryAllBalancesRequest{
		Address: addr,
	}

	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := banktypes.NewQueryClient(conn)

	grpcRsp, err := client.AllBalances(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp.GetBalances(), nil
}

func queryShentuDenomBalance(endpoint, addr, denom string) (sdk.Coin, error) {
	var zeroCoin sdk.Coin
	grpcReq := &banktypes.QueryBalanceRequest{
		Address: addr,
		Denom:   denom,
	}

	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := banktypes.NewQueryClient(conn)

	grpcRsp, err := client.Balance(context.Background(), grpcReq)
	if err != nil {
		return zeroCoin, fmt.Errorf("failed to execute request: %w", err)
	}

	return *grpcRsp.GetBalance(), nil
}

func configFile(filename string) string {
	return filepath.Join(shentuHome, "config", filename)
}

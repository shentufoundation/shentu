package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ory/dockertest/v3/docker"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type flagOption func(map[string]string)

func withExtraFlag(key, value string) flagOption {
	return func(flags map[string]string) {
		flags[key] = value
	}
}

func applyOptions(options []flagOption) map[string]string {
	flags := make(map[string]string)
	for _, opt := range options {
		opt(flags)
	}
	return flags
}

func (s *IntegrationTestSuite) execShentuTxCmd(ctx context.Context, c *chain, cmd []string, valIdx int, validation func([]byte, []byte) bool) ([]byte, []byte) {
	if validation == nil {
		validation = s.execValidationDefault(s.chainA, 0)
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
		s.Require().FailNowf("Exec validation failed", "stdout: %s, stderr: %s", string(stdOut), string(stdErr))
	}
	return stdOut, stdErr
}

func (s *IntegrationTestSuite) executeHermesCommand(ctx context.Context, hermesCmd []string, validation func([]byte, []byte) bool) ([]byte, []byte) {
	if validation == nil {
		validation = s.execValidationHermes()
	}
	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)
	s.T().Logf("Executing hermes command: %s", strings.Join(hermesCmd, " "))
	exec, err := s.dkrPool.Client.CreateExec(docker.CreateExecOptions{
		Context:      ctx,
		AttachStdout: true,
		AttachStderr: true,
		Container:    s.hermesResource.Container.ID,
		User:         "root",
		Cmd:          hermesCmd,
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
		s.Require().FailNowf("Exec validation failed", "stdout: %s, stderr: %s", string(stdOut), string(stdErr))
	}
	return stdOut, stdErr
}

func (s *IntegrationTestSuite) execValidationDefault(chain *chain, valIdx int) func([]byte, []byte) bool {
	return func(stdout, stderr []byte) bool {
		var txResp sdk.TxResponse
		if err := cdc.UnmarshalJSON(stdout, &txResp); err != nil {
			return false
		}
		if strings.Contains(txResp.String(), "code: 0") || txResp.Code == 0 {
			endpoint := fmt.Sprintf("http://%s", s.valResources[chain.id][valIdx].GetHostPort("1317/tcp"))
			s.Require().Eventually(
				func() bool {
					err := queryShentuTx(endpoint, txResp.TxHash)
					return err == nil
				},
				time.Minute,
				5*time.Second,
				"stdout: %s, stderr: %s",
				string(stdout), string(stderr),
			)
			return true
		}
		return false
	}
}

func (s *IntegrationTestSuite) execValidationError(chain *chain, valIdx int) func([]byte, []byte) bool {
	return func(stdout, stderr []byte) bool {
		var txResp sdk.TxResponse
		if err := cdc.UnmarshalJSON(stdout, &txResp); err != nil {
			return true
		}
		if txResp.Code != 0 {
			return true
		}
		endpoint := fmt.Sprintf("http://%s", s.valResources[chain.id][valIdx].GetHostPort("1317/tcp"))
		s.Require().Eventually(
			func() bool {
				err := queryShentuTx(endpoint, txResp.TxHash)
				return err != nil
			},
			time.Minute,
			5*time.Second,
			"stdout: %s, stderr: %s",
			string(stdout), string(stderr),
		)
		return true
	}
}

func (s *IntegrationTestSuite) execValidationHermes() func([]byte, []byte) bool {
	return func(stdout, stderr []byte) bool {
		var out map[string]interface{}
		lines := bytes.Split(stdout, []byte("\n"))
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			err := json.Unmarshal(line, &out)
			if err != nil {
				return false
			}
			if s := out["status"]; s != nil && s != "success" {
				return false
			} else if s == "success" {
				return true
			}
		}
		return true
	}
}

func (s *IntegrationTestSuite) execBankSend(c *chain, valIdx int, from, to string, amount, fees sdk.Coin, expectError bool, opt ...flagOption) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bank",
		"send",
		from,
		to,
		amount.String(),
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	flags := applyOptions(opt)
	for k, v := range flags {
		cmd = append(cmd, fmt.Sprintf("--%s", k), fmt.Sprintf("%s", v))
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed bank send from %s to %s %s", from, to, amount.String())
}

func (s *IntegrationTestSuite) execBankMultiSend(c *chain, valIdx int, from string, to []string, amount, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bank",
		"multi-send",
		from,
	}
	cmd2 := []string{
		amount.String(),
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	cmd = append(cmd, to...)
	cmd = append(cmd, cmd2...)
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed bank multi-send from %s to %s %s", from, strings.Join(to, ","), amount.String())
}

func (s *IntegrationTestSuite) execSetWithdrawAddress(c *chain, valIdx int, delegator, withdrawer string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"distribution",
		"set-withdraw-addr",
		withdrawer,
		"--from", delegator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed set withdraw address of %s to %s", delegator, withdrawer)
}

func (s *IntegrationTestSuite) execWithdrawReward(c *chain, valIdx int, delegator, validator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"distribution",
		"withdraw-rewards",
		validator,
		"--from", delegator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed withdraw rewards of %s from %s", delegator, validator)
}

func (s *IntegrationTestSuite) execFundCommunityPool(c *chain, valIdx int, account string, amount, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"distribution",
		"fund-community-pool",
		amount.String(),
		"--from", account,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed fund community pool from %s %s", account, amount.String())
}

func (s *IntegrationTestSuite) execFeeGrant(c *chain, valIdx int, granter, grantee string, limit, fees sdk.Coin, expectError bool, opt ...flagOption) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"feegrant",
		"grant",
		granter,
		grantee,
		"--from", granter,
		"--spend-limit", limit.String(),
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}

	flags := applyOptions(opt)
	for k, v := range flags {
		cmd = append(cmd, fmt.Sprintf("--%s", k), fmt.Sprintf("%s", v))
	}

	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed fee grant from %s to %s %s", granter, grantee, limit.String())
}

func (s *IntegrationTestSuite) execDelegate(c *chain, valIdx int, delegator, validator string, amount, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"staking",
		"delegate",
		validator,
		amount.String(),
		"--from", delegator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed delegate from %s to %s %s", delegator, validator, amount.String())
}

func (s *IntegrationTestSuite) execCreateProgram(c *chain, valIdx int, programID, name, desc, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"create-program",
		"--program-id", programID,
		"--name", name,
		"--detail", desc,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed create program %s by %s", programID, creator)
}

func (s *IntegrationTestSuite) execEditProgram(c *chain, valIdx int, name, desc, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"edit-program",
		"--name", name,
		"--detail", desc,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed edit program %s by %s", name, creator)
}

func (s *IntegrationTestSuite) execActivateProgram(c *chain, valIdx int, programID, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"activate-program",
		programID,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed activate program %s by %s", programID, creator)
}

func (s *IntegrationTestSuite) execCloseProgram(c *chain, valIdx int, programID, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"close-program",
		programID,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed close program %s by %s", programID, creator)
}

func (s *IntegrationTestSuite) execSubmitFinding(c *chain, valIdx int, programID, findingID, severity, desc, poc, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"submit-finding",
		"--program-id", programID,
		"--finding-id", findingID,
		"--severity-level", severity,
		"--desc", desc,
		"--poc", poc,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed submit finding %s by %s", findingID, creator)
}

func (s *IntegrationTestSuite) execEditFinding(c *chain, valIdx int, findingID, severity, desc, poc, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"edit-finding",
		"--finding-id", findingID,
		"--severity-level", severity,
		"--desc", desc,
		"--poc", poc,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed edit finding %s by %s", findingID, creator)
}

func (s *IntegrationTestSuite) execActivateFinding(c *chain, valIdx int, findingID, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"activate-finding",
		findingID,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed activate finding %s by %s", findingID, creator)
}

func (s *IntegrationTestSuite) execConfirmFinding(c *chain, valIdx int, findingID, fingerprint, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"confirm-finding",
		findingID,
		"--fingerprint", fingerprint,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed confirm finding %s by %s", findingID, creator)
}

func (s *IntegrationTestSuite) execEditPayment(c *chain, valIdx int, findingID, payment, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"edit-finding",
		"--finding-id", findingID,
		"--payment-hash", payment,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed edit payment of finding %s by %s", findingID, creator)
}

func (s *IntegrationTestSuite) execConfirmPayment(c *chain, valIdx int, findingID, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"confirm-finding-paid",
		findingID,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed confirm payment of finding %s by %s", findingID, creator)
}

func (s *IntegrationTestSuite) execPublishFinding(c *chain, valIdx int, findingID, desc, poc, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"publish-finding",
		findingID,
		"--desc", desc,
		"--poc", poc,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed publish finding %s by %s", findingID, creator)
}

func (s *IntegrationTestSuite) execCloseFinding(c *chain, valIdx int, findingID, creator string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"bounty",
		"close-finding",
		findingID,
		"--from", creator,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully executed close finding %s by %s", findingID, creator)
}

func (s *IntegrationTestSuite) execIssueCertificate(c *chain, valIdx int, content, certificate, desc, certifier string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"cert",
		"issue-certificate",
		certificate,
		content,
		"--description", desc,
		"--from", certifier,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully issued %s certificate to %s by %s", certificate, content, certifier)
	if !expectError {
		certificateCounter++
	}
}

func (s *IntegrationTestSuite) execRevokeCertificate(c *chain, valIdx int, certificateID, certifier string, fees sdk.Coin, expectError bool) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"cert",
		"revoke-certificate",
		certificateID,
		"--from", certifier,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	validation := s.execValidationDefault(c, valIdx)
	if expectError {
		validation = s.execValidationError(c, valIdx)
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, validation)
	s.T().Logf("Successfully revoked certificate %s by %s", certificateID, certifier)
}

func (s *IntegrationTestSuite) execSubmitProposal(c *chain, valIdx int, proposalFileName, proposer string, fees sdk.Coin) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"gov",
		"submit-proposal",
		configFile(proposalFileName),
		"--from", proposer,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	proposalCounter++
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, s.execValidationDefault(c, valIdx))
	s.T().Logf("Successfully submitted proposal %d", proposalCounter)
}

func (s *IntegrationTestSuite) execVoteProposal(c *chain, valIdx int, proposalID uint64, voter, option string, fees sdk.Coin) {
	cmd := []string{
		shentuBinary,
		txCommand,
		"gov",
		"vote",
		fmt.Sprintf("%d", proposalID),
		option,
		"--from", voter,
		"--fees", fees.String(),
		"--chain-id", c.id,
		"--keyring-backend", "test",
		"--output", "json",
		"--yes",
	}
	s.execShentuTxCmd(context.Background(), c, cmd, valIdx, s.execValidationDefault(c, valIdx))
	s.T().Logf("Successfully voted proposal %d %s by %s", proposalID, option, voter)
}

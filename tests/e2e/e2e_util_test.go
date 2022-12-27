package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ory/dockertest/v3/docker"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
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

func (s *IntegrationTestSuite) getLatestBlockHeight(c *chain, valIdx int) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	type statusInfo struct {
		StatusInfo struct {
			LatestHeight string `json:"latest_block_height"`
		} `json:"SyncInfo"`
	}

	var latestHeight int
	command := []string{
		shentuBinary,
		"status",
	}

	s.execShentuTxCmd(ctx, c, command, valIdx, func(stdOut []byte, stdErr []byte) bool {
		var (
			err   error
			block statusInfo
		)
		s.Require().NoError(json.Unmarshal(stdErr, &block))
		latestHeight, err = strconv.Atoi(block.StatusInfo.LatestHeight)
		s.Require().NoError(err)
		return latestHeight > 0
	})
	return latestHeight
}

func (s *IntegrationTestSuite) executeIssueCertificate(c *chain, valIdx int, certificateType, content, certifierAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx issue certificate %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		certtypes.ModuleName,
		"issue-certificate",
		certificateType,
		content,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, certifierAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully issue %s certificate to %s", certifierAddr, certificateType, content)
}

func (s *IntegrationTestSuite) executeSubmitUpgradeProposal(c *chain, valIdx, upgradeHeight int, submitterAddr, proposalName, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx submit proposal %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		govtypes.ModuleName,
		"submit-proposal",
		"software-upgrade",
		proposalName,
		fmt.Sprintf("--upgrade-height=%d", upgradeHeight),
		fmt.Sprintf("--title=\"title of %s\"", proposalName),
		fmt.Sprintf("--description=\"description of %s\"", proposalName),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully submit %s proposal", submitterAddr, proposalName)
}

func (s *IntegrationTestSuite) executeSubmitClaimProposal(c *chain, valIdx int, proposalFile, submitterAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx submit proposal %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		govtypes.ModuleName,
		"submit-proposal",
		"shield-claim",
		proposalFile,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully submit claim proposal %s", submitterAddr, proposalFile)
}

func (s *IntegrationTestSuite) executeDepositProposal(c *chain, valIdx int, submitterAddr string, proposalId int, amount, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx deposit proposal %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		govtypes.ModuleName,
		"deposit",
		fmt.Sprintf("%d", proposalId),
		amount,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully deposit proposal %d %s", submitterAddr, proposalId, amount)
}

func (s *IntegrationTestSuite) executeVoteProposal(c *chain, valIdx int, submitterAddr string, proposalId int, vote, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx vote proposal %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		govtypes.ModuleName,
		"vote",
		fmt.Sprintf("%d", proposalId),
		vote,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully vote proposal %d %s", submitterAddr, proposalId, vote)
}

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

func (s *IntegrationTestSuite) executeDepositCollateral(c *chain, valIdx int, submitterAddr, amount, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx shield deposit collateral %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		shieldtypes.ModuleName,
		"deposit-collateral",
		amount,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully deposit %s", submitterAddr, amount)
}

func (s *IntegrationTestSuite) executeCreatePool(c *chain, valIdx int, poolAmount, poolName, sponsorAddr, shieldLimit, submitterAddr, nativeAmount, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx shield create pool %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		shieldtypes.ModuleName,
		"create-pool",
		poolAmount,
		poolName,
		sponsorAddr,
		fmt.Sprintf("--%s=%s", "shield-limit", shieldLimit),
		fmt.Sprintf("--%s=%s", "native-deposit", nativeAmount),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully create pool %s", submitterAddr, poolName)
}

func (s *IntegrationTestSuite) executePurchaseShield(c *chain, valIdx, poolId int, shieldAmount, description, submitterAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx shield purchase %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		shieldtypes.ModuleName,
		"purchase",
		fmt.Sprintf("%d", poolId),
		shieldAmount,
		description,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully purchase shield at pool %d", submitterAddr, poolId)
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
			endpoint := fmt.Sprintf("http://%s", s.valResources[chain.id][valIdx].GetHostPort("1317/tcp"))
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

func (s *IntegrationTestSuite) writeClaimProposal(c *chain, valIdx, poolId, purchaseId int, fileName string) string {
	type ClaimLoss struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
	type ClaimProposal struct {
		PoolId      int         `json:"pool_id"`
		PurchaseId  int         `json:"purchase_id"`
		Evidence    string      `json:"evidence"`
		Description string      `json:"description"`
		Loss        []ClaimLoss `json:"loss"`
	}

	var loss = ClaimLoss{
		Denom:  "uctk",
		Amount: "100000000",
	}
	var proposal = &ClaimProposal{
		PoolId:      poolId,
		PurchaseId:  purchaseId,
		Evidence:    "Attack happened on <time> caused loss of <amount> to <account> by <txhashes>",
		Description: "Details of the attack",
	}
	proposal.Loss = append(proposal.Loss, loss)

	proposalByte, err := json.Marshal(proposal)
	s.Require().NoError(err)

	path := filepath.Join(c.validators[valIdx].configDir(), "config", fileName)

	_, err = os.Create(path)
	s.Require().NoError(err)

	os.WriteFile(path, proposalByte, 0o600)
	return path
}

func queryShentuTx(endpoint, txHash string) error {
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", endpoint, txHash))
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("tx query returned non-200 status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	txResp := result["tx_response"].(map[string]interface{})
	if v := txResp["code"]; v.(float64) != 0 {
		return fmt.Errorf("tx %s failed with status code %v", txHash, v)
	}

	return nil
}

func queryShentuAllBalances(endpoint, addr string) (sdk.Coins, error) {
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", endpoint, addr))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var balancesResp banktypes.QueryAllBalancesResponse
	if err := cdc.UnmarshalJSON(bz, &balancesResp); err != nil {
		return nil, err
	}

	return balancesResp.Balances, nil
}

func queryShentuDenomBalance(endpoint, addr, denom string) (sdk.Coin, error) {
	var zeroCoin sdk.Coin

	path := fmt.Sprintf(
		"%s/cosmos/bank/v1beta1/balances/%s/by_denom?denom=%s",
		endpoint, addr, denom,
	)
	resp, err := http.Get(path)
	if err != nil {
		return zeroCoin, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return zeroCoin, err
	}

	var balanceResp banktypes.QueryBalanceResponse
	if err := cdc.UnmarshalJSON(bz, &balanceResp); err != nil {
		return zeroCoin, err
	}

	return *balanceResp.Balance, nil
}

func queryDelegation(endpoint, validatorAddr, delegatorAddr string) (stakingtypes.QueryDelegationResponse, error) {
	var res stakingtypes.QueryDelegationResponse

	path := fmt.Sprintf(
		"%s/cosmos/staking/v1beta1/validators/%s/delegations/%s",
		endpoint, validatorAddr, delegatorAddr,
	)

	resp, err := http.Get(path)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(bz, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryProposal(endpoint string, proposalId int) (govtypes.QueryProposalResponse, error) {
	var res govtypes.QueryProposalResponse
	path := fmt.Sprintf(
		"%s/shentu/gov/v1alpha1/proposals/%d",
		endpoint, proposalId,
	)

	resp, err := http.Get(path)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(bz, &res); err != nil {
		return res, err
	}

	return res, nil
}

func queryCertificate(endpoint string, certificateId int) (certtypes.QueryCertificateResponse, error) {
	var res certtypes.QueryCertificateResponse
	path := fmt.Sprintf(
		"%s/shentu/cert/v1alpha1/certificate/%d",
		endpoint, certificateId,
	)

	resp, err := http.Get(path)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(bz, &res); err != nil {
		return res, err
	}

	return res, nil
}

func queryShieldPool(endpoint string, poolId int) (shieldtypes.QueryPoolResponse, error) {
	var res shieldtypes.QueryPoolResponse
	path := fmt.Sprintf(
		"%s/shentu/shield/v1alpha1/pool/%d",
		endpoint, poolId,
	)

	resp, err := http.Get(path)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(bz, &res); err != nil {
		return res, err
	}

	return res, nil
}

func queryShieldStatus(endpoint string) (shieldtypes.QueryShieldStatusResponse, error) {
	var res shieldtypes.QueryShieldStatusResponse
	path := fmt.Sprintf(
		"%s/shentu/shield/v1alpha1/status",
		endpoint,
	)

	resp, err := http.Get(path)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(bz, &res); err != nil {
		return res, err
	}

	return res, nil
}

func queryShieldPurchase(endpoint, purchaser string, poolId int) (shieldtypes.QueryPurchaseListResponse, error) {
	var res shieldtypes.QueryPurchaseListResponse
	path := fmt.Sprintf(
		"%s/shentu/shield/v1alpha1/purchase_list/%d/%s",
		endpoint, poolId, purchaser,
	)

	resp, err := http.Get(path)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(bz, &res); err != nil {
		return res, err
	}

	return res, nil
}

func queryShieldReimbursement(endpoint string, proposalId int) (shieldtypes.QueryReimbursementResponse, error) {
	var res shieldtypes.QueryReimbursementResponse
	path := fmt.Sprintf(
		"%s/shentu/shield/v1alpha1/proposal/%d/reimbursement",
		endpoint, proposalId,
	)

	resp, err := http.Get(path)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(bz, &res); err != nil {
		return res, err
	}

	return res, nil
}

func configFile(filename string) string {
	return filepath.Join(shentuHome, "config", filename)
}

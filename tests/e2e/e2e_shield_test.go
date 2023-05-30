package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

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

func queryShieldPool(endpoint string, poolId int) (*shieldtypes.QueryPoolResponse, error) {
	grpcReq := &shieldtypes.QueryPoolRequest{
		PoolId: uint64(poolId),
	}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := shieldtypes.NewQueryClient(conn)

	grpcRsp, err := client.Pool(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}

func queryShieldStatus(endpoint string) (*shieldtypes.QueryShieldStatusResponse, error) {
	grpcReq := &shieldtypes.QueryShieldStatusRequest{}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := shieldtypes.NewQueryClient(conn)

	grpcRsp, err := client.ShieldStatus(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}

func queryShieldPurchase(endpoint, purchaser string, poolId int) (*shieldtypes.QueryPurchaseListResponse, error) {
	grpcReq := &shieldtypes.QueryPurchaseListRequest{
		PoolId:    uint64(poolId),
		Purchaser: purchaser,
	}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := shieldtypes.NewQueryClient(conn)

	grpcRsp, err := client.PurchaseList(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}

func queryShieldReimbursement(endpoint string, proposalId int) (*shieldtypes.QueryReimbursementResponse, error) {
	grpcReq := &shieldtypes.QueryReimbursementRequest{
		ProposalId: uint64(proposalId),
	}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := shieldtypes.NewQueryClient(conn)

	grpcRsp, err := client.Reimbursement(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}

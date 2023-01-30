package e2e

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/tempfile"

	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"
	bountycli "github.com/shentufoundation/shentu/v2/x/bounty/client/cli"
	bountytypes "github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (s *IntegrationTestSuite) executeCreateProgram(c *chain, valIdx int, creatorAddr, desc, encKeyFile, commissionRate, deposit, endTime, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty create program on %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"create-program",
		fmt.Sprintf("--%s=%s", bountycli.FlagDesc, desc),
		fmt.Sprintf("--%s=%s", bountycli.FlagEncKeyFile, encKeyFile),
		fmt.Sprintf("--%s=%s", bountycli.FlagCommissionRate, commissionRate),
		fmt.Sprintf("--%s=%s", bountycli.FlagDeposit, deposit),
		fmt.Sprintf("--%s=%s", bountycli.FlagSubmissionEndTime, endTime),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFrom, creatorAddr),
		fmt.Sprintf("--%s=%s", sdkflags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", sdkflags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully create program", creatorAddr)
}

func (s *IntegrationTestSuite) executeSubmitFinding(c *chain, valIdx, programId int, submitAddr, desc, title, poc, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty submit finding on %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"submit-finding",
		fmt.Sprintf("--%s=%s", bountycli.FlagFindingDesc, ""),
		fmt.Sprintf("--%s=%s", bountycli.FlagFindingTitle, ""),
		fmt.Sprintf("--%s=%s", bountycli.FlagFindingPoc, ""),
		fmt.Sprintf("--%s=%d", bountycli.FlagProgramID, programId),
		fmt.Sprintf("--%s=%d", bountycli.FlagFindingSeverityLevel, bountytypes.SeverityLevelMajor),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFrom, submitAddr),
		fmt.Sprintf("--%s=%s", sdkflags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", sdkflags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully submit finding", submitAddr)
}

func (s *IntegrationTestSuite) executeAcceptFinding(c *chain, valIdx, findingId int, hostAddr, comment, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty acctpe finding on %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"accept-finding",
		fmt.Sprintf("%d", findingId),
		fmt.Sprintf("--%s=%s", bountycli.FlagComment, comment),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFrom, hostAddr),
		fmt.Sprintf("--%s=%s", sdkflags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", sdkflags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully accept finding", hostAddr)
}

func (s *IntegrationTestSuite) executeRejectFinding(c *chain, valIdx, findingId int, hostAddr, comment, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty reject finding on %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"reject-finding",
		fmt.Sprintf("%d", findingId),
		fmt.Sprintf("--%s=%s", bountycli.FlagComment, comment),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFrom, hostAddr),
		fmt.Sprintf("--%s=%s", sdkflags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", sdkflags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFees, fees),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully reject finding", hostAddr)
}

func queryBountyProgram(endpoint string, programID int) (*bountytypes.QueryProgramResponse, error) {
	grpcReq := &bountytypes.QueryProgramRequest{
		ProgramId: uint64(programID),
	}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := bountytypes.NewQueryClient(conn)

	grpcRsp, err := client.Program(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}

func queryBountyFinding(endpoint string, findingID int) (*bountytypes.QueryFindingResponse, error) {
	grpcReq := &bountytypes.QueryFindingRequest{
		FindingId: uint64(findingID),
	}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := bountytypes.NewQueryClient(conn)

	grpcRsp, err := client.Finding(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}

func generateBountyKeyFile(filePath string) error {
	priKey, _, err := bountycli.GenerateKey()
	if err != nil {
		return err
	}
	if filePath == "" {
		return fmt.Errorf("empty key file path")
	}

	decKeyBz := crypto.FromECDSA(priKey.ExportECDSA())
	if err := tempfile.WriteFileAtomic(filePath, decKeyBz, 0666); err != nil {
		return err
	}
	return nil
}

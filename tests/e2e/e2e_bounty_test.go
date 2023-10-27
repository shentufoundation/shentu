package e2e

import (
	"context"
	"fmt"
	"strings"
	"time"

	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"
	bountycli "github.com/shentufoundation/shentu/v2/x/bounty/client/cli"
	bountytypes "github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (s *IntegrationTestSuite) executeCreateProgram(c *chain, valIdx int, pid, name, desc, creatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty open program %s on %s", pid, c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"create-program",
		fmt.Sprintf("--%s=%s", bountycli.FlagProgramID, pid),
		fmt.Sprintf("--%s=%s", bountycli.FlagName, name),
		fmt.Sprintf("--%s=%s", bountycli.FlagDesc, desc),
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

func (s *IntegrationTestSuite) executeOpenProgram(c *chain, valIdx int, pid, creatorAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty create program %s on %s", pid, c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"open-program",
		fmt.Sprintf("%s", pid),
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

func (s *IntegrationTestSuite) executeSubmitFinding(c *chain, valIdx int, pid, fid, submitAddr, title, desc, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty submit finding on %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"submit-finding",
		fmt.Sprintf("--%s=%s", bountycli.FlagProgramID, pid),
		fmt.Sprintf("--%s=%s", bountycli.FlagFindingID, fid),
		fmt.Sprintf("--%s=%s", bountycli.FlagFindingTitle, title),
		fmt.Sprintf("--%s=%s", bountycli.FlagDesc, desc),
		fmt.Sprintf("--%s=%d", bountycli.FlagFindingSeverityLevel, bountytypes.SeverityLevelMedium),
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

func (s *IntegrationTestSuite) executeAcceptFinding(c *chain, valIdx int, findingId, hostAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty acctpe finding on %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"accept-finding",
		fmt.Sprintf("%s", findingId),
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

func (s *IntegrationTestSuite) executeRejectFinding(c *chain, valIdx int, findingId, hostAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty reject finding on %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"reject-finding",
		fmt.Sprintf("%s", findingId),
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

//func (s *IntegrationTestSuite) executeReleaseFinding(c *chain, valIdx, findingId int, hostAddr, keyFile, fees string) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
//	defer cancel()
//
//	s.T().Logf("Executing shentu bounty release finding on %s", c.id)
//
//	command := []string{
//		shentuBinary,
//		txCommand,
//		bountytypes.ModuleName,
//		"release-finding",
//		fmt.Sprintf("%d", findingId),
//		fmt.Sprintf("--%s=%s", bountycli.FlagEncKeyFile, keyFile),
//		fmt.Sprintf("--%s=%s", sdkflags.FlagFrom, hostAddr),
//		fmt.Sprintf("--%s=%s", sdkflags.FlagChainID, c.id),
//		fmt.Sprintf("--%s=%s", sdkflags.FlagGas, "auto"),
//		fmt.Sprintf("--%s=%s", sdkflags.FlagFees, fees),
//		"--keyring-backend=test",
//		"--output=json",
//		"-y",
//	}
//
//	s.T().Logf("cmd: %s", strings.Join(command, " "))
//
//	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
//	s.T().Logf("%s successfully release finding", hostAddr)
//}

func (s *IntegrationTestSuite) executeEndProgram(c *chain, valIdx int, programId, hostAddr, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu bounty close program on %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		bountytypes.ModuleName,
		"close-program",
		fmt.Sprintf("%s", programId),
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
	s.T().Logf("%s successfully end program", hostAddr)
}

func queryBountyProgram(endpoint, programID string) (*bountytypes.QueryProgramResponse, error) {
	grpcReq := &bountytypes.QueryProgramRequest{
		ProgramId: programID,
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

func queryBountyFinding(endpoint, findingID string) (*bountytypes.QueryFindingResponse, error) {
	grpcReq := &bountytypes.QueryFindingRequest{
		FindingId: findingID,
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

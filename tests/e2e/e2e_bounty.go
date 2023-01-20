package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

	s.T().Logf("Executing shentu tx issue certificate %s", c.id)

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

func queryBountyProgram(endpoint string, programID int) (bountytypes.QueryProgramResponse, error) {
	var programResp bountytypes.QueryProgramResponse

	path := fmt.Sprintf(
		"%s/shentu/bounty/v1/programs/%d",
		endpoint, programID,
	)
	resp, err := http.Get(path)
	if err != nil {
		return programResp, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return programResp, err
	}

	if err := cdc.UnmarshalJSON(bz, &programResp); err != nil {
		return programResp, err
	}

	return programResp, nil
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

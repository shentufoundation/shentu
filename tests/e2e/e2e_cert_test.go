package e2e

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

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

func queryCertificate(endpoint string, certificateId int) (*certtypes.QueryCertificateResponse, error) {
	grpcReq := &certtypes.QueryCertificateRequest{
		CertificateId: uint64(certificateId),
	}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := certtypes.NewQueryClient(conn)

	grpcRsp, err := client.Certificate(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}

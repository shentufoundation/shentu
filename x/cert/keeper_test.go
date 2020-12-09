package cert_test

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

func randomString(length int) string {
	out := make([]byte, length)
	rand.Read(out)
	return string(out)
}

func Test_IssueAndRevokeCertificate(t *testing.T) {
	t.Run("Testing issuing and revoking a certificate", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})
		certifier := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))[0]
		app.CertKeeper.SetCertifier(ctx, types.NewCertifier(certifier, "", certifier, ""))

		cert := types.NewCompilationCertificate(types.CertificateTypeCompilation, 
					"content1", "compiler1", "bytecodehash1", "description", certifier)
		id, err := app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)

		id2, err := app.CertKeeper.IssueCertificate(ctx, types.NewCompilationCertificate(types.CertificateTypeCompilation, 
					"content2", "compiler1", "bytecodehash1", "description", certifier))
		require.NoError(t, err)

		retrieved_cert, err := app.CertKeeper.GetCertificateByID(ctx, id)
		require.NoError(t, err)
		require.True(t, reflect.DeepEqual(cert, retrieved_cert))

		retrieved_ids := app.CertKeeper.GetCertifierCertIDs(ctx, certifier)
		require.True(t, retrieved_ids[0] == id && retrieved_ids[1] == id2)

		contentCertID, found := app.CertKeeper.GetContentCertID(ctx, types.CertificateTypeCompilation, cert.RequestContent())
		require.True(t, found && contentCertID == id)
		
		err = app.CertKeeper.RevokeCertificate(ctx, cert, certifier)
		require.NoError(t, err)

		_, err = app.CertKeeper.GetCertificateByID(ctx, id)
		require.Error(t, err)

		retrieved_ids = app.CertKeeper.GetCertifierCertIDs(ctx, certifier)
		require.True(t, retrieved_ids[0] == id2)

		_, found = app.CertKeeper.GetContentCertID(ctx, types.CertificateTypeCompilation, cert.RequestContent())
		require.True(t, !found)
	})
}

func Test_IterationByCertifier(t *testing.T) {
	t.Run("Testing certifier-based iteration", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})
		addrs := simapp.AddTestAddrs(app, ctx, 5, sdk.NewInt(10000))
		for _, addr := range addrs {
			app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, "", addr, ""))
		}

		// Store certificates
		addr0Count := 0
		addr2Count := 0
		for i := 1; i < 500; i++ {
			index := rand.Intn(4)
			if index == 0 {
				addr0Count++
			} else if index == 2 {
				addr2Count++
			}
			length := rand.Intn(10) + 10
			s := randomString(length)
			cert := types.NewCompilationCertificate(types.CertificateTypeCompilation, s, "compiler1", "bytecodehash1", "", addrs[index])
			_, err := app.CertKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err)
		}

		certs := app.CertKeeper.GetCertificatesByCertifier(ctx, addrs[0])
		require.Equal(t, addr0Count, len(certs))

		certs = app.CertKeeper.GetCertificatesByCertifier(ctx, addrs[2])
		require.Equal(t, addr2Count, len(certs))
	})
}

func Test_CertificateQueries(t *testing.T) {
	t.Run("Testing various queries on certifications", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})
		addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))
		for _, addr := range addrs {
			app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, "", addr, ""))
		}

		var cert *types.CompilationCertificate
		content1, content2  := "duplicate content", "duplicate content 2"
		
		count1 := 0 // number of certificates with duplicate content
		count2 := 0 // number of certificates certified by addr[0]
		
		for _, certType := range types.CertificateTypes {
			index := rand.Intn(2) // random address index
			if index == 1 { // 50% - 50% between two contents
				cert = types.NewCompilationCertificate(certType, content1, "compiler1", "bytecodehash1", "", addrs[index])
				count1++
			} else {
				cert = types.NewCompilationCertificate(certType, content2, "compiler1", "bytecodehash1", "", addrs[index])
			}
			if index == 0 { // addr[0] certifier
				count2++
			}
			_, err := app.CertKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err)
		}
				
		// Queries by request content
		contentToFind, _ := types.NewRequestContent("sourcecodehash", content1)
		certs := app.CertKeeper.GetCertificatesByContent(ctx, contentToFind)
		require.Equal(t, count1, len(certs))

		queryParams := types.NewQueryCertificatesParams(1, 100, nil, "sourcecodehash", content1)
		_, certs_filtered, err := app.CertKeeper.GetCertificatesFiltered(ctx, queryParams)
		require.NoError(t, err)
		reflect.DeepEqual(certs, certs_filtered)

		// Queries by certifier
		queryParams = types.NewQueryCertificatesParams(1, 100, addrs[0], "", "")
		_, certs, err = app.CertKeeper.GetCertificatesFiltered(ctx, queryParams)
		require.NoError(t, err)
		require.Equal(t, count2, len(certs))
	})
}

func Test_IsCertified(t *testing.T) {
	t.Run("Testing the function IsCertified", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})
		addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))
		app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "", addrs[0], ""))

		certType := "auditing"
		contentTypeStr := "address"
		contentStr := "certik1k4gj07sgy6x3k6ms31aztgu9aajjkaw3ktsydag"

		isCertified := app.CertKeeper.IsCertified(ctx, contentTypeStr, contentStr, certType)
		require.Equal(t, false, isCertified)

		cert, err := types.NewGeneralCertificate(certType, contentTypeStr, contentStr,
			"Audited by CertiK", addrs[0])
		require.NoError(t, err)

		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)

		isCertified = app.CertKeeper.IsCertified(ctx, contentTypeStr, contentStr, certType)
		require.Equal(t, true, isCertified)
	})
}

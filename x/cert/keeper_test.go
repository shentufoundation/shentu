package cert_test

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cert/types"
)

func Test_GetNewCertificateID(t *testing.T) {
	t.Run("Testing GetNewCertificateID", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
		addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

		// Set and Get a certificate
		c1 := types.NewCompilationCertificate(types.CertificateTypeCompilation, "sourcodehash0",
			"compiler1", "bytecodehash1", "", addrs[0])

		id1 := app.CertKeeper.GetNextCertificateID(ctx)

		c1.SetCertificateID(id1)
		app.CertKeeper.SetNextCertificateID(ctx, id1+1)
		app.CertKeeper.SetCertificate(ctx, c1)

		data, _ := app.CertKeeper.GetCertificateByID(ctx, id1)
		if data == nil {
			t.Errorf("Could not retrieve data from the store")
		}
		if !reflect.DeepEqual(data, c1) {
			t.Errorf("Retrieved data different from the original data")
		}

		// Set an identical certificate
		c2 := types.NewCompilationCertificate(types.CertificateTypeCompilation, "sourcodehash0",
			"compiler1", "bytecodehash1", "", addrs[0])
		id2 := app.CertKeeper.GetNextCertificateID(ctx)
		c2.SetCertificateID(id2)
		app.CertKeeper.SetNextCertificateID(ctx, id2+1)
		app.CertKeeper.SetCertificate(ctx, c2)

		data, _ = app.CertKeeper.GetCertificateByID(ctx, id2)
		if data == nil {
			t.Errorf("Could not retrieve data from the store")
		}
		if !reflect.DeepEqual(data, c2) {
			t.Errorf("Retrieved data different from the original data")
		}

		// Delete the first certificate and add the third certificate
		id := c1.ID()
		app.CertKeeper.DeleteCertificate(ctx, c1)

		c3 := types.NewCompilationCertificate(types.CertificateTypeCompilation, "sourcodehash0",
			"compiler1", "bytecodehash1", "", addrs[0])
		id3 := app.CertKeeper.GetNextCertificateID(ctx)
		require.Equal(t, id+2, id3)

		c3.SetCertificateID(id3)
		app.CertKeeper.SetCertificate(ctx, c3)

		data, _ = app.CertKeeper.GetCertificateByID(ctx, id3)
		if data == nil {
			t.Errorf("Could not retrieve data from the store")
		}
		if !reflect.DeepEqual(data, c3) {
			t.Errorf("Retrieved data different from the original data")
		}
	})
}

func randomString(length int) string {
	out := make([]byte, length)
	rand.Read(out)
	return string(out)
}

func Test_IterationByCertifier(t *testing.T) {
	t.Run("Testing certifier-based iteration", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
		addrs := simapp.AddTestAddrs(app, ctx, 5, sdk.NewInt(10000))
		for _, addr := range addrs {
			app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, "", addr, ""))
		}

		// Store certificates
		addr0Count := 0
		addr2Count := 0
		for i := 1; i < 50000; i++ {
			index := rand.Intn(4)
			if index == 0 {
				addr0Count++
			} else if index == 2 {
				addr2Count++
			}
			length := rand.Intn(10) + 10
			s := randomString(length)
			cert := types.NewCompilationCertificate(types.CertificateTypeCompilation, s,
				"compiler1", "bytecodehash1", "", addrs[index])
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
		ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
		addrs := simapp.AddTestAddrs(app, ctx, 5, sdk.NewInt(10000))
		for _, addr := range addrs {
			app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, "", addr, ""))
		}

		// Store certificates
		count := 0
		count2 := 0 // For counting certificates with given address and content
		count3 := 0 // For counting certificates with given address
		dupContent := "duplicate content"
		totalCerts := 1000
		for i := 1; i < totalCerts; i++ {
			index := rand.Intn(4) // random address index

			var cert *types.CompilationCertificate
			dup := rand.Intn(100)
			if dup > 95 {
				cert = types.NewCompilationCertificate(types.CertificateTypeCompilation, dupContent,
					"compiler1", "bytecodehash1", "", addrs[index])
				count++
				if index == 0 {
					count2++
					count3++
				}
			} else {
				length := rand.Intn(10) + 10
				s := randomString(length)
				cert = types.NewCompilationCertificate(types.CertificateTypeCompilation, s, "compiler1",
					"bytecodehash1", "", addrs[index])
				if index == 0 {
					count3++
				}
			}
			_, err := app.CertKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err)
		}

		// Test GetCertificatesByContent()
		contentToFind, _ := types.NewRequestContent("sourcecodehash", dupContent)
		certs := app.CertKeeper.GetCertificatesByContent(ctx, contentToFind)
		require.Equal(t, count, len(certs))

		// Test GetCertificatesFiltered()
		// Query by content only
		queryParams := types.NewQueryCertificatesParams(1, totalCerts, nil, "sourcecodehash", dupContent)
		total, certs, err := app.CertKeeper.GetCertificatesFiltered(ctx, queryParams)
		require.NoError(t, err)
		require.Equal(t, uint64(count), total)
		require.Equal(t, count, len(certs))

		// Query by content and certifier
		queryParams = types.NewQueryCertificatesParams(1, totalCerts, addrs[0], "sourcecodehash", dupContent)
		total, certs, err = app.CertKeeper.GetCertificatesFiltered(ctx, queryParams)
		require.NoError(t, err)
		require.Equal(t, uint64(count2), total)
		require.Equal(t, count2, len(certs))

		// Query by certifier only
		queryParams = types.NewQueryCertificatesParams(1, totalCerts, addrs[0], "", "")
		total, certs, err = app.CertKeeper.GetCertificatesFiltered(ctx, queryParams)
		require.NoError(t, err)
		require.Equal(t, uint64(count3), total)
		require.Equal(t, count3, len(certs))
	})
}

func Test_IsCertified(t *testing.T) {
	t.Run("Testing the function IsCertified", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
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

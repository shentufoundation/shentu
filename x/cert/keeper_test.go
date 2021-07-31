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
		c1, err := types.NewCertificate("compilation", "sourcodehash0", "compiler1",
			"bytecodehash1", "", addrs[0])
		require.NoError(t, err)

		id1 := app.CertKeeper.GetNextCertificateID(ctx)

		c1.CertificateId = id1
		app.CertKeeper.SetNextCertificateID(ctx, id1+1)
		app.CertKeeper.SetCertificate(ctx, c1)

		data, err := app.CertKeeper.GetCertificateByID(ctx, id1)
		if err != nil {
			t.Errorf("Could not retrieve data from the store")
		}
		if !reflect.DeepEqual(data, c1) {
			t.Errorf("Retrieved data different from the original data")
		}

		// Set an identical certificate
		c2, err := types.NewCertificate("compilation", "sourcodehash0", "compiler1",
			"bytecodehash1", "", addrs[0])
		require.NoError(t, err)
		id2 := app.CertKeeper.GetNextCertificateID(ctx)
		c2.CertificateId = id2
		app.CertKeeper.SetNextCertificateID(ctx, id2+1)
		app.CertKeeper.SetCertificate(ctx, c2)

		data, err = app.CertKeeper.GetCertificateByID(ctx, id2)
		if err != nil {
			t.Errorf("Could not retrieve data from the store")
		}
		if !reflect.DeepEqual(data, c2) {
			t.Errorf("Retrieved data different from the original data")
		}

		// Delete the first certificate and add the third certificate
		id := c1.CertificateId
		app.CertKeeper.DeleteCertificate(ctx, c1)

		c3, err := types.NewCertificate("compilation", "sourcodehash0", "compiler1",
			"bytecodehash1", "", addrs[0])
		require.NoError(t, err)
		id3 := app.CertKeeper.GetNextCertificateID(ctx)
		require.Equal(t, id+2, id3)

		c3.CertificateId = id3
		app.CertKeeper.SetCertificate(ctx, c3)

		data, err = app.CertKeeper.GetCertificateByID(ctx, id3)
		if err != nil {
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
			cert, err := types.NewCertificate("compilation", s, "compiler1",
				"bytecodehash1", "", addrs[index])
			require.NoError(t, err)
			_, err = app.CertKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err)
		}

		certs := app.CertKeeper.GetCertificatesByCertifier(ctx, addrs[0])
		require.Equal(t, addr0Count, len(certs))

		certs = app.CertKeeper.GetCertificatesByCertifier(ctx, addrs[2])
		require.Equal(t, addr2Count, len(certs))
	})
}

func Test_IsCertified(t *testing.T) {
	t.Run("Testing the function IsCertified", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
		addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))
		app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "", addrs[0], ""))

		certType := "auditing"
		contentStr := "certik1k4gj07sgy6x3k6ms31aztgu9aajjkaw3ktsydag"

		isCertified := app.CertKeeper.IsCertified(ctx, contentStr, certType)
		require.Equal(t, false, isCertified)

		cert, err := types.NewCertificate(certType, contentStr,
			"", "", "Audited by CertiK", addrs[0])
		require.NoError(t, err)

		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)

		isCertified = app.CertKeeper.IsCertified(ctx, contentStr, certType)
		require.Equal(t, true, isCertified)
	})
}

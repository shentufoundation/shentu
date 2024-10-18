package keeper_test

//
//import (
//	"github.com/stretchr/testify/require"
//	"testing"
//	"time"
//
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//
//	shentuapp "github.com/shentufoundation/shentu/v2/app"
//	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
//	"github.com/shentufoundation/shentu/v2/x/cert/types"
//)
//
//func TestMsgServer_Certificate(t *testing.T) {
//	ctx, certKeeper, msgServer, addrs := DoInit(t)
//
//	content := types.AssembleContent("identity", addrs[1].String())
//	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), types.NewMsgIssueCertificate(content, "", "", "", addrs[0]))
//	require.NoError(t, err)
//
//	certificateID := uint64(1)
//	_, err = msgServer.RevokeCertificate(sdk.WrapSDKContext(ctx), types.NewMsgRevokeCertificate(addrs[0], certificateID, ""))
//	require.NoError(t, err)
//
//	_, err = certKeeper.GetCertificateByID(ctx, certificateID)
//	require.Error(t, err)
//}
//
//func DoInit(t *testing.T) (sdk.Context, keeper.Keeper, types.MsgServer, []sdk.AccAddress) {
//	app := shentuapp.Setup(t, false)
//	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
//	addrs := shentuapp.AddTestAddrs(app, ctx, 5, sdk.NewInt(80000*1e6))
//	ok := app.CertKeeper
//	msgServer := keeper.NewMsgServerImpl(ok)
//	ctx = ctx.WithValue("msgServer", msgServer).WithValue("t", t).WithValue("ok", ok)
//
//	ok.SetCertifier(ctx, types.Certifier{
//		Address: addrs[0].String(),
//	})
//	return ctx, ok, msgServer, addrs
//}

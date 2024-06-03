package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	. "github.com/shentufoundation/shentu/v2/x/cvm/keeper"
)

func TestNewBlockChain(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cvmk := app.CVMKeeper

	bc := NewBlockChain(ctx, cvmk)
	require.NotNil(t, bc.LastBlockTime())
}

func TestBlockchain_BlockHash(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cvmk := app.CVMKeeper

	bc := NewBlockChain(ctx, cvmk)
	ctxHash := ctx.BlockHeader().LastBlockId.Hash
	bz, err := bc.BlockHash(0)
	require.Nil(t, err)
	require.Equal(t, ctxHash, bz)

	bz, err = bc.BlockHash(10)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestBlockchain_LastBlockHeight(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cvmk := app.CVMKeeper
	bc := NewBlockChain(ctx, cvmk)

	ctxHeight := ctx.BlockHeader().Height
	require.Equal(t, uint64(ctxHeight), bc.LastBlockHeight())
}

func TestBlockchain_LastBlockTime(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cvmk := app.CVMKeeper
	bc := NewBlockChain(ctx, cvmk)

	ctxTime := ctx.BlockHeader().Time
	require.Equal(t, ctxTime, bc.LastBlockTime())
}

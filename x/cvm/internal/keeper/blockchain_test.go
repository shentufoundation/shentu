package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cvm/internal/keeper"
)

func NewGasMeter(limit uint64) sdk.GasMeter {
	return sdk.NewGasMeter(limit)
}

func TestNewBlockChain(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	bc := keeper.NewBlockChain(ctx, app.CvmKeeper)
	require.NotNil(t, bc.LastBlockTime())
}

func TestBlockchain_BlockHash(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	bc := keeper.NewBlockChain(ctx, app.CvmKeeper)

	ctxHash := ctx.BlockHeader().LastBlockId.Hash
	bz, err := bc.BlockHash(0)
	require.Nil(t, err)
	require.Equal(t, ctxHash, bz)

	bz, err = bc.BlockHash(10)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestBlockchain_LastBlockHeight(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	bc := keeper.NewBlockChain(ctx, app.CvmKeeper)

	ctxHeight := ctx.BlockHeader().Height
	require.Equal(t, uint64(ctxHeight), bc.LastBlockHeight())
}

func TestBlockchain_LastBlockTime(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	bc := keeper.NewBlockChain(ctx, app.CvmKeeper)

	ctxTime := ctx.BlockHeader().Time
	require.Equal(t, ctxTime, bc.LastBlockTime())
}

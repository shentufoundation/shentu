package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"reflect"
	"encoding/hex"

	"github.com/spf13/cobra"

	"github.com/gogo/protobuf/proto"
	"github.com/cosmos/cosmos-sdk/client"
	codec "github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"

	sdk "github.com/cosmos/cosmos-sdk/types"
	
	shentuapp "github.com/shentufoundation/shentu/v2/app"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func PubkeyToAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pubkey2addr <pubkey_json>",
		Short: "Convert pubkey to the address",
		Long: `sample: 
shentud pubkey2addr '{"@type": "/cosmos.crypto.ed25519.PubKey", "key": "lzn2Q4Z4DfUdrIDdxVOcXQS44gEdlLpL8T3QnO4brZk="}'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var pk cryptotypes.PubKey
			if err := cliCtx.Codec.UnmarshalInterfaceJSON([]byte(args[0]), &pk); err != nil {
				return err
			}
			sdk.AccAddress(pk.Address()).String()
			cliCtx.PrintString(
				sdk.AccAddress(pk.Address()).String() + "\n" +
					sdk.ValAddress(pk.Address()).String() + "\n" +
					sdk.ConsAddress(pk.Address()).String() + "\n",
			)
			return nil
		},
	}
	return cmd
}

func CheckStoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-store <rep-pattern>",
		Short: "match dumped store with the given regular expression pattern",
		Long: `sample (matching bech32 address prefix [shentu]):
shentud check-store 'shentu[a-z]{0,10}1'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			serverCtx := server.GetServerContextFromCmd(cmd)
			home := serverCtx.Config.RootDir
			db, err := sdk.NewLevelDB("application", filepath.Join(home, "data"))
			if err != nil {
				return err
			}
			app := newApp(serverCtx.Logger, db, nil, serverCtx.Viper)
			shentuApp, _ := app.(*shentuapp.ShentuApp)
			ctx := sdk.NewContext(shentuApp.CommitMultiStore().CacheMultiStore(),
				tmproto.Header{Height: shentuApp.LastBlockHeight()},
				false,
				shentuApp.Logger(),
			)
			jstr := checkKeys(ctx, shentuApp, cliCtx)

			cliCtx.PrintString("------------------- " + jstr)
			return nil
		},
	}
	return cmd
}

type So struct { 
	strings.Builder 
	wrt io.Writer
}

func (so *So) WriteStarter(str, sep string) {
	txt := "\n\""+str+"\":"+sep
	if so.wrt != nil {
		so.wrt.Write([]byte(txt))
	} else {
		so.WriteString(txt)
	}
}

func (so *So) WriteEnder(edr string) {
	txt := "\n"+edr
	if so.wrt != nil {
		so.wrt.Write([]byte(txt))
	} else {
		so.WriteString(txt)
	}
}

func (so *So) WriteString(txt string) {
	if so.wrt != nil {
		so.wrt.Write([]byte(txt))
	} else {
		so.Builder.WriteString(txt)
	}
}

func checkKeys(ctx sdk.Context, app *shentuapp.ShentuApp, cliCtx client.Context) string {
	cdc := app.Codec()
	var so = So{wrt: os.Stdout}
	for skn, ks := range allKeys {
		store := ctx.KVStore(app.GetKey(skn))
		so.WriteStarter(skn+"###", "{")
		for _, k := range ks {
			iter := sdk.KVStorePrefixIterator(store, k.prefix)
			so.WriteStarter(hex.EncodeToString(k.prefix), "[")
			for ; iter.Valid(); iter.Next() {
				var msg proto.Message
				if k.marshalWay != 3 {
					iv, ok := k.ptr.(codec.ProtoMarshaler)
					if !ok {
						panic("failed to cast to codec.ProtoMarshaler")
					}
					if k.marshalWay == 2 {
						cdc.MustUnmarshalLengthPrefixed(iter.Value(), iv)
					} else if k.marshalWay == 1{
						cdc.MustUnmarshal(iter.Value(), iv)
					}
					msg = iv.(proto.Message)
				} else if k.marshalWay == 3 {
					cdc.UnmarshalInterface(iter.Value(), k.ptr)
					msg = reflect.ValueOf(k.ptr).Elem().Interface().(proto.Message)
				}
				vstr := string(cdc.MustMarshalJSON(msg))
				so.WriteString(vstr+",")
			}
			so.WriteEnder("]")
			iter.Close()
		}
		so.WriteEnder("}")
	}
	return so.String()
}

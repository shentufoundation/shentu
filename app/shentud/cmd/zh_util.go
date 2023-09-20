package cmd

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
	"strconv"
	"sort"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	codec "github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/gogo/protobuf/proto"
	gogotypes "github.com/gogo/protobuf/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

const (
	ModulesFlag = "modules"
	ModulesFlagP = "m"
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
	ks := make([]string, 0, len(allKeys))
	for k := range allKeys {
		ks = append(ks, k)
	}
	modulesTxt := strings.Join(ks, ",")
	cmd := &cobra.Command{
		Use:   "check-store <rep-pattern>",
		Short: "match dumped store with the given regular expression pattern",
		Long: `sample (matching bech32 address prefix [shentu]):
shentud check-store 'shentu[a-z]{0,10}1'
All modules are: ` + modulesTxt,
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			so := So{wrt: os.Stdout}
			if len(args) == 1 {
				so.rep, err = regexp.Compile(args[0])
				if err != nil {
					return err
				}
			}
			ms, _ := cmd.Flags().GetString(ModulesFlag)
			var ks map[string]bool
			if ms != "" {
				ks = make(map[string]bool)
				for _, k := range strings.Split(ms, ",") {
					ks[k] = true
				}
			}
			for k := range ks {
				if _, ok := allKeys[k]; !ok {
					return fmt.Errorf("module %s not in the list", k)
				}
			}
			
			serverCtx := server.GetServerContextFromCmd(cmd)
			home := serverCtx.Config.RootDir
			db, err := sdk.NewLevelDB("application", filepath.Join(home, "data"))
			if err != nil {
				return err
			}
			app := newApp(serverCtx.Logger, db, nil, serverCtx.Viper)
			shentuApp, _ := app.(*shentuapp.ShentuApp)
			ctx := sdk.NewContext(
				shentuApp.CommitMultiStore().CacheMultiStore(),
				tmproto.Header{Height: shentuApp.LastBlockHeight()},
				false,
				shentuApp.Logger(),
			)

			jstr := checkKeys(ctx, shentuApp, cliCtx, so, ks)
			// printIBCKeys(ctx, shentuApp, cliCtx, so)
			cliCtx.PrintString(jstr + "\n")
			return nil
		},
	}
	cmd.Flags().StringP(ModulesFlag, ModulesFlagP, "", "specify comma separated module names that will be checked")
	return cmd
}
type Layer struct {
	encloseToken string
	starting bool
}

type So struct { 
	strings.Builder 
	wrt io.Writer
	rep *regexp.Regexp
	layers []Layer
}

func (so *So) WriteStarter(str, sep string) {
	txt := sep
	if str != "" {
		txt = "\n\""+str+"\":"+txt
	}
	if !so.HitFirst() {
		txt = ", " + txt
	}
	so.WriteString(txt)
}

func (so *So) StartObj(key string) {
	so.WriteStarter(key, "{")
	so.layers = append(so.layers, Layer{"}", true})
}

func (so *So) StartArray(key string) {
	so.WriteStarter(key, "[")
	so.layers = append(so.layers, Layer{"]", true})
}

func (so *So) End() {
	if len(so.layers) <= 0 {
		panic("empty tokens when calling End")
	}
	innermostLayer := so.layers[len(so.layers)-1]
	so.layers  = so.layers[:len(so.layers)-1]
	so.WriteString(innermostLayer.encloseToken)
}

func (so *So) HitFirst() bool {
	if len(so.layers) <= 0 {
		return true
	}
	isFirst := so.layers[len(so.layers)-1].starting
	so.layers[len(so.layers)-1].starting = false
	return isFirst
}

func (so *So) WriteKV(key, value []byte) error {
	if so.rep != nil {
		if so.rep.FindIndex(value) != nil {
			so.WriteString("found the pattern! sample: ", string(value))
			return fmt.Errorf("pattern found")
		}
		return nil
	}

	wrappedKey := "{\"k\":"+binToStr(key)+", \"v\":"
	if !so.HitFirst() {
		wrappedKey = ", " + wrappedKey
	}
	if so.wrt != nil {
		so.wrt.Write([]byte(wrappedKey))
	} else {
		so.Builder.WriteString(wrappedKey)
	}
	so.WriteString(string(value), "}")
	return nil
}

func (so *So) WriteString(strs ...string) {
	if so.wrt != nil {
		so.wrt.Write([]byte(strings.Join(strs, "")))
	} else {
		so.Builder.WriteString(strings.Join(strs, ""))
	}
}

func clear(itr interface{}) {
	if itr != nil {
		p := reflect.ValueOf(itr).Elem()
		p.Set(reflect.Zero(p.Type()))
	}
}

func byteToStr(bys []byte) string {
	if utf8.Valid(bys) {
		return string(bys)
	}
	if sdk.VerifyAddressFormat(bys) != nil {
		return "addr-"+hex.EncodeToString(bys)
	}
	return hex.EncodeToString(bys)
}

func checkKeys(ctx sdk.Context, app *shentuapp.ShentuApp, cliCtx client.Context, so So, mKeys map[string]bool) string {
	cdc := app.Codec()
	allModules := make([]string, 0, len(allKeys))
	for k := range allKeys {
		allModules = append(allModules, k)
	}
	sort.Strings(allModules)

	so.StartObj("")
	for _, skn := range allModules {
		if mKeys != nil && !mKeys[skn] {
			continue
		}
		store := ctx.KVStore(app.GetKey(skn))
		so.StartObj("module-"+skn)
		for _, k := range allKeys[skn] {
			iter := sdk.KVStorePrefixIterator(store, k.prefix)
			so.StartArray("key-"+hex.EncodeToString(k.prefix))
			for ; iter.Valid(); iter.Next() {
				var msg proto.Message
				if k.marshalWay == 3 {
					cdc.UnmarshalInterface(iter.Value(), k.ptr)
					msg = reflect.ValueOf(k.ptr).Elem().Interface().(proto.Message)
				} else if k.marshalWay == 4 {
					msg = &gogotypes.StringValue{Value: byteToStr(iter.Value())}
				} else {
					iv, ok := k.ptr.(codec.ProtoMarshaler)
					if !ok {
						panic("failed to cast to codec.ProtoMarshaler")
					}
					if k.marshalWay == 2 {
						cdc.MustUnmarshalLengthPrefixed(iter.Value(), iv)
					} else if k.marshalWay == 1{
						cdc.MustUnmarshal(iter.Value(), iv)
					} else {
						panic("unknow marshalway!")
					}
					msg = iv.(proto.Message)
				}
				if so.WriteKV(iter.Key(), cdc.MustMarshalJSON(msg)) != nil {
					clear(k.ptr)
					break
				}
				clear(k.ptr)
			}
			so.End()
			iter.Close()
		}
		so.End()
	}
	so.End()
	return so.String()
}

//for printable ascii, print as is
//otherwise, print [hex]
//escape double-quote character
func binToStr(binData []byte) string {
	sb := strings.Builder{}
	for _, b := range binData {
		if b >= 32 && b <= 126 {
			sb.WriteString(fmt.Sprintf("%c", b))
		} else {
			sb.WriteString(fmt.Sprintf("[%02x]", b))
		}
	}
	return strconv.Quote(sb.String()) 
}

func printIBCKeys(ctx sdk.Context, app *shentuapp.ShentuApp, cliCtx client.Context, so So) string {
	// cdc := app.Codec()
	ibcKeys := []string{
		ibchost.StoreKey,
		ibctransfertypes.StoreKey,
		icahosttypes.StoreKey,
	}
	for _, skn := range ibcKeys {
		store := ctx.KVStore(app.GetKey(skn))
		iter := sdk.KVStorePrefixIterator(store, nil)
		so.StartObj(skn)
		for ; iter.Valid(); iter.Next() {
			so.WriteString("\n", binToStr(iter.Key()))
		}
		so.End()
		iter.Close()
	}
	return "//////////////"
}

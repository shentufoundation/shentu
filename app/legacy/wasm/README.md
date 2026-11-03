# Historical wasm parameter messages

This package preserves `cosmwasm.wasm.v1.MsgUpdateParams`, `Params`,
`AccessConfig`, and `AccessType` from wasmd v0.55.1 (Apache-2.0), the version
removed in Shentu v2.19.0. It supports decoding historical governance proposals
and exporting/importing their genesis JSON without restoring the wasm module.
Only this known historical message is registered; other wasm messages are not
silently accepted or discarded. Add their original schemas if needed.

The source schemas are in `proto/cosmwasm/wasm/v1`. Wire field numbers, protobuf
names, enum values, and JSON field names are unchanged. The Go package is local;
stringers are generated instead of depending on wasm's handwritten helpers.
No message service, keeper, or execution handler is registered.

Regenerate from the repository root using Cosmos `protoc-gen-gocosmos`:

```sh
go build -o /tmp/shentu-protoc-gen-gocosmos github.com/cosmos/gogoproto/protoc-gen-gocosmos
protoc -I proto -I third_party/proto \
  --plugin=protoc-gen-gocosmos=/tmp/shentu-protoc-gen-gocosmos \
  --gocosmos_out=paths=source_relative:/tmp \
  proto/cosmwasm/wasm/v1/types.proto proto/cosmwasm/wasm/v1/tx.proto
cp /tmp/cosmwasm/wasm/v1/{types,tx}.pb.go app/legacy/wasm/
```

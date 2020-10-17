package types

// Client corresponds to the name of the client chain as used in task, e.g. `eth` for task `eth:0x...`.
type Client string

// DefaultClient sets the default client chain name for tasks.
const DefaultClient = "eth"

// Strategy specifies primitive aggregation strategy configuration.
type Strategy struct {
	Type       string      `mapstructure:"type"`
	Primitives []Primitive `mapstructure:"primitive"`
}

// Primitive specifies primitive configuration.
type Primitive struct {
	PrimitiveContractAddr string  `mapstructure:"primitive_contract_address"`
	Weight                float32 `mapstructure:"weight"`
}

// PrimitiveScore groups returned score and related info from primitive.
type PrimitiveScore struct {
	Primitive
	Score uint8
}

// PrimitivePayload specifies primitive query payload.
type PrimitivePayload struct {
	Contract string `json:"contract"` // original requested contract in the event
	Client   Client `json:"client"`   // contract client chain
	Address  string `json:"address"`  // contract address
	Function string `json:"function"` // contract function
}

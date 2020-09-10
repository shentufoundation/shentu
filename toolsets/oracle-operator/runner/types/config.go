package types

// Config defines the data type for configurations.
type Config struct {
	Combination Strategy `mapstructure:"strategy"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{}
}

// Strategy defines primitive combination.
type Strategy struct {
	Type       string      `mapstructure:"type"`
	Primitives []Primitive `mapstructure:"primitive"`
}

// Primitive defines primitive address and how to use that primitive.
type Primitive struct {
	PrimitiveContractAddr string  `mapstructure:"primitive_contract_address"`
	Weight                float32 `mapstructure:"weight"`
}

// AbiEntry defines the data type for abi entry.
type AbiEntry struct {
	Name string `json:"name"`
	Type string `json:"stateMutability"`
}

// MsgPrimitiveResponse defines returned score and related info from primitive.
type MsgPrimitiveResponse struct {
	Score uint8
	Primitive
}

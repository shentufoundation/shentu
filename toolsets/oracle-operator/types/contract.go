package types

// PrimitiveContractFnName is Security Primitive Contract Interface function name.
const PrimitiveContractFnName = "getInsight"

// ABIEntry defines the data type for abi entry.
type ABIEntry struct {
	Name string `json:"name"`
	Type string `json:"stateMutability"`
}

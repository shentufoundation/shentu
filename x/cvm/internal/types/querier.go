package types

import "github.com/hyperledger/burrow/acm"

// querier keys
const (
	QueryCode     = "code"
	QueryStorage  = "storage"
	QueryAbi      = "abi"
	QueryAddrMeta = "address-meta"
	QueryMeta     = "meta"
	QueryView     = "view"
	QueryAccount  = "account"
)

// QueryResView is the query result payload for a storage query.
type QueryResView struct {
	Ret []byte `json:"ret"`
}

// QueryResCode is the query result payload for a contract code query.
type QueryResCode struct {
	Code acm.Bytecode `json:"code"`
}

// String implements fmt.Stringer.
func (q QueryResCode) String() string {
	return q.Code.String()
}

// QueryResStorage is the query result payload for a storage query.
type QueryResStorage struct {
	Value []byte `json:"value"`
}

// String implements fmt.Stringer.
func (q QueryResStorage) String() string {
	return string(q.Value)
}

// QueryResAbi is the query result payload for a contract code ABI query.
type QueryResAbi struct {
	Abi []byte `json:"abi"`
}

// String implements fmt.Stringer.
func (q QueryResAbi) String() string {
	return string(q.Abi)
}

// QueryResAddrMeta is the query result payload for a contract code ABI query.
type QueryResAddrMeta struct {
	Metahash string `json:"metahash"`
}

// String implements fmt.Stringer.
func (q QueryResAddrMeta) String() string {
	return q.Metahash
}

// QueryResMeta is the query result payload for a contract code ABI query.
type QueryResMeta struct {
	Meta string `json:"meta"`
}

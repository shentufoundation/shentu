package types

import (
	"encoding/json"
	"strings"
)

type AddOrRemove bool

const (
	Add    AddOrRemove = false
	Remove AddOrRemove = true
)

func (aor AddOrRemove) String() string {
	switch aor {
	case Add:
		return "add"
	case Remove:
		return "remove"
	default:
		panic("invalid AddOrRemove value")
	}
}

func AddOrRemoveFromString(str string) (AddOrRemove, error) {
	switch strings.ToLower(str) {
	case "add":
		return Add, nil
	case "remove":
		return Remove, nil
	default:
		return Add, ErrAddOrRemove
	}
}

// MarshalJSON marshals to JSON using string.
func (aor AddOrRemove) MarshalJSON() ([]byte, error) {
	return json.Marshal(aor.String())
}

// UnmarshalJSON unmarshals from JSON using string values.
func (aor *AddOrRemove) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	value, err := AddOrRemoveFromString(s)
	if err != nil {
		return err
	}

	*aor = value
	return nil
}

func (aor AddOrRemove) ToProto() CertifierUpdateOperation {
	switch aor {
	case Add:
		return CertifierUpdateOperationAdd
	case Remove:
		return CertifierUpdateOperationRemove
	default:
		return CertifierUpdateOperationUnspecified
	}
}

func AddOrRemoveFromProto(op CertifierUpdateOperation) (AddOrRemove, error) {
	switch op {
	case CertifierUpdateOperationAdd:
		return Add, nil
	case CertifierUpdateOperationRemove:
		return Remove, nil
	default:
		return Add, ErrAddOrRemove
	}
}

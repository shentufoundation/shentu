package types

import (
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	proto "github.com/gogo/protobuf/proto"
)

// EncryptionKey is the interface for all kinds of Program EncryptionKey.
type EncryptionKey interface {
	proto.Message

	GetEncryptionKey() []byte
}

type EncryptedDesc interface {
	proto.Message

	GetEncryptedDesc() []byte
}

type EncryptedPoc interface {
	proto.Message

	GetEncryptedPoc() []byte
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces.
func (p Program) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	var pubKey EncryptionKey
	return unpacker.UnpackAny(p.EncryptionKey, &pubKey)
}

// GetEncryptionKey returns EncryptionKey of the Program.
func (p Program) GetEncryptionKey() EncryptionKey {
	pubKey, ok := p.EncryptionKey.GetCachedValue().(EncryptionKey)
	if !ok {
		return nil
	}
	return pubKey
}

func (f Finding) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	var desc EncryptedDesc
	var poc EncryptedPoc
	err := unpacker.UnpackAny(f.EncryptedDesc, &desc)
	if err != nil {
		return err
	}

	return unpacker.UnpackAny(f.EncryptedPoc, &poc)
}

func (f Finding) GetEncryptedDesc() EncryptedDesc {
	desc, ok := f.EncryptedDesc.GetCachedValue().(EncryptedDesc)
	if !ok {
		return nil
	}
	return desc
}

func (f Finding) GetEncryptedPoc() EncryptedPoc {
	poc, ok := f.EncryptedPoc.GetCachedValue().(EncryptedPoc)
	if !ok {
		return nil
	}
	return poc
}

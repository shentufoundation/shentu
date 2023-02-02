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

type FindingDesc interface {
	proto.Message

	GetFindingDesc() []byte
}

type FindingPoc interface {
	proto.Message

	GetFindingPoc() []byte
}

type FindingComment interface {
	proto.Message

	GetFindingComment() []byte
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
	var desc FindingDesc
	var poc FindingPoc
	var comment FindingComment
	err := unpacker.UnpackAny(f.FindingDesc, &desc)
	if err != nil {
		return err
	}
	if err = unpacker.UnpackAny(f.FindingComment, &comment); err != nil {
		return err
	}

	return unpacker.UnpackAny(f.FindingPoc, &poc)
}

func (f Finding) GetFindingDesc() FindingDesc {
	desc, ok := f.FindingDesc.GetCachedValue().(FindingDesc)
	if !ok {
		return nil
	}
	return desc
}

func (f Finding) GetFindingPoc() FindingPoc {
	poc, ok := f.FindingPoc.GetCachedValue().(FindingPoc)
	if !ok {
		return nil
	}
	return poc
}

func (f Finding) GetFindingComment() FindingComment {
	comment, ok := f.FindingComment.GetCachedValue().(FindingComment)
	if !ok {
		return nil
	}
	return comment
}

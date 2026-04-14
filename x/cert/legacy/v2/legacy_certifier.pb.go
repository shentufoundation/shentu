package v2

// LegacyCertifier mirrors the v1 on-wire shape of x/cert Certifier:
//   field 1: address (string)
//   field 3: proposer (string) -- removed from the current proto
//   field 4: description (string)
// It exists so the v1->v2 prefix migration can still translate the
// proposer bech32 prefix against historical data without reintroducing
// the field on the live Certifier type.

import (
	"fmt"
	"io"

	"github.com/cosmos/gogoproto/proto"
)

type LegacyCertifier struct {
	Address     string
	Proposer    string
	Description string
}

func (m *LegacyCertifier) Reset()         { *m = LegacyCertifier{} }
func (m *LegacyCertifier) String() string { return proto.CompactTextString(m) }
func (*LegacyCertifier) ProtoMessage()    {}

func (m *LegacyCertifier) Size() (n int) {
	if m == nil {
		return 0
	}
	if l := len(m.Address); l > 0 {
		n += 1 + l + sovLegacyCertifier(uint64(l))
	}
	if l := len(m.Proposer); l > 0 {
		n += 1 + l + sovLegacyCertifier(uint64(l))
	}
	if l := len(m.Description); l > 0 {
		n += 1 + l + sovLegacyCertifier(uint64(l))
	}
	return n
}

func (m *LegacyCertifier) Marshal() ([]byte, error) {
	size := m.Size()
	dAtA := make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyCertifier) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyCertifier) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintLegacyCertifier(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x22 // field 4, wire type 2
	}
	if len(m.Proposer) > 0 {
		i -= len(m.Proposer)
		copy(dAtA[i:], m.Proposer)
		i = encodeVarintLegacyCertifier(dAtA, i, uint64(len(m.Proposer)))
		i--
		dAtA[i] = 0x1a // field 3, wire type 2
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintLegacyCertifier(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa // field 1, wire type 2
	}
	return len(dAtA) - i, nil
}

func (m *LegacyCertifier) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return fmt.Errorf("LegacyCertifier: int overflow")
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType != 2 {
			return fmt.Errorf("LegacyCertifier: unexpected wireType %d for field %d", wireType, fieldNum)
		}
		var stringLen uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return fmt.Errorf("LegacyCertifier: int overflow")
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			stringLen |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		intStringLen := int(stringLen)
		if intStringLen < 0 {
			return fmt.Errorf("LegacyCertifier: invalid length")
		}
		postIndex := iNdEx + intStringLen
		if postIndex < 0 || postIndex > l {
			return io.ErrUnexpectedEOF
		}
		val := string(dAtA[iNdEx:postIndex])
		iNdEx = postIndex
		switch fieldNum {
		case 1:
			m.Address = val
		case 3:
			m.Proposer = val
		case 4:
			m.Description = val
		}
	}
	return nil
}

func sovLegacyCertifier(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}

func encodeVarintLegacyCertifier(dAtA []byte, offset int, v uint64) int {
	offset -= sovLegacyCertifier(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: shentu/bounty/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the bounty module's genesis state
type GenesisState struct {
	Programs []*Program `protobuf:"bytes,1,rep,name=programs,proto3" json:"programs,omitempty"`
	Findings []*Finding `protobuf:"bytes,2,rep,name=findings,proto3" json:"findings,omitempty"`
	// starting_theorem_id is the ID of the starting theorem
	StartingTheoremId uint64 `protobuf:"varint,3,opt,name=starting_theorem_id,json=startingTheoremId,proto3" json:"starting_theorem_id,omitempty"`
	// theorems defines all the theorems present at genesis.
	Theorems []*Theorem `protobuf:"bytes,4,rep,name=theorems,proto3" json:"theorems,omitempty"`
	// proofs defines all the proofs present at genesis.
	Proofs []*Proof `protobuf:"bytes,5,rep,name=proofs,proto3" json:"proofs,omitempty"`
	// grants defines all the grants present at genesis.
	Grants   []*Grant   `protobuf:"bytes,6,rep,name=grants,proto3" json:"grants,omitempty"`
	Deposits []*Deposit `protobuf:"bytes,7,rep,name=deposits,proto3" json:"deposits,omitempty"`
	Rewards  []*Reward  `protobuf:"bytes,8,rep,name=rewards,proto3" json:"rewards,omitempty"`
	Params   *Params    `protobuf:"bytes,9,opt,name=params,proto3" json:"params,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_186d656250aa7272, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetPrograms() []*Program {
	if m != nil {
		return m.Programs
	}
	return nil
}

func (m *GenesisState) GetFindings() []*Finding {
	if m != nil {
		return m.Findings
	}
	return nil
}

func (m *GenesisState) GetStartingTheoremId() uint64 {
	if m != nil {
		return m.StartingTheoremId
	}
	return 0
}

func (m *GenesisState) GetTheorems() []*Theorem {
	if m != nil {
		return m.Theorems
	}
	return nil
}

func (m *GenesisState) GetProofs() []*Proof {
	if m != nil {
		return m.Proofs
	}
	return nil
}

func (m *GenesisState) GetGrants() []*Grant {
	if m != nil {
		return m.Grants
	}
	return nil
}

func (m *GenesisState) GetDeposits() []*Deposit {
	if m != nil {
		return m.Deposits
	}
	return nil
}

func (m *GenesisState) GetRewards() []*Reward {
	if m != nil {
		return m.Rewards
	}
	return nil
}

func (m *GenesisState) GetParams() *Params {
	if m != nil {
		return m.Params
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "shentu.bounty.v1.GenesisState")
}

func init() { proto.RegisterFile("shentu/bounty/v1/genesis.proto", fileDescriptor_186d656250aa7272) }

var fileDescriptor_186d656250aa7272 = []byte{
	// 356 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xb1, 0x4e, 0xe3, 0x30,
	0x18, 0x80, 0x9b, 0x6b, 0x2f, 0xed, 0xe5, 0x6e, 0x38, 0x02, 0x12, 0xa6, 0x12, 0x51, 0xc5, 0xd4,
	0x29, 0xa1, 0x45, 0xbc, 0x00, 0x42, 0x54, 0x88, 0xa5, 0x32, 0x4c, 0x2c, 0x95, 0x4b, 0x1c, 0xd7,
	0x43, 0xfd, 0x47, 0xb6, 0x53, 0xe8, 0x5b, 0xf0, 0x58, 0x8c, 0x1d, 0x19, 0x51, 0xfa, 0x22, 0x28,
	0xb6, 0x9b, 0x81, 0x36, 0x5b, 0x94, 0xef, 0xfb, 0xfe, 0xdf, 0x51, 0x1c, 0x44, 0x6a, 0x41, 0x85,
	0x2e, 0x92, 0x39, 0x14, 0x42, 0xaf, 0x93, 0xd5, 0x28, 0x61, 0x54, 0x50, 0xc5, 0x55, 0x9c, 0x4b,
	0xd0, 0x10, 0xfe, 0xb7, 0x3c, 0xb6, 0x3c, 0x5e, 0x8d, 0xfa, 0x27, 0x0c, 0x18, 0x18, 0x98, 0x54,
	0x4f, 0xd6, 0xeb, 0x9f, 0xef, 0xcd, 0x71, 0x85, 0xc1, 0x17, 0x65, 0x3b, 0xf8, 0x37, 0xb1, 0x83,
	0x1f, 0x35, 0xd1, 0x34, 0xbc, 0x0e, 0x7a, 0xb9, 0x04, 0x26, 0xc9, 0x52, 0x21, 0x6f, 0xd0, 0x1e,
	0xfe, 0x1d, 0x9f, 0xc5, 0x3f, 0x57, 0xc5, 0x53, 0x6b, 0xe0, 0x5a, 0xad, 0xb2, 0x8c, 0x8b, 0x94,
	0x0b, 0xa6, 0xd0, 0xaf, 0xa6, 0xec, 0xce, 0x1a, 0xb8, 0x56, 0xc3, 0x38, 0x38, 0x56, 0x9a, 0x48,
	0xcd, 0x05, 0x9b, 0xe9, 0x05, 0x05, 0x49, 0x97, 0x33, 0x9e, 0xa2, 0xf6, 0xc0, 0x1b, 0x76, 0xf0,
	0xd1, 0x0e, 0x3d, 0x59, 0x72, 0x9f, 0x56, 0x6b, 0x9c, 0xa6, 0x50, 0xa7, 0x69, 0x8d, 0xd3, 0x71,
	0xad, 0x86, 0x49, 0xe0, 0xe7, 0x12, 0x20, 0x53, 0xe8, 0xb7, 0x89, 0x4e, 0x0f, 0x7e, 0x12, 0x64,
	0xd8, 0x69, 0x55, 0xc0, 0x24, 0x11, 0x5a, 0x21, 0xbf, 0x29, 0x98, 0x54, 0x1c, 0x3b, 0xad, 0x3a,
	0x58, 0x4a, 0x73, 0x50, 0x5c, 0x2b, 0xd4, 0x6d, 0x3a, 0xd8, 0xad, 0x35, 0x70, 0xad, 0x86, 0xe3,
	0xa0, 0x2b, 0xe9, 0x2b, 0x91, 0xa9, 0x42, 0x3d, 0x53, 0xa1, 0xfd, 0x0a, 0x1b, 0x01, 0xef, 0xc4,
	0xf0, 0x32, 0xf0, 0x73, 0x62, 0xfe, 0xcf, 0x9f, 0x81, 0x77, 0x38, 0x99, 0x1a, 0x8e, 0x9d, 0x77,
	0xf3, 0xf0, 0x51, 0x46, 0xde, 0xa6, 0x8c, 0xbc, 0xaf, 0x32, 0xf2, 0xde, 0xb7, 0x51, 0x6b, 0xb3,
	0x8d, 0x5a, 0x9f, 0xdb, 0xa8, 0xf5, 0x3c, 0x62, 0x5c, 0x2f, 0x8a, 0x79, 0xfc, 0x02, 0xcb, 0xc4,
	0x4e, 0xc9, 0xa0, 0x10, 0x29, 0xd1, 0x1c, 0x84, 0x7b, 0x91, 0xbc, 0xed, 0xee, 0x8e, 0x5e, 0xe7,
	0x54, 0xcd, 0x7d, 0x73, 0x71, 0xae, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0xee, 0x81, 0x9c, 0xec,
	0xa1, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Params != nil {
		{
			size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
	}
	if len(m.Rewards) > 0 {
		for iNdEx := len(m.Rewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Rewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.Deposits) > 0 {
		for iNdEx := len(m.Deposits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Deposits[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.Grants) > 0 {
		for iNdEx := len(m.Grants) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Grants[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Proofs) > 0 {
		for iNdEx := len(m.Proofs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Proofs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Theorems) > 0 {
		for iNdEx := len(m.Theorems) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Theorems[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.StartingTheoremId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.StartingTheoremId))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Findings) > 0 {
		for iNdEx := len(m.Findings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Findings[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Programs) > 0 {
		for iNdEx := len(m.Programs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Programs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Programs) > 0 {
		for _, e := range m.Programs {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Findings) > 0 {
		for _, e := range m.Findings {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.StartingTheoremId != 0 {
		n += 1 + sovGenesis(uint64(m.StartingTheoremId))
	}
	if len(m.Theorems) > 0 {
		for _, e := range m.Theorems {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Proofs) > 0 {
		for _, e := range m.Proofs {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Grants) > 0 {
		for _, e := range m.Grants {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Deposits) > 0 {
		for _, e := range m.Deposits {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Rewards) > 0 {
		for _, e := range m.Rewards {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.Params != nil {
		l = m.Params.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Programs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Programs = append(m.Programs, &Program{})
			if err := m.Programs[len(m.Programs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Findings", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Findings = append(m.Findings, &Finding{})
			if err := m.Findings[len(m.Findings)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartingTheoremId", wireType)
			}
			m.StartingTheoremId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartingTheoremId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Theorems", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Theorems = append(m.Theorems, &Theorem{})
			if err := m.Theorems[len(m.Theorems)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proofs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Proofs = append(m.Proofs, &Proof{})
			if err := m.Proofs[len(m.Proofs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Grants", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Grants = append(m.Grants, &Grant{})
			if err := m.Grants[len(m.Grants)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposits", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Deposits = append(m.Deposits, &Deposit{})
			if err := m.Deposits[len(m.Deposits)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rewards = append(m.Rewards, &Reward{})
			if err := m.Rewards[len(m.Rewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Params == nil {
				m.Params = &Params{}
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)

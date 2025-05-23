// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: shentu/cert/v1alpha1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types"
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

type GenesisState struct {
	Certifiers        []Certifier   `protobuf:"bytes,1,rep,name=certifiers,proto3" json:"certifiers" yaml:"certifiers"`
	Platforms         []Platform    `protobuf:"bytes,2,rep,name=platforms,proto3" json:"platforms" yaml:"platforms"`
	Certificates      []Certificate `protobuf:"bytes,3,rep,name=certificates,proto3" json:"certificates" yaml:"certificates"`
	Libraries         []Library     `protobuf:"bytes,4,rep,name=libraries,proto3" json:"libraries" yaml:"libraries"`
	NextCertificateId uint64        `protobuf:"varint,5,opt,name=next_certificate_id,json=nextCertificateId,proto3" json:"next_certificate_id,omitempty" yaml:"next_certificate_id"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_860284e2a718f650, []int{0}
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

func init() {
	proto.RegisterType((*GenesisState)(nil), "shentu.cert.v1alpha1.GenesisState")
}

func init() {
	proto.RegisterFile("shentu/cert/v1alpha1/genesis.proto", fileDescriptor_860284e2a718f650)
}

var fileDescriptor_860284e2a718f650 = []byte{
	// 418 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x3f, 0x8f, 0xd3, 0x30,
	0x18, 0x87, 0x13, 0x7a, 0x20, 0xce, 0xdc, 0xc0, 0xe5, 0x6e, 0xc8, 0x15, 0xe1, 0x1c, 0x9e, 0x3a,
	0xc5, 0x0a, 0x6c, 0x1d, 0xc3, 0x80, 0x2a, 0x21, 0x84, 0x82, 0x60, 0xe8, 0x52, 0x39, 0xa9, 0x93,
	0x5a, 0x4a, 0xe2, 0x28, 0x76, 0xaa, 0xe6, 0x1b, 0x30, 0xf2, 0x11, 0xfa, 0x71, 0x3a, 0x76, 0x83,
	0xa9, 0x42, 0xed, 0xc2, 0xdc, 0x4f, 0x80, 0x62, 0xa7, 0xff, 0x50, 0xc4, 0x66, 0xfb, 0x7d, 0xde,
	0xe7, 0xfd, 0xbd, 0x92, 0x01, 0x12, 0x33, 0x9a, 0xcb, 0x0a, 0x47, 0xb4, 0x94, 0x78, 0xee, 0x91,
	0xb4, 0x98, 0x11, 0x0f, 0x27, 0x34, 0xa7, 0x82, 0x09, 0xb7, 0x28, 0xb9, 0xe4, 0xd6, 0xbd, 0x66,
	0xdc, 0x86, 0x71, 0x0f, 0x4c, 0xff, 0x3e, 0xe1, 0x09, 0x57, 0x00, 0x6e, 0x4e, 0x9a, 0xed, 0x3f,
	0x24, 0x9c, 0x27, 0x29, 0xc5, 0xea, 0x16, 0x56, 0x31, 0x26, 0x79, 0xdd, 0x96, 0x60, 0xc4, 0x45,
	0xc6, 0x05, 0x0e, 0x89, 0xa0, 0x78, 0xee, 0x85, 0x54, 0x12, 0x0f, 0x47, 0x9c, 0xe5, 0x87, 0x56,
	0x5d, 0x9f, 0x68, 0xa7, 0xbe, 0xb4, 0x25, 0xa7, 0x33, 0xa5, 0xca, 0xa3, 0x00, 0xf4, 0xb3, 0x07,
	0x6e, 0x3e, 0xe8, 0xd0, 0x5f, 0x24, 0x91, 0xd4, 0x1a, 0x03, 0xd0, 0x94, 0x59, 0xcc, 0x68, 0x29,
	0x6c, 0xf3, 0xb1, 0x37, 0x78, 0xf1, 0xd6, 0x71, 0xbb, 0x16, 0x71, 0xdf, 0x1f, 0x38, 0xff, 0x61,
	0xb5, 0x71, 0x8c, 0xfd, 0xc6, 0xb9, 0xad, 0x49, 0x96, 0x0e, 0xd1, 0x49, 0x80, 0x82, 0x33, 0x9b,
	0xf5, 0x0d, 0x5c, 0x17, 0x29, 0x91, 0x31, 0x2f, 0x33, 0x61, 0x3f, 0x51, 0x6a, 0xd8, 0xad, 0xfe,
	0xdc, 0x62, 0xbe, 0xdd, 0x9a, 0x5f, 0x6a, 0xf3, 0xb1, 0x1d, 0x05, 0x27, 0x95, 0x15, 0x82, 0x9b,
	0x76, 0x4a, 0x44, 0x24, 0x15, 0x76, 0x4f, 0xa9, 0xdf, 0xfc, 0x37, 0x75, 0x43, 0xfa, 0xaf, 0x5a,
	0xfb, 0xdd, 0x45, 0x6e, 0x25, 0x41, 0xc1, 0x85, 0xd3, 0xfa, 0x0a, 0xae, 0x53, 0x16, 0x96, 0xa4,
	0x64, 0x54, 0xd8, 0x57, 0x6a, 0xc0, 0xeb, 0xee, 0x01, 0x1f, 0x15, 0x56, 0xff, 0x1b, 0xfd, 0xd8,
	0x8d, 0x82, 0x93, 0xc9, 0xfa, 0x04, 0xee, 0x72, 0xba, 0x90, 0x93, 0xb3, 0x59, 0x13, 0x36, 0xb5,
	0x9f, 0x3e, 0x9a, 0x83, 0x2b, 0x1f, 0xee, 0x37, 0x4e, 0x5f, 0x77, 0x77, 0x40, 0x28, 0xb8, 0x6d,
	0x5e, 0xcf, 0xf6, 0x19, 0x4d, 0x87, 0xcf, 0xbf, 0x2f, 0x1d, 0xe3, 0xcf, 0xd2, 0x31, 0xfc, 0xd1,
	0x6a, 0x0b, 0xcd, 0xf5, 0x16, 0x9a, 0xbf, 0xb7, 0xd0, 0xfc, 0xb1, 0x83, 0xc6, 0x7a, 0x07, 0x8d,
	0x5f, 0x3b, 0x68, 0x8c, 0x71, 0xc2, 0xe4, 0xac, 0x0a, 0xdd, 0x88, 0x67, 0x58, 0x6f, 0x10, 0xf3,
	0x2a, 0x9f, 0x12, 0xc9, 0x78, 0xde, 0x3e, 0xe0, 0x85, 0xfe, 0x32, 0xb2, 0x2e, 0xa8, 0x08, 0x9f,
	0xa9, 0xbf, 0xf2, 0xee, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x08, 0x42, 0xa3, 0xc7, 0xf4, 0x02,
	0x00, 0x00,
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
	if m.NextCertificateId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.NextCertificateId))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Libraries) > 0 {
		for iNdEx := len(m.Libraries) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Libraries[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Certificates) > 0 {
		for iNdEx := len(m.Certificates) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Certificates[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Platforms) > 0 {
		for iNdEx := len(m.Platforms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Platforms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Certifiers) > 0 {
		for iNdEx := len(m.Certifiers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Certifiers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Certifiers) > 0 {
		for _, e := range m.Certifiers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Platforms) > 0 {
		for _, e := range m.Platforms {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Certificates) > 0 {
		for _, e := range m.Certificates {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Libraries) > 0 {
		for _, e := range m.Libraries {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.NextCertificateId != 0 {
		n += 1 + sovGenesis(uint64(m.NextCertificateId))
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
				return fmt.Errorf("proto: wrong wireType = %d for field Certifiers", wireType)
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
			m.Certifiers = append(m.Certifiers, Certifier{})
			if err := m.Certifiers[len(m.Certifiers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Platforms", wireType)
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
			m.Platforms = append(m.Platforms, Platform{})
			if err := m.Platforms[len(m.Platforms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Certificates", wireType)
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
			m.Certificates = append(m.Certificates, Certificate{})
			if err := m.Certificates[len(m.Certificates)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Libraries", wireType)
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
			m.Libraries = append(m.Libraries, Library{})
			if err := m.Libraries[len(m.Libraries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextCertificateId", wireType)
			}
			m.NextCertificateId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextCertificateId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: topchain/pin/pin_request.proto

package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	proto "github.com/cosmos/gogoproto/proto"
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

type PinRequest struct {
	Index     string `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	CroId     string `protobuf:"bytes,2,opt,name=croId,proto3" json:"croId,omitempty"`
	Amount    string `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Duration  string `protobuf:"bytes,4,opt,name=duration,proto3" json:"duration,omitempty"`
	Submitter string `protobuf:"bytes,5,opt,name=submitter,proto3" json:"submitter,omitempty"`
}

func (m *PinRequest) Reset()         { *m = PinRequest{} }
func (m *PinRequest) String() string { return proto.CompactTextString(m) }
func (*PinRequest) ProtoMessage()    {}
func (*PinRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_bdf8b5b4dd4bd537, []int{0}
}
func (m *PinRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PinRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PinRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PinRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PinRequest.Merge(m, src)
}
func (m *PinRequest) XXX_Size() int {
	return m.Size()
}
func (m *PinRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PinRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PinRequest proto.InternalMessageInfo

func (m *PinRequest) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *PinRequest) GetCroId() string {
	if m != nil {
		return m.CroId
	}
	return ""
}

func (m *PinRequest) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func (m *PinRequest) GetDuration() string {
	if m != nil {
		return m.Duration
	}
	return ""
}

func (m *PinRequest) GetSubmitter() string {
	if m != nil {
		return m.Submitter
	}
	return ""
}

func init() {
	proto.RegisterType((*PinRequest)(nil), "topchain.pin.PinRequest")
}

func init() { proto.RegisterFile("topchain/pin/pin_request.proto", fileDescriptor_bdf8b5b4dd4bd537) }

var fileDescriptor_bdf8b5b4dd4bd537 = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2b, 0xc9, 0x2f, 0x48,
	0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x2f, 0x80, 0xe0, 0xf8, 0xa2, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12,
	0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x1e, 0x98, 0xbc, 0x5e, 0x41, 0x66, 0x9e, 0x52, 0x17,
	0x23, 0x17, 0x57, 0x40, 0x66, 0x5e, 0x10, 0x44, 0x89, 0x90, 0x08, 0x17, 0x6b, 0x66, 0x5e, 0x4a,
	0x6a, 0x85, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x84, 0x03, 0x12, 0x4d, 0x2e, 0xca, 0xf7,
	0x4c, 0x91, 0x60, 0x82, 0x88, 0x82, 0x39, 0x42, 0x62, 0x5c, 0x6c, 0x89, 0xb9, 0xf9, 0xa5, 0x79,
	0x25, 0x12, 0xcc, 0x60, 0x61, 0x28, 0x4f, 0x48, 0x8a, 0x8b, 0x23, 0xa5, 0xb4, 0x28, 0xb1, 0x24,
	0x33, 0x3f, 0x4f, 0x82, 0x05, 0x2c, 0x03, 0xe7, 0x0b, 0xc9, 0x70, 0x71, 0x16, 0x97, 0x26, 0xe5,
	0x66, 0x96, 0x94, 0xa4, 0x16, 0x49, 0xb0, 0x82, 0x25, 0x11, 0x02, 0x4e, 0x7a, 0x27, 0x1e, 0xc9,
	0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e,
	0xcb, 0x31, 0xdc, 0x78, 0x2c, 0xc7, 0x10, 0x25, 0x02, 0xf7, 0x54, 0x05, 0xd8, 0x5b, 0x25, 0x95,
	0x05, 0xa9, 0xc5, 0x49, 0x6c, 0x60, 0x1f, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x6e, 0x0e,
	0x0c, 0xb0, 0xf3, 0x00, 0x00, 0x00,
}

func (m *PinRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PinRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PinRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Submitter) > 0 {
		i -= len(m.Submitter)
		copy(dAtA[i:], m.Submitter)
		i = encodeVarintPinRequest(dAtA, i, uint64(len(m.Submitter)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Duration) > 0 {
		i -= len(m.Duration)
		copy(dAtA[i:], m.Duration)
		i = encodeVarintPinRequest(dAtA, i, uint64(len(m.Duration)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Amount) > 0 {
		i -= len(m.Amount)
		copy(dAtA[i:], m.Amount)
		i = encodeVarintPinRequest(dAtA, i, uint64(len(m.Amount)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.CroId) > 0 {
		i -= len(m.CroId)
		copy(dAtA[i:], m.CroId)
		i = encodeVarintPinRequest(dAtA, i, uint64(len(m.CroId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintPinRequest(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPinRequest(dAtA []byte, offset int, v uint64) int {
	offset -= sovPinRequest(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PinRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovPinRequest(uint64(l))
	}
	l = len(m.CroId)
	if l > 0 {
		n += 1 + l + sovPinRequest(uint64(l))
	}
	l = len(m.Amount)
	if l > 0 {
		n += 1 + l + sovPinRequest(uint64(l))
	}
	l = len(m.Duration)
	if l > 0 {
		n += 1 + l + sovPinRequest(uint64(l))
	}
	l = len(m.Submitter)
	if l > 0 {
		n += 1 + l + sovPinRequest(uint64(l))
	}
	return n
}

func sovPinRequest(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPinRequest(x uint64) (n int) {
	return sovPinRequest(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PinRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPinRequest
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
			return fmt.Errorf("proto: PinRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PinRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPinRequest
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
				return ErrInvalidLengthPinRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPinRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CroId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPinRequest
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
				return ErrInvalidLengthPinRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPinRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CroId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPinRequest
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
				return ErrInvalidLengthPinRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPinRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPinRequest
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
				return ErrInvalidLengthPinRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPinRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Duration = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Submitter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPinRequest
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
				return ErrInvalidLengthPinRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPinRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Submitter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPinRequest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPinRequest
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
func skipPinRequest(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPinRequest
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
					return 0, ErrIntOverflowPinRequest
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
					return 0, ErrIntOverflowPinRequest
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
				return 0, ErrInvalidLengthPinRequest
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPinRequest
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPinRequest
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPinRequest        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPinRequest          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPinRequest = fmt.Errorf("proto: unexpected end of group")
)
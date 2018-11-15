// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: endpoint.proto

package v1 // import "github.com/solo-io/supergloo/pkg/api/external/gloo/v1"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

import bytes "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

//
// @solo-kit:resource.short_name=ep
// @solo-kit:resource.plural_name=endpoints
// @solo-kit:resource.resource_groups=api.gloo.solo.io
//
// Endpoints represent dynamically discovered address/ports where an upstream service is listening
type Endpoint struct {
	// List of the upstreams the endpoint belongs to
	Upstreams []*core.ResourceRef `protobuf:"bytes,1,rep,name=upstreams" json:"upstreams,omitempty"`
	// Address of the endpoint (ip or hostname)
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// listening port for the endpoint
	Port uint32 `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	// Metadata contains the object metadata for this resource
	Metadata             core.Metadata `protobuf:"bytes,7,opt,name=metadata" json:"metadata"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}
func (*Endpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_94e1c2da017cfd23, []int{0}
}
func (m *Endpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoint.Unmarshal(m, b)
}
func (m *Endpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoint.Marshal(b, m, deterministic)
}
func (dst *Endpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoint.Merge(dst, src)
}
func (m *Endpoint) XXX_Size() int {
	return xxx_messageInfo_Endpoint.Size(m)
}
func (m *Endpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoint.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoint proto.InternalMessageInfo

func (m *Endpoint) GetUpstreams() []*core.ResourceRef {
	if m != nil {
		return m.Upstreams
	}
	return nil
}

func (m *Endpoint) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Endpoint) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Endpoint) GetMetadata() core.Metadata {
	if m != nil {
		return m.Metadata
	}
	return core.Metadata{}
}

func init() {
	proto.RegisterType((*Endpoint)(nil), "gloo.solo.io.Endpoint")
}
func (this *Endpoint) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Endpoint)
	if !ok {
		that2, ok := that.(Endpoint)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Upstreams) != len(that1.Upstreams) {
		return false
	}
	for i := range this.Upstreams {
		if !this.Upstreams[i].Equal(that1.Upstreams[i]) {
			return false
		}
	}
	if this.Address != that1.Address {
		return false
	}
	if this.Port != that1.Port {
		return false
	}
	if !this.Metadata.Equal(&that1.Metadata) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}

func init() { proto.RegisterFile("endpoint.proto", fileDescriptor_endpoint_94e1c2da017cfd23) }

var fileDescriptor_endpoint_94e1c2da017cfd23 = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xb1, 0x4e, 0xf3, 0x30,
	0x14, 0x85, 0x7f, 0xff, 0xad, 0x68, 0xeb, 0x02, 0x83, 0x85, 0x50, 0xe8, 0x00, 0x11, 0x53, 0x06,
	0xb0, 0x95, 0x22, 0x04, 0x12, 0x5b, 0x25, 0x46, 0x16, 0x8f, 0x6c, 0x6e, 0x72, 0x1b, 0xac, 0x26,
	0xb9, 0x96, 0xed, 0x54, 0x3c, 0x12, 0x12, 0x2f, 0xc2, 0x53, 0x30, 0xf0, 0x24, 0x28, 0x4e, 0x02,
	0x42, 0x62, 0x60, 0xf2, 0xb5, 0xcf, 0x77, 0xe4, 0x73, 0x2e, 0x3d, 0x84, 0x3a, 0x37, 0xa8, 0x6b,
	0xcf, 0x8d, 0x45, 0x8f, 0x6c, 0xbf, 0x28, 0x11, 0xb9, 0xc3, 0x12, 0xb9, 0xc6, 0xc5, 0x51, 0x81,
	0x05, 0x06, 0x41, 0xb4, 0x53, 0xc7, 0x2c, 0xd2, 0x42, 0xfb, 0xa7, 0x66, 0xcd, 0x33, 0xac, 0x44,
	0x4b, 0x5e, 0x6a, 0xec, 0xce, 0xad, 0xf6, 0x42, 0x19, 0x2d, 0x76, 0xa9, 0xa8, 0xc0, 0xab, 0x5c,
	0x79, 0xd5, 0x5b, 0x2e, 0xfe, 0x60, 0xb1, 0xb0, 0xe9, 0xe8, 0xf3, 0x57, 0x42, 0xa7, 0xf7, 0x7d,
	0x2e, 0x76, 0x43, 0x67, 0x8d, 0x71, 0xde, 0x82, 0xaa, 0x5c, 0x44, 0xe2, 0x51, 0x32, 0x5f, 0x9e,
	0xf0, 0x0c, 0x2d, 0x0c, 0x29, 0xb9, 0x04, 0x87, 0x8d, 0xcd, 0x40, 0xc2, 0x46, 0x7e, 0xb3, 0x2c,
	0xa2, 0x13, 0x95, 0xe7, 0x16, 0x9c, 0x8b, 0xfe, 0xc7, 0x24, 0x99, 0xc9, 0xe1, 0xca, 0x18, 0x1d,
	0x1b, 0xb4, 0x3e, 0x1a, 0xc5, 0x24, 0x39, 0x90, 0x61, 0x66, 0xb7, 0x74, 0x3a, 0x64, 0x8e, 0x26,
	0x31, 0x49, 0xe6, 0xcb, 0xe3, 0x9f, 0xbf, 0x3c, 0xf4, 0xea, 0x6a, 0xfc, 0xf6, 0x7e, 0xf6, 0x4f,
	0x7e, 0xd1, 0xab, 0xbb, 0x97, 0x8f, 0x53, 0xf2, 0x78, 0xfd, 0x5b, 0xc3, 0xc6, 0x80, 0x6d, 0x17,
	0x2a, 0xcc, 0xb6, 0x08, 0x35, 0xe1, 0xd9, 0x83, 0xad, 0x55, 0x29, 0xc2, 0xeb, 0x2e, 0x5d, 0xef,
	0x85, 0xc6, 0x57, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x71, 0x9c, 0xfc, 0x3a, 0x88, 0x01, 0x00,
	0x00,
}
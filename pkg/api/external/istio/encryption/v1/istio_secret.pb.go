// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: istio_secret.proto

package v1 // import "github.com/solo-io/supergloo/pkg/api/external/istio/encryption/v1"

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
// @solo-kit:resource.short_name=ics
// @solo-kit:resource.plural_name=istiocerts
// @solo-kit:resource.resource_groups=translator.supergloo.solo.io,install.supergloo.solo.io
// Secret containing CA Certs for Istio
// Structured TLS Secret that istio uses for non-default root certificates
type IstioCacertsSecret struct {
	RootCert  string `protobuf:"bytes,1,opt,name=root_cert,json=root-cert.pem,proto3" json:"root_cert,omitempty"`
	CertChain string `protobuf:"bytes,2,opt,name=cert_chain,json=cert-chain.pem,proto3" json:"cert_chain,omitempty"`
	CaCert    string `protobuf:"bytes,3,opt,name=ca_cert,json=ca-cert.pem,proto3" json:"ca_cert,omitempty"`
	CaKey     string `protobuf:"bytes,4,opt,name=ca_key,json=ca-key.pem,proto3" json:"ca_key,omitempty"`
	// Metadata contains the object metadata for this resource
	Metadata             core.Metadata `protobuf:"bytes,7,opt,name=metadata" json:"metadata"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *IstioCacertsSecret) Reset()         { *m = IstioCacertsSecret{} }
func (m *IstioCacertsSecret) String() string { return proto.CompactTextString(m) }
func (*IstioCacertsSecret) ProtoMessage()    {}
func (*IstioCacertsSecret) Descriptor() ([]byte, []int) {
	return fileDescriptor_istio_secret_51f67eb0eba329ab, []int{0}
}
func (m *IstioCacertsSecret) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IstioCacertsSecret.Unmarshal(m, b)
}
func (m *IstioCacertsSecret) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IstioCacertsSecret.Marshal(b, m, deterministic)
}
func (dst *IstioCacertsSecret) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IstioCacertsSecret.Merge(dst, src)
}
func (m *IstioCacertsSecret) XXX_Size() int {
	return xxx_messageInfo_IstioCacertsSecret.Size(m)
}
func (m *IstioCacertsSecret) XXX_DiscardUnknown() {
	xxx_messageInfo_IstioCacertsSecret.DiscardUnknown(m)
}

var xxx_messageInfo_IstioCacertsSecret proto.InternalMessageInfo

func (m *IstioCacertsSecret) GetRootCert() string {
	if m != nil {
		return m.RootCert
	}
	return ""
}

func (m *IstioCacertsSecret) GetCertChain() string {
	if m != nil {
		return m.CertChain
	}
	return ""
}

func (m *IstioCacertsSecret) GetCaCert() string {
	if m != nil {
		return m.CaCert
	}
	return ""
}

func (m *IstioCacertsSecret) GetCaKey() string {
	if m != nil {
		return m.CaKey
	}
	return ""
}

func (m *IstioCacertsSecret) GetMetadata() core.Metadata {
	if m != nil {
		return m.Metadata
	}
	return core.Metadata{}
}

func init() {
	proto.RegisterType((*IstioCacertsSecret)(nil), "encryption.istio.io.IstioCacertsSecret")
}
func (this *IstioCacertsSecret) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*IstioCacertsSecret)
	if !ok {
		that2, ok := that.(IstioCacertsSecret)
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
	if this.RootCert != that1.RootCert {
		return false
	}
	if this.CertChain != that1.CertChain {
		return false
	}
	if this.CaCert != that1.CaCert {
		return false
	}
	if this.CaKey != that1.CaKey {
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

func init() { proto.RegisterFile("istio_secret.proto", fileDescriptor_istio_secret_51f67eb0eba329ab) }

var fileDescriptor_istio_secret_51f67eb0eba329ab = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0xc6, 0x09, 0x54, 0x2d, 0x75, 0x05, 0x43, 0x40, 0x28, 0xaa, 0x10, 0x54, 0x9d, 0xba, 0xc4,
	0x56, 0x61, 0x61, 0xa5, 0x0c, 0x88, 0x81, 0xa5, 0x6c, 0x2c, 0x91, 0x6b, 0x4e, 0xa9, 0x95, 0x3f,
	0x67, 0x39, 0xd7, 0x8a, 0xbc, 0x11, 0x8f, 0xc2, 0xc4, 0x23, 0x30, 0xf0, 0x24, 0xc8, 0x17, 0xa0,
	0x0b, 0x93, 0x4f, 0xdf, 0xf7, 0xf9, 0xb3, 0xef, 0x27, 0x62, 0xdb, 0x90, 0xc5, 0xac, 0x01, 0xe3,
	0x81, 0xa4, 0xf3, 0x48, 0x18, 0x9f, 0x40, 0x6d, 0x7c, 0xeb, 0xc8, 0x62, 0x2d, 0xd9, 0x96, 0x16,
	0xc7, 0xa7, 0x39, 0xe6, 0xc8, 0xbe, 0x0a, 0x53, 0x17, 0x1d, 0xcf, 0x73, 0x4b, 0xeb, 0xcd, 0x4a,
	0x1a, 0xac, 0x54, 0x83, 0x25, 0xa6, 0x16, 0xbb, 0xb3, 0xb0, 0xa4, 0xb4, 0xb3, 0x6a, 0x3b, 0x57,
	0x15, 0x90, 0x7e, 0xd1, 0xa4, 0xbb, 0x2b, 0xd3, 0x8f, 0x48, 0xc4, 0x0f, 0xa1, 0xf5, 0x4e, 0x1b,
	0xf0, 0xd4, 0x3c, 0xf1, 0xd3, 0xf1, 0x44, 0x0c, 0x3d, 0x22, 0x65, 0x41, 0x4b, 0xa2, 0x49, 0x34,
	0x1b, 0x2e, 0x8f, 0x82, 0x90, 0x06, 0x41, 0x3a, 0xa8, 0xe2, 0xa9, 0x10, 0x61, 0xce, 0xcc, 0x5a,
	0xdb, 0x3a, 0xd9, 0xe7, 0xc8, 0x71, 0x50, 0x52, 0x56, 0x38, 0x73, 0x2e, 0x06, 0x46, 0x77, 0x1d,
	0x07, 0x1c, 0x18, 0x19, 0xbd, 0x6b, 0x18, 0x8b, 0xbe, 0xd1, 0x59, 0x01, 0x6d, 0xd2, 0x63, 0x53,
	0x18, 0x9d, 0x16, 0xd0, 0xb2, 0x77, 0x23, 0x0e, 0x7f, 0x3f, 0x9a, 0x0c, 0x26, 0xd1, 0x6c, 0x74,
	0x75, 0x26, 0x0d, 0x7a, 0x90, 0x61, 0x1d, 0x69, 0x51, 0x3e, 0xfe, 0xb8, 0x8b, 0xde, 0xfb, 0xe7,
	0xe5, 0xde, 0xf2, 0x2f, 0xbd, 0xb8, 0x7f, 0xfb, 0xba, 0x88, 0x9e, 0x6f, 0xff, 0x23, 0xb1, 0x71,
	0xe0, 0xf3, 0x12, 0x51, 0xb9, 0x22, 0x67, 0x1c, 0xf0, 0x4a, 0xe0, 0x6b, 0x5d, 0x2a, 0xe6, 0xaa,
	0x76, 0xa0, 0xd5, 0x76, 0xbe, 0xea, 0x33, 0xa0, 0xeb, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x91,
	0xae, 0x5a, 0xc0, 0x94, 0x01, 0x00, 0x00,
}

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/solo-io/supergloo/api/external/prometheus/v1/config.proto

package v1 // import "github.com/solo-io/supergloo/pkg/api/external/prometheus/v1"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import types "github.com/gogo/protobuf/types"
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
// @solo-kit:resource.short_name=pcf
// @solo-kit:resource.plural_name=prometheusconfigs
//
// Prometheus Config
type PrometheusConfig struct {
	// json_name must refer to the data key in the configmap we expect
	Prometheus *types.Struct `protobuf:"bytes,1,opt,name=prometheus,json=prometheus.yml" json:"prometheus,omitempty"`
	// Metadata contains the object metadata for this resource
	Metadata             core.Metadata `protobuf:"bytes,7,opt,name=metadata" json:"metadata"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *PrometheusConfig) Reset()         { *m = PrometheusConfig{} }
func (m *PrometheusConfig) String() string { return proto.CompactTextString(m) }
func (*PrometheusConfig) ProtoMessage()    {}
func (*PrometheusConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_23b416b7a516d18d, []int{0}
}
func (m *PrometheusConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrometheusConfig.Unmarshal(m, b)
}
func (m *PrometheusConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrometheusConfig.Marshal(b, m, deterministic)
}
func (dst *PrometheusConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrometheusConfig.Merge(dst, src)
}
func (m *PrometheusConfig) XXX_Size() int {
	return xxx_messageInfo_PrometheusConfig.Size(m)
}
func (m *PrometheusConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_PrometheusConfig.DiscardUnknown(m)
}

var xxx_messageInfo_PrometheusConfig proto.InternalMessageInfo

func (m *PrometheusConfig) GetPrometheus() *types.Struct {
	if m != nil {
		return m.Prometheus
	}
	return nil
}

func (m *PrometheusConfig) GetMetadata() core.Metadata {
	if m != nil {
		return m.Metadata
	}
	return core.Metadata{}
}

func init() {
	proto.RegisterType((*PrometheusConfig)(nil), "config.prometheus.io.PrometheusConfig")
}
func (this *PrometheusConfig) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PrometheusConfig)
	if !ok {
		that2, ok := that.(PrometheusConfig)
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
	if !this.Prometheus.Equal(that1.Prometheus) {
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

func init() {
	proto.RegisterFile("github.com/solo-io/supergloo/api/external/prometheus/v1/config.proto", fileDescriptor_config_23b416b7a516d18d)
}

var fileDescriptor_config_23b416b7a516d18d = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0x5d, 0x10, 0x95, 0x08, 0x22, 0x4b, 0xd1, 0x52, 0x44, 0xc5, 0x93, 0x17, 0x67, 0x58,
	0xbd, 0x08, 0x3d, 0x59, 0xbd, 0x0a, 0x52, 0x6f, 0xde, 0xb2, 0x6b, 0x9a, 0x86, 0x66, 0x3b, 0x21,
	0x99, 0x14, 0x7d, 0x04, 0xdf, 0xc4, 0x47, 0xf1, 0x29, 0x3c, 0xf8, 0x24, 0xb2, 0xd9, 0x6d, 0xeb,
	0x41, 0x7a, 0x4a, 0x86, 0xf9, 0x87, 0xef, 0xe3, 0x17, 0x0f, 0xda, 0xf0, 0x34, 0x96, 0x50, 0x51,
	0x8d, 0x81, 0x2c, 0x5d, 0x19, 0xc2, 0x10, 0x9d, 0xf2, 0xda, 0x12, 0xa1, 0x74, 0x06, 0xd5, 0x1b,
	0x2b, 0x3f, 0x97, 0x16, 0x9d, 0xa7, 0x5a, 0xf1, 0x54, 0xc5, 0x80, 0x8b, 0x02, 0x2b, 0x9a, 0x4f,
	0x8c, 0x06, 0xe7, 0x89, 0x29, 0xef, 0xad, 0xa7, 0x2e, 0x02, 0x86, 0x06, 0x3d, 0x4d, 0x9a, 0x52,
	0x00, 0x9b, 0x5f, 0x9b, 0x1d, 0x9c, 0x68, 0x22, 0x6d, 0x15, 0xa6, 0xa9, 0x8c, 0x13, 0x0c, 0xec,
	0x63, 0xc5, 0xdd, 0xb6, 0xf8, 0xcf, 0xa7, 0x79, 0x67, 0x86, 0x93, 0xce, 0xa2, 0xc0, 0x5a, 0xb1,
	0x7c, 0x95, 0x2c, 0xdb, 0x93, 0x8b, 0x8f, 0x4c, 0x1c, 0x3e, 0xad, 0xc0, 0xf7, 0xc9, 0x24, 0x1f,
	0x0a, 0xb1, 0x96, 0xe9, 0x67, 0xe7, 0xd9, 0xe5, 0xfe, 0xf5, 0x31, 0xb4, 0x68, 0x58, 0xa2, 0xe1,
	0x39, 0xa1, 0xc7, 0x07, 0x7f, 0xbc, 0xdf, 0x6b, 0x9b, 0xdf, 0x8a, 0xbd, 0x25, 0xa3, 0xbf, 0x9b,
	0x4e, 0x8f, 0xa0, 0x22, 0xaf, 0xa0, 0x31, 0x01, 0x43, 0xf0, 0xd8, 0x6d, 0x47, 0xdb, 0x5f, 0xdf,
	0x67, 0x5b, 0xe3, 0x55, 0x7a, 0x74, 0xf7, 0xf9, 0x73, 0x9a, 0xbd, 0x0c, 0x37, 0x96, 0xea, 0x66,
	0x7a, 0x43, 0xb1, 0xe5, 0x4e, 0xb2, 0xbb, 0xf9, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x81, 0x5d, 0x3a,
	0x6e, 0x9a, 0x01, 0x00, 0x00,
}

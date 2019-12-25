// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/imtanmoy/authN/api/protos/organization.proto

package organization

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Organization struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Organization) Reset()         { *m = Organization{} }
func (m *Organization) String() string { return proto.CompactTextString(m) }
func (*Organization) ProtoMessage()    {}
func (*Organization) Descriptor() ([]byte, []int) {
	return fileDescriptor_f45fed1d089d7880, []int{0}
}

func (m *Organization) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Organization.Unmarshal(m, b)
}
func (m *Organization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Organization.Marshal(b, m, deterministic)
}
func (m *Organization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Organization.Merge(m, src)
}
func (m *Organization) XXX_Size() int {
	return xxx_messageInfo_Organization.Size(m)
}
func (m *Organization) XXX_DiscardUnknown() {
	xxx_messageInfo_Organization.DiscardUnknown(m)
}

var xxx_messageInfo_Organization proto.InternalMessageInfo

func (m *Organization) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Organization) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Organization)(nil), "organization.Organization")
}

func init() {
	proto.RegisterFile("github.com/imtanmoy/authN/api/protos/organization.proto", fileDescriptor_f45fed1d089d7880)
}

var fileDescriptor_f45fed1d089d7880 = []byte{
	// 170 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x4f, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0xcc, 0x2d, 0x49, 0xcc, 0xcb, 0xcd, 0xaf, 0xd4,
	0x4f, 0x2c, 0x2d, 0xc9, 0xa8, 0xd4, 0x4f, 0x2c, 0xc8, 0xd4, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x2f,
	0xd6, 0xcf, 0x2f, 0x4a, 0x4f, 0xcc, 0xcb, 0xac, 0x4a, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x03, 0x8b,
	0x09, 0xf1, 0x20, 0x8b, 0x29, 0x19, 0x71, 0xf1, 0xf8, 0x23, 0xf1, 0x85, 0xf8, 0xb8, 0x98, 0x32,
	0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83, 0x98, 0x32, 0x53, 0x84, 0x84, 0xb8, 0x58, 0xfc,
	0x12, 0x73, 0x53, 0x25, 0x98, 0x14, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0xa3, 0x64, 0x2e, 0x61,
	0x64, 0x3d, 0xc1, 0xa9, 0x45, 0x65, 0x99, 0xc9, 0xa9, 0x42, 0x3e, 0x5c, 0x42, 0xce, 0x45, 0xa9,
	0x89, 0x25, 0xa9, 0x28, 0x06, 0x4a, 0xe9, 0xa1, 0xb8, 0x01, 0x59, 0x4e, 0x0a, 0x8f, 0x9c, 0x93,
	0x46, 0x94, 0x1a, 0x6e, 0x1f, 0x22, 0x6b, 0x4b, 0x62, 0x03, 0xfb, 0xcb, 0x18, 0x10, 0x00, 0x00,
	0xff, 0xff, 0x8f, 0x1b, 0x60, 0xb1, 0x12, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrganizationServiceClient is the client API for OrganizationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrganizationServiceClient interface {
	CreateOrganization(ctx context.Context, in *Organization, opts ...grpc.CallOption) (*Organization, error)
}

type organizationServiceClient struct {
	cc *grpc.ClientConn
}

func NewOrganizationServiceClient(cc *grpc.ClientConn) OrganizationServiceClient {
	return &organizationServiceClient{cc}
}

func (c *organizationServiceClient) CreateOrganization(ctx context.Context, in *Organization, opts ...grpc.CallOption) (*Organization, error) {
	out := new(Organization)
	err := c.cc.Invoke(ctx, "/organization.OrganizationService/CreateOrganization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrganizationServiceServer is the server API for OrganizationService service.
type OrganizationServiceServer interface {
	CreateOrganization(context.Context, *Organization) (*Organization, error)
}

// UnimplementedOrganizationServiceServer can be embedded to have forward compatible implementations.
type UnimplementedOrganizationServiceServer struct {
}

func (*UnimplementedOrganizationServiceServer) CreateOrganization(ctx context.Context, req *Organization) (*Organization, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrganization not implemented")
}

func RegisterOrganizationServiceServer(s *grpc.Server, srv OrganizationServiceServer) {
	s.RegisterService(&_OrganizationService_serviceDesc, srv)
}

func _OrganizationService_CreateOrganization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Organization)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationServiceServer).CreateOrganization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/organization.OrganizationService/CreateOrganization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationServiceServer).CreateOrganization(ctx, req.(*Organization))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrganizationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "organization.OrganizationService",
	HandlerType: (*OrganizationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrganization",
			Handler:    _OrganizationService_CreateOrganization_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/imtanmoy/authN/api/protos/organization.proto",
}

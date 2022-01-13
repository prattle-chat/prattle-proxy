// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: group.proto

package server

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Group is used for most group operations, such as:
//
//   1. Creating a group (note: ID is ignored; IDs are generated)
//   2. Viewing group info
type Group struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Owners      []string `protobuf:"bytes,2,rep,name=owners,proto3" json:"owners,omitempty"`
	Members     []string `protobuf:"bytes,3,rep,name=members,proto3" json:"members,omitempty"`
	IsOpen      bool     `protobuf:"varint,4,opt,name=is_open,json=isOpen,proto3" json:"is_open,omitempty"`
	IsBroadcast bool     `protobuf:"varint,5,opt,name=is_broadcast,json=isBroadcast,proto3" json:"is_broadcast,omitempty"`
}

func (x *Group) Reset() {
	*x = Group{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Group) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Group) ProtoMessage() {}

func (x *Group) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Group.ProtoReflect.Descriptor instead.
func (*Group) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{0}
}

func (x *Group) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Group) GetOwners() []string {
	if x != nil {
		return x.Owners
	}
	return nil
}

func (x *Group) GetMembers() []string {
	if x != nil {
		return x.Members
	}
	return nil
}

func (x *Group) GetIsOpen() bool {
	if x != nil {
		return x.IsOpen
	}
	return false
}

func (x *Group) GetIsBroadcast() bool {
	if x != nil {
		return x.IsBroadcast
	}
	return false
}

// GroupUser is used on operations where a group ID and a user ID
// is needed, such as invitations, booting users, etc.
//
// When a GroupUser message passes between peers, the field 'For' should
// be set to the caller, as determined by the origin server; this is to
// ensure that federated peers can determine whether an operation is
// permitted.
//
// This does, of course, put a certain level of trust in the sanctity of
// the peering connection.
type GroupUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	For     string `protobuf:"bytes,3,opt,name=for,proto3" json:"for,omitempty"`
}

func (x *GroupUser) Reset() {
	*x = GroupUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupUser) ProtoMessage() {}

func (x *GroupUser) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupUser.ProtoReflect.Descriptor instead.
func (*GroupUser) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{1}
}

func (x *GroupUser) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *GroupUser) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GroupUser) GetFor() string {
	if x != nil {
		return x.For
	}
	return ""
}

var File_group_proto protoreflect.FileDescriptor

var file_group_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x85, 0x01, 0x0a, 0x05, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x77, 0x6e,
	0x65, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x17, 0x0a,
	0x07, 0x69, 0x73, 0x5f, 0x6f, 0x70, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x69, 0x73, 0x4f, 0x70, 0x65, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x62, 0x72, 0x6f,
	0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73,
	0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x22, 0x51, 0x0a, 0x09, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x55, 0x73, 0x65, 0x72, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x6f,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x66, 0x6f, 0x72, 0x32, 0xee, 0x02, 0x0a,
	0x06, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x26, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x0c, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x1a,
	0x0c, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x00, 0x12,
	0x32, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x10, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x12, 0x28, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x2e, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x0c, 0x2e,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x00, 0x12, 0x34, 0x0a,
	0x06, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x12, 0x10, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x10, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x55, 0x73, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x38,
	0x0a, 0x0a, 0x44, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x10, 0x2e, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x05, 0x4c, 0x65, 0x61, 0x76,
	0x65, 0x12, 0x10, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x55,
	0x73, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x2e, 0x5a,
	0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x61, 0x74,
	0x74, 0x6c, 0x65, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x70, 0x72, 0x61, 0x74, 0x74, 0x6c, 0x65,
	0x2d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_group_proto_rawDescOnce sync.Once
	file_group_proto_rawDescData = file_group_proto_rawDesc
)

func file_group_proto_rawDescGZIP() []byte {
	file_group_proto_rawDescOnce.Do(func() {
		file_group_proto_rawDescData = protoimpl.X.CompressGZIP(file_group_proto_rawDescData)
	})
	return file_group_proto_rawDescData
}

var file_group_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_group_proto_goTypes = []interface{}{
	(*Group)(nil),         // 0: group.Group
	(*GroupUser)(nil),     // 1: group.GroupUser
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_group_proto_depIdxs = []int32{
	0, // 0: group.Groups.Create:input_type -> group.Group
	1, // 1: group.Groups.Join:input_type -> group.GroupUser
	1, // 2: group.Groups.Info:input_type -> group.GroupUser
	1, // 3: group.Groups.Invite:input_type -> group.GroupUser
	1, // 4: group.Groups.PromoteUser:input_type -> group.GroupUser
	1, // 5: group.Groups.DemoteUser:input_type -> group.GroupUser
	1, // 6: group.Groups.Leave:input_type -> group.GroupUser
	0, // 7: group.Groups.Create:output_type -> group.Group
	2, // 8: group.Groups.Join:output_type -> google.protobuf.Empty
	0, // 9: group.Groups.Info:output_type -> group.Group
	2, // 10: group.Groups.Invite:output_type -> google.protobuf.Empty
	2, // 11: group.Groups.PromoteUser:output_type -> google.protobuf.Empty
	2, // 12: group.Groups.DemoteUser:output_type -> google.protobuf.Empty
	2, // 13: group.Groups.Leave:output_type -> google.protobuf.Empty
	7, // [7:14] is the sub-list for method output_type
	0, // [0:7] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_group_proto_init() }
func file_group_proto_init() {
	if File_group_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_group_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Group); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_group_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupUser); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_group_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_group_proto_goTypes,
		DependencyIndexes: file_group_proto_depIdxs,
		MessageInfos:      file_group_proto_msgTypes,
	}.Build()
	File_group_proto = out.File
	file_group_proto_rawDesc = nil
	file_group_proto_goTypes = nil
	file_group_proto_depIdxs = nil
}

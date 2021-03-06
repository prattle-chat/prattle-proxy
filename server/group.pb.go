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

// Group is used for creating and reading groups; with the following caveats:
//
//    1. The ID field is ignored when creating a new group; group IDs are
//       minted serverside (in much the same way user IDs are)
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

// JoinRequest is used to join a group. Clients need only set
// the value of group_id. When a prattle server relays a JoinRequest
// to a peer, the relaying server must set the value of joiner to
// the ID of the originator of the call
type JoinRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Joiner  string `protobuf:"bytes,2,opt,name=joiner,proto3" json:"joiner,omitempty"`
}

func (x *JoinRequest) Reset() {
	*x = JoinRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinRequest) ProtoMessage() {}

func (x *JoinRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use JoinRequest.ProtoReflect.Descriptor instead.
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{1}
}

func (x *JoinRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *JoinRequest) GetJoiner() string {
	if x != nil {
		return x.Joiner
	}
	return ""
}

// InfoRequest is used to retrieve a group. Clients need only set
// the value of group_id. When a prattle server relays an InfoRequest
// to a peer, the relaying server must set the value of requester to
// the ID of the originator of the call
type InfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
}

func (x *InfoRequest) Reset() {
	*x = InfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoRequest) ProtoMessage() {}

func (x *InfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoRequest.ProtoReflect.Descriptor instead.
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{2}
}

func (x *InfoRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

// LeaveRequest is used to leave a group. Clients need only set
// the value of group_id. When a prattle server relays an LeaveRequest
// to a peer, the relaying server must set the value of leaver to
// the ID of the originator of the call
type LeaveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Leaver  string `protobuf:"bytes,2,opt,name=leaver,proto3" json:"leaver,omitempty"`
}

func (x *LeaveRequest) Reset() {
	*x = LeaveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaveRequest) ProtoMessage() {}

func (x *LeaveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaveRequest.ProtoReflect.Descriptor instead.
func (*LeaveRequest) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{3}
}

func (x *LeaveRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *LeaveRequest) GetLeaver() string {
	if x != nil {
		return x.Leaver
	}
	return ""
}

// InviteRequest is used to invite a user (the invitee) to a group.
type InviteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Invitee string `protobuf:"bytes,2,opt,name=invitee,proto3" json:"invitee,omitempty"`
}

func (x *InviteRequest) Reset() {
	*x = InviteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InviteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InviteRequest) ProtoMessage() {}

func (x *InviteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InviteRequest.ProtoReflect.Descriptor instead.
func (*InviteRequest) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{4}
}

func (x *InviteRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *InviteRequest) GetInvitee() string {
	if x != nil {
		return x.Invitee
	}
	return ""
}

// PromoteRequest is used to promote a user (the promotee) within a group,
// from a normal member to a owner
//
// The value of `promoter` is ignored when an invitation request comes
// from a non-peered connection. It is the responsibility of the server
// which recieves that first, non-peered/ relayed, request to set the
// promoter field to the originator of the call
type PromoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId  string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Promotee string `protobuf:"bytes,2,opt,name=promotee,proto3" json:"promotee,omitempty"`
}

func (x *PromoteRequest) Reset() {
	*x = PromoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PromoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PromoteRequest) ProtoMessage() {}

func (x *PromoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PromoteRequest.ProtoReflect.Descriptor instead.
func (*PromoteRequest) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{5}
}

func (x *PromoteRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *PromoteRequest) GetPromotee() string {
	if x != nil {
		return x.Promotee
	}
	return ""
}

// DemoteRequest is used to demote a user (the demotee) to a group,
// from owner -> normal user, and normal user -> booted
//
// The value of `demoter` is ignored when an invitation request comes
// from a non-peered connection. It is the responsibility of the server
// which recieves that first, non-peered/ relayed, request to set the
// demoter field to the originator of the call
type DemoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Demotee string `protobuf:"bytes,2,opt,name=demotee,proto3" json:"demotee,omitempty"`
}

func (x *DemoteRequest) Reset() {
	*x = DemoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DemoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DemoteRequest) ProtoMessage() {}

func (x *DemoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DemoteRequest.ProtoReflect.Descriptor instead.
func (*DemoteRequest) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{6}
}

func (x *DemoteRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *DemoteRequest) GetDemotee() string {
	if x != nil {
		return x.Demotee
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
	0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x22, 0x40, 0x0a, 0x0b, 0x4a, 0x6f, 0x69,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6a, 0x6f, 0x69, 0x6e, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6a, 0x6f, 0x69, 0x6e, 0x65, 0x72, 0x22, 0x28, 0x0a, 0x0b, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x49, 0x64, 0x22, 0x41, 0x0a, 0x0c, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6c, 0x65, 0x61, 0x76, 0x65, 0x72, 0x22, 0x44, 0x0a, 0x0d, 0x49, 0x6e, 0x76, 0x69,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x65, 0x22, 0x47,
	0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70,
	0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x65, 0x22, 0x44, 0x0a, 0x0d, 0x44, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x65, 0x32, 0x82, 0x03,
	0x0a, 0x06, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x26, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x12, 0x0c, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x1a, 0x0c, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x00,
	0x12, 0x34, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x12, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2a, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12,
	0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x22, 0x00, 0x12, 0x38, 0x0a, 0x06, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x12, 0x14, 0x2e, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x2e, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0b,
	0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x2e, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x0a,
	0x44, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x2e, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x2e, 0x44, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x05, 0x4c, 0x65,
	0x61, 0x76, 0x65, 0x12, 0x13, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x4c, 0x65, 0x61, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x70, 0x72, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x70, 0x72,
	0x61, 0x74, 0x74, 0x6c, 0x65, 0x2d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_group_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_group_proto_goTypes = []interface{}{
	(*Group)(nil),          // 0: group.Group
	(*JoinRequest)(nil),    // 1: group.JoinRequest
	(*InfoRequest)(nil),    // 2: group.InfoRequest
	(*LeaveRequest)(nil),   // 3: group.LeaveRequest
	(*InviteRequest)(nil),  // 4: group.InviteRequest
	(*PromoteRequest)(nil), // 5: group.PromoteRequest
	(*DemoteRequest)(nil),  // 6: group.DemoteRequest
	(*emptypb.Empty)(nil),  // 7: google.protobuf.Empty
}
var file_group_proto_depIdxs = []int32{
	0, // 0: group.Groups.Create:input_type -> group.Group
	1, // 1: group.Groups.Join:input_type -> group.JoinRequest
	2, // 2: group.Groups.Info:input_type -> group.InfoRequest
	4, // 3: group.Groups.Invite:input_type -> group.InviteRequest
	5, // 4: group.Groups.PromoteUser:input_type -> group.PromoteRequest
	6, // 5: group.Groups.DemoteUser:input_type -> group.DemoteRequest
	3, // 6: group.Groups.Leave:input_type -> group.LeaveRequest
	0, // 7: group.Groups.Create:output_type -> group.Group
	7, // 8: group.Groups.Join:output_type -> google.protobuf.Empty
	0, // 9: group.Groups.Info:output_type -> group.Group
	7, // 10: group.Groups.Invite:output_type -> google.protobuf.Empty
	7, // 11: group.Groups.PromoteUser:output_type -> google.protobuf.Empty
	7, // 12: group.Groups.DemoteUser:output_type -> google.protobuf.Empty
	7, // 13: group.Groups.Leave:output_type -> google.protobuf.Empty
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
			switch v := v.(*JoinRequest); i {
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
		file_group_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoRequest); i {
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
		file_group_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaveRequest); i {
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
		file_group_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InviteRequest); i {
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
		file_group_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PromoteRequest); i {
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
		file_group_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DemoteRequest); i {
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
			NumMessages:   7,
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

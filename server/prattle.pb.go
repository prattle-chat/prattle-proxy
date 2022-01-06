// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: prattle.proto

package server

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// SignupRequest contains a password to be used during Signup
//
// This message could be merged with OTPAndKey, but it makes a certain
// amount of sense to keep this small, explicit, and less complex by
// making it only used in a single place
type SignupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *SignupRequest) Reset() {
	*x = SignupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prattle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignupRequest) ProtoMessage() {}

func (x *SignupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_prattle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignupRequest.ProtoReflect.Descriptor instead.
func (*SignupRequest) Descriptor() ([]byte, []int) {
	return file_prattle_proto_rawDescGZIP(), []int{0}
}

func (x *SignupRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// SignupResponse contains the information necessary to connect to a Proxy;
// namely: the user's new ID (including domain name information), and a
// value which can seed an OTP app/authenticator/etc. to act as a password.
type SignupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Seed int64  `protobuf:"varint,2,opt,name=seed,proto3" json:"seed,omitempty"`
}

func (x *SignupResponse) Reset() {
	*x = SignupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prattle_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignupResponse) ProtoMessage() {}

func (x *SignupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_prattle_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignupResponse.ProtoReflect.Descriptor instead.
func (*SignupResponse) Descriptor() ([]byte, []int) {
	return file_prattle_proto_rawDescGZIP(), []int{1}
}

func (x *SignupResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SignupResponse) GetSeed() int64 {
	if x != nil {
		return x.Seed
	}
	return 0
}

// OTPAndKey encapsulates an OTP value and a key.
// It is used:
//
//   1. On finalise: where value is a valid OTP value, and key is a user ID
//   2. On token request: where value is a valid OTP value, and key is password
//   3. On returning a minted token (as a response from above): where value is empty,
//      and key is a token to auth with
//   4. On looking up a public key: where key is a user key
type OTPAndKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64  `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	Key   string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *OTPAndKey) Reset() {
	*x = OTPAndKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prattle_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OTPAndKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OTPAndKey) ProtoMessage() {}

func (x *OTPAndKey) ProtoReflect() protoreflect.Message {
	mi := &file_prattle_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OTPAndKey.ProtoReflect.Descriptor instead.
func (*OTPAndKey) Descriptor() ([]byte, []int) {
	return file_prattle_proto_rawDescGZIP(), []int{2}
}

func (x *OTPAndKey) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *OTPAndKey) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// MessageWrappper wraps an encoded/ encrypted message to be forwarded to a recipient
//
// Encoded is expected to be formed by taking a 'Message' and encypting it with the
// recipient's public key. Because of this, that message is where all of the important
// and useful metadatas are stored, such as sender
type MessageWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Recipient string `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Encoded   []byte `protobuf:"bytes,2,opt,name=encoded,proto3" json:"encoded,omitempty"`
}

func (x *MessageWrapper) Reset() {
	*x = MessageWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prattle_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageWrapper) ProtoMessage() {}

func (x *MessageWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_prattle_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageWrapper.ProtoReflect.Descriptor instead.
func (*MessageWrapper) Descriptor() ([]byte, []int) {
	return file_prattle_proto_rawDescGZIP(), []int{3}
}

func (x *MessageWrapper) GetRecipient() string {
	if x != nil {
		return x.Recipient
	}
	return ""
}

func (x *MessageWrapper) GetEncoded() []byte {
	if x != nil {
		return x.Encoded
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// recipient in this context is used to determine whether a message was sent to a
	// user directly, or to a group
	Recipient string                 `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Sender    string                 `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Datetime  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=datetime,proto3" json:"datetime,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prattle_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_prattle_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_prattle_proto_rawDescGZIP(), []int{4}
}

func (x *Message) GetRecipient() string {
	if x != nil {
		return x.Recipient
	}
	return ""
}

func (x *Message) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *Message) GetDatetime() *timestamppb.Timestamp {
	if x != nil {
		return x.Datetime
	}
	return nil
}

var File_prattle_proto protoreflect.FileDescriptor

var file_prattle_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2b, 0x0a,
	0x0d, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x34, 0x0a, 0x0e, 0x53, 0x69,
	0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x65, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x65, 0x65, 0x64,
	0x22, 0x33, 0x0a, 0x09, 0x4f, 0x54, 0x50, 0x41, 0x6e, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x48, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x63, 0x69,
	0x70, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x22,
	0x77, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x12, 0x36, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08,
	0x64, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x32, 0xd3, 0x02, 0x0a, 0x05, 0x50, 0x72, 0x6f,
	0x78, 0x79, 0x12, 0x2b, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x0e, 0x2e, 0x53,
	0x69, 0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x53,
	0x69, 0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x30, 0x0a, 0x08, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x73, 0x65, 0x12, 0x0a, 0x2e, 0x4f, 0x54,
	0x50, 0x41, 0x6e, 0x64, 0x4b, 0x65, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x21, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x0a, 0x2e, 0x4f, 0x54, 0x50,
	0x41, 0x6e, 0x64, 0x4b, 0x65, 0x79, 0x1a, 0x0a, 0x2e, 0x4f, 0x54, 0x50, 0x41, 0x6e, 0x64, 0x4b,
	0x65, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x0c, 0x41, 0x64, 0x64, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x4b, 0x65, 0x79, 0x12, 0x0a, 0x2e, 0x4f, 0x54, 0x50, 0x41, 0x6e, 0x64, 0x4b, 0x65, 0x79,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x09, 0x53, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x0f, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x22, 0x00, 0x30, 0x01, 0x12, 0x25, 0x0a, 0x09, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65,
	0x79, 0x12, 0x0a, 0x2e, 0x4f, 0x54, 0x50, 0x41, 0x6e, 0x64, 0x4b, 0x65, 0x79, 0x1a, 0x0a, 0x2e,
	0x4f, 0x54, 0x50, 0x41, 0x6e, 0x64, 0x4b, 0x65, 0x79, 0x22, 0x00, 0x12, 0x31, 0x0a, 0x04, 0x53,
	0x65, 0x6e, 0x64, 0x12, 0x0f, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x2e,
	0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x61,
	0x74, 0x74, 0x6c, 0x65, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x70, 0x72, 0x61, 0x74, 0x74, 0x6c,
	0x65, 0x2d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_prattle_proto_rawDescOnce sync.Once
	file_prattle_proto_rawDescData = file_prattle_proto_rawDesc
)

func file_prattle_proto_rawDescGZIP() []byte {
	file_prattle_proto_rawDescOnce.Do(func() {
		file_prattle_proto_rawDescData = protoimpl.X.CompressGZIP(file_prattle_proto_rawDescData)
	})
	return file_prattle_proto_rawDescData
}

var file_prattle_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_prattle_proto_goTypes = []interface{}{
	(*SignupRequest)(nil),         // 0: SignupRequest
	(*SignupResponse)(nil),        // 1: SignupResponse
	(*OTPAndKey)(nil),             // 2: OTPAndKey
	(*MessageWrapper)(nil),        // 3: MessageWrapper
	(*Message)(nil),               // 4: Message
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 6: google.protobuf.Empty
}
var file_prattle_proto_depIdxs = []int32{
	5, // 0: Message.datetime:type_name -> google.protobuf.Timestamp
	0, // 1: Proxy.Signup:input_type -> SignupRequest
	2, // 2: Proxy.Finalise:input_type -> OTPAndKey
	2, // 3: Proxy.Token:input_type -> OTPAndKey
	2, // 4: Proxy.AddPublicKey:input_type -> OTPAndKey
	6, // 5: Proxy.Subscribe:input_type -> google.protobuf.Empty
	2, // 6: Proxy.PublicKey:input_type -> OTPAndKey
	3, // 7: Proxy.Send:input_type -> MessageWrapper
	1, // 8: Proxy.Signup:output_type -> SignupResponse
	6, // 9: Proxy.Finalise:output_type -> google.protobuf.Empty
	2, // 10: Proxy.Token:output_type -> OTPAndKey
	6, // 11: Proxy.AddPublicKey:output_type -> google.protobuf.Empty
	3, // 12: Proxy.Subscribe:output_type -> MessageWrapper
	2, // 13: Proxy.PublicKey:output_type -> OTPAndKey
	6, // 14: Proxy.Send:output_type -> google.protobuf.Empty
	8, // [8:15] is the sub-list for method output_type
	1, // [1:8] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_prattle_proto_init() }
func file_prattle_proto_init() {
	if File_prattle_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_prattle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignupRequest); i {
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
		file_prattle_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignupResponse); i {
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
		file_prattle_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OTPAndKey); i {
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
		file_prattle_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageWrapper); i {
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
		file_prattle_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_prattle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_prattle_proto_goTypes,
		DependencyIndexes: file_prattle_proto_depIdxs,
		MessageInfos:      file_prattle_proto_msgTypes,
	}.Build()
	File_prattle_proto = out.File
	file_prattle_proto_rawDesc = nil
	file_prattle_proto_goTypes = nil
	file_prattle_proto_depIdxs = nil
}

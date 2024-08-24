// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: counter.proto

package counter

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetUnreadCountTotalV1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetUnreadCountTotalV1Request) Reset() {
	*x = GetUnreadCountTotalV1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_counter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUnreadCountTotalV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUnreadCountTotalV1Request) ProtoMessage() {}

func (x *GetUnreadCountTotalV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_counter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUnreadCountTotalV1Request.ProtoReflect.Descriptor instead.
func (*GetUnreadCountTotalV1Request) Descriptor() ([]byte, []int) {
	return file_counter_proto_rawDescGZIP(), []int{0}
}

func (x *GetUnreadCountTotalV1Request) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetUnreadCountTotalV1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *GetUnreadCountTotalV1Response) Reset() {
	*x = GetUnreadCountTotalV1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_counter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUnreadCountTotalV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUnreadCountTotalV1Response) ProtoMessage() {}

func (x *GetUnreadCountTotalV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_counter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUnreadCountTotalV1Response.ProtoReflect.Descriptor instead.
func (*GetUnreadCountTotalV1Response) Descriptor() ([]byte, []int) {
	return file_counter_proto_rawDescGZIP(), []int{1}
}

func (x *GetUnreadCountTotalV1Response) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type GetUnreadCountV1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentUserId string `protobuf:"bytes,1,opt,name=current_user_id,json=currentUserId,proto3" json:"current_user_id,omitempty"`
	ChatUserId    string `protobuf:"bytes,2,opt,name=chat_user_id,json=chatUserId,proto3" json:"chat_user_id,omitempty"`
}

func (x *GetUnreadCountV1Request) Reset() {
	*x = GetUnreadCountV1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_counter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUnreadCountV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUnreadCountV1Request) ProtoMessage() {}

func (x *GetUnreadCountV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_counter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUnreadCountV1Request.ProtoReflect.Descriptor instead.
func (*GetUnreadCountV1Request) Descriptor() ([]byte, []int) {
	return file_counter_proto_rawDescGZIP(), []int{2}
}

func (x *GetUnreadCountV1Request) GetCurrentUserId() string {
	if x != nil {
		return x.CurrentUserId
	}
	return ""
}

func (x *GetUnreadCountV1Request) GetChatUserId() string {
	if x != nil {
		return x.ChatUserId
	}
	return ""
}

type GetUnreadCountV1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *GetUnreadCountV1Response) Reset() {
	*x = GetUnreadCountV1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_counter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUnreadCountV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUnreadCountV1Response) ProtoMessage() {}

func (x *GetUnreadCountV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_counter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUnreadCountV1Response.ProtoReflect.Descriptor instead.
func (*GetUnreadCountV1Response) Descriptor() ([]byte, []int) {
	return file_counter_proto_rawDescGZIP(), []int{3}
}

func (x *GetUnreadCountV1Response) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type MarkMessagesAsReadV1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	MessageIds []int64 `protobuf:"varint,2,rep,packed,name=message_ids,json=messageIds,proto3" json:"message_ids,omitempty"`
}

func (x *MarkMessagesAsReadV1Request) Reset() {
	*x = MarkMessagesAsReadV1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_counter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarkMessagesAsReadV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkMessagesAsReadV1Request) ProtoMessage() {}

func (x *MarkMessagesAsReadV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_counter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkMessagesAsReadV1Request.ProtoReflect.Descriptor instead.
func (*MarkMessagesAsReadV1Request) Descriptor() ([]byte, []int) {
	return file_counter_proto_rawDescGZIP(), []int{4}
}

func (x *MarkMessagesAsReadV1Request) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *MarkMessagesAsReadV1Request) GetMessageIds() []int64 {
	if x != nil {
		return x.MessageIds
	}
	return nil
}

type MarkMessagesAsReadV1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MarkMessagesAsReadV1Response) Reset() {
	*x = MarkMessagesAsReadV1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_counter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarkMessagesAsReadV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkMessagesAsReadV1Response) ProtoMessage() {}

func (x *MarkMessagesAsReadV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_counter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkMessagesAsReadV1Response.ProtoReflect.Descriptor instead.
func (*MarkMessagesAsReadV1Response) Descriptor() ([]byte, []int) {
	return file_counter_proto_rawDescGZIP(), []int{5}
}

var File_counter_proto protoreflect.FileDescriptor

var file_counter_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x22, 0x37, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x55,
	0x6e, 0x72, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x56,
	0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x35, 0x0a, 0x1d, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x63, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x55,
	0x6e, 0x72, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c, 0x63,
	0x68, 0x61, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x63, 0x68, 0x61, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x30, 0x0a,
	0x18, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x56,
	0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x57, 0x0a, 0x1b, 0x4d, 0x61, 0x72, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x41,
	0x73, 0x52, 0x65, 0x61, 0x64, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0a, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x73, 0x22, 0x1e, 0x0a, 0x1c, 0x4d, 0x61, 0x72, 0x6b,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x41, 0x73, 0x52, 0x65, 0x61, 0x64, 0x56, 0x31,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xb6, 0x02, 0x0a, 0x0e, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x66, 0x0a, 0x15, 0x47,
	0x65, 0x74, 0x55, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x6f, 0x74,
	0x61, 0x6c, 0x56, 0x31, 0x12, 0x25, 0x2e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x47,
	0x65, 0x74, 0x55, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x6f, 0x74,
	0x61, 0x6c, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x57, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x72, 0x65, 0x61, 0x64,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x31, 0x12, 0x20, 0x2e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x63, 0x0a, 0x14,
	0x4d, 0x61, 0x72, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x41, 0x73, 0x52, 0x65,
	0x61, 0x64, 0x56, 0x31, 0x12, 0x24, 0x2e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x4d,
	0x61, 0x72, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x41, 0x73, 0x52, 0x65, 0x61,
	0x64, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x41, 0x73, 0x52, 0x65, 0x61, 0x64, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6f, 0x72, 0x6f, 0x63, 0x68, 0x69, 0x2d, 0x6b, 0x65, 0x79, 0x64, 0x72, 0x65, 0x61, 0x6d, 0x2f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_counter_proto_rawDescOnce sync.Once
	file_counter_proto_rawDescData = file_counter_proto_rawDesc
)

func file_counter_proto_rawDescGZIP() []byte {
	file_counter_proto_rawDescOnce.Do(func() {
		file_counter_proto_rawDescData = protoimpl.X.CompressGZIP(file_counter_proto_rawDescData)
	})
	return file_counter_proto_rawDescData
}

var file_counter_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_counter_proto_goTypes = []interface{}{
	(*GetUnreadCountTotalV1Request)(nil),  // 0: counter.GetUnreadCountTotalV1Request
	(*GetUnreadCountTotalV1Response)(nil), // 1: counter.GetUnreadCountTotalV1Response
	(*GetUnreadCountV1Request)(nil),       // 2: counter.GetUnreadCountV1Request
	(*GetUnreadCountV1Response)(nil),      // 3: counter.GetUnreadCountV1Response
	(*MarkMessagesAsReadV1Request)(nil),   // 4: counter.MarkMessagesAsReadV1Request
	(*MarkMessagesAsReadV1Response)(nil),  // 5: counter.MarkMessagesAsReadV1Response
}
var file_counter_proto_depIdxs = []int32{
	0, // 0: counter.CounterService.GetUnreadCountTotalV1:input_type -> counter.GetUnreadCountTotalV1Request
	2, // 1: counter.CounterService.GetUnreadCountV1:input_type -> counter.GetUnreadCountV1Request
	4, // 2: counter.CounterService.MarkMessagesAsReadV1:input_type -> counter.MarkMessagesAsReadV1Request
	1, // 3: counter.CounterService.GetUnreadCountTotalV1:output_type -> counter.GetUnreadCountTotalV1Response
	3, // 4: counter.CounterService.GetUnreadCountV1:output_type -> counter.GetUnreadCountV1Response
	5, // 5: counter.CounterService.MarkMessagesAsReadV1:output_type -> counter.MarkMessagesAsReadV1Response
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_counter_proto_init() }
func file_counter_proto_init() {
	if File_counter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_counter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUnreadCountTotalV1Request); i {
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
		file_counter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUnreadCountTotalV1Response); i {
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
		file_counter_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUnreadCountV1Request); i {
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
		file_counter_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUnreadCountV1Response); i {
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
		file_counter_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarkMessagesAsReadV1Request); i {
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
		file_counter_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarkMessagesAsReadV1Response); i {
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
			RawDescriptor: file_counter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_counter_proto_goTypes,
		DependencyIndexes: file_counter_proto_depIdxs,
		MessageInfos:      file_counter_proto_msgTypes,
	}.Build()
	File_counter_proto = out.File
	file_counter_proto_rawDesc = nil
	file_counter_proto_goTypes = nil
	file_counter_proto_depIdxs = nil
}
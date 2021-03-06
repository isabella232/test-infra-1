// Code generated by protoc-gen-go. DO NOT EDIT.
// source: scheduling/v1/scheduling_service.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	testing "github.com/grpc/test-infra/proto/grpc/testing"
	longrunning "google.golang.org/genproto/googleapis/longrunning"
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

type Component_Kind int32

const (
	// No kind was specified, so the service will attempt to infer one.
	//
	// This inference works as follows:
	//   1. If this field belongs to the driver_container of the
	//      StartTestSessionRequest, this kind is inferred to be a driver.
	//   2. Otherwise, this field is inferred to be a client worker.
	Component_UNSPECIFIED Component_Kind = 0
	// Orchestrates the test amongst workers and reports the result.
	Component_DRIVER Component_Kind = 1
	// Processes incoming requests from client containers.
	Component_SERVER Component_Kind = 2
	// Sends outgoing requests to a server container.
	Component_CLIENT Component_Kind = 3
)

var Component_Kind_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "DRIVER",
	2: "SERVER",
	3: "CLIENT",
}

var Component_Kind_value = map[string]int32{
	"UNSPECIFIED": 0,
	"DRIVER":      1,
	"SERVER":      2,
	"CLIENT":      3,
}

func (x Component_Kind) String() string {
	return proto.EnumName(Component_Kind_name, int32(x))
}

func (Component_Kind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cbfd1c3e45490e21, []int{0, 0}
}

type Event_Kind int32

const (
	// An infrastructure problem. This is outside of the test runner's control,
	// and a bug should be filed if it is encountered.
	Event_INTERNAL_ERROR Event_Kind = 0
	// Waiting for an executor to process the test session.
	Event_QUEUE Event_Kind = 1
	// An executor has been assigned to provision and monitor the session;
	// however, work has not yet begun on this event's subject.
	Event_ACCEPT Event_Kind = 2
	// Reserving and configuring resources for the event's subject.
	Event_PROVISION Event_Kind = 3
	// Resources are responding with a healthy signal. However, this is not
	// indicative of a running test.
	Event_RUN Event_Kind = 4
	// Resources terminated or released as expected. It does not indicate that
	// the test was successful or the result was recorded.
	Event_DONE Event_Kind = 5
	// An irrecoverable error has caused the subject of the event to terminate.
	Event_ERROR Event_Kind = 6
)

var Event_Kind_name = map[int32]string{
	0: "INTERNAL_ERROR",
	1: "QUEUE",
	2: "ACCEPT",
	3: "PROVISION",
	4: "RUN",
	5: "DONE",
	6: "ERROR",
}

var Event_Kind_value = map[string]int32{
	"INTERNAL_ERROR": 0,
	"QUEUE":          1,
	"ACCEPT":         2,
	"PROVISION":      3,
	"RUN":            4,
	"DONE":           5,
	"ERROR":          6,
}

func (x Event_Kind) String() string {
	return proto.EnumName(Event_Kind_name, int32(x))
}

func (Event_Kind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cbfd1c3e45490e21, []int{3, 0}
}

type Component struct {
	// The name and version of the container image as it appears in a registry.
	//
	// When using GCR, this will likely be a fully-qualified domain name and path
	// plus a tag or sha256 hash. For example, a Java worker image may be
	// similar to one of the following:
	//
	//   - gcr.io/grpc-testing/java_worker:v3.2.19
	//   - gcr.io/grpc-testing/java_worker@sha256:82b6360a84c19f23ed7ee9...
	//
	// The "latest" tag is automatically set by the registry, and there
	// are no guarantees that it will point to a specific image. It should be
	// avoided.
	ContainerImage string `protobuf:"bytes,1,opt,name=container_image,json=containerImage,proto3" json:"container_image,omitempty"`
	// The kind which assigns the responsibilities of this container.
	Kind Component_Kind `protobuf:"varint,2,opt,name=kind,proto3,enum=grpc.testing.benchmarking.scheduling.v1.Component_Kind" json:"kind,omitempty"`
	// The pool where the component should run. This must match the "pool" label
	// on a node to be scheduled on it.
	Pool                 string   `protobuf:"bytes,3,opt,name=pool,proto3" json:"pool,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Component) Reset()         { *m = Component{} }
func (m *Component) String() string { return proto.CompactTextString(m) }
func (*Component) ProtoMessage()    {}
func (*Component) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbfd1c3e45490e21, []int{0}
}

func (m *Component) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Component.Unmarshal(m, b)
}
func (m *Component) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Component.Marshal(b, m, deterministic)
}
func (m *Component) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Component.Merge(m, src)
}
func (m *Component) XXX_Size() int {
	return xxx_messageInfo_Component.Size(m)
}
func (m *Component) XXX_DiscardUnknown() {
	xxx_messageInfo_Component.DiscardUnknown(m)
}

var xxx_messageInfo_Component proto.InternalMessageInfo

func (m *Component) GetContainerImage() string {
	if m != nil {
		return m.ContainerImage
	}
	return ""
}

func (m *Component) GetKind() Component_Kind {
	if m != nil {
		return m.Kind
	}
	return Component_UNSPECIFIED
}

func (m *Component) GetPool() string {
	if m != nil {
		return m.Pool
	}
	return ""
}

type StartTestSessionRequest struct {
	// The test scenario to run as a part of the test session.
	Scenario *testing.Scenario `protobuf:"bytes,1,opt,name=scenario,proto3" json:"scenario,omitempty"`
	// The component which orchestrates running the test scenarios amongst
	// workers.
	Driver *Component `protobuf:"bytes,2,opt,name=driver,proto3" json:"driver,omitempty"`
	// The components that are required as part of the test. Normally, this will
	// involve at least one server and a number of clients.
	Workers []*Component `protobuf:"bytes,3,rep,name=workers,proto3" json:"workers,omitempty"`
	// An optional location where the result should be written.
	//
	// Types that are valid to be assigned to ResultDestination:
	//	*StartTestSessionRequest_BqResultsTable
	//	*StartTestSessionRequest_GcsResultsFileUri
	ResultDestination    isStartTestSessionRequest_ResultDestination `protobuf_oneof:"result_destination"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *StartTestSessionRequest) Reset()         { *m = StartTestSessionRequest{} }
func (m *StartTestSessionRequest) String() string { return proto.CompactTextString(m) }
func (*StartTestSessionRequest) ProtoMessage()    {}
func (*StartTestSessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbfd1c3e45490e21, []int{1}
}

func (m *StartTestSessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartTestSessionRequest.Unmarshal(m, b)
}
func (m *StartTestSessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartTestSessionRequest.Marshal(b, m, deterministic)
}
func (m *StartTestSessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartTestSessionRequest.Merge(m, src)
}
func (m *StartTestSessionRequest) XXX_Size() int {
	return xxx_messageInfo_StartTestSessionRequest.Size(m)
}
func (m *StartTestSessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StartTestSessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StartTestSessionRequest proto.InternalMessageInfo

func (m *StartTestSessionRequest) GetScenario() *testing.Scenario {
	if m != nil {
		return m.Scenario
	}
	return nil
}

func (m *StartTestSessionRequest) GetDriver() *Component {
	if m != nil {
		return m.Driver
	}
	return nil
}

func (m *StartTestSessionRequest) GetWorkers() []*Component {
	if m != nil {
		return m.Workers
	}
	return nil
}

type isStartTestSessionRequest_ResultDestination interface {
	isStartTestSessionRequest_ResultDestination()
}

type StartTestSessionRequest_BqResultsTable struct {
	BqResultsTable string `protobuf:"bytes,4,opt,name=bq_results_table,json=bqResultsTable,proto3,oneof"`
}

type StartTestSessionRequest_GcsResultsFileUri struct {
	GcsResultsFileUri string `protobuf:"bytes,5,opt,name=gcs_results_file_uri,json=gcsResultsFileUri,proto3,oneof"`
}

func (*StartTestSessionRequest_BqResultsTable) isStartTestSessionRequest_ResultDestination() {}

func (*StartTestSessionRequest_GcsResultsFileUri) isStartTestSessionRequest_ResultDestination() {}

func (m *StartTestSessionRequest) GetResultDestination() isStartTestSessionRequest_ResultDestination {
	if m != nil {
		return m.ResultDestination
	}
	return nil
}

func (m *StartTestSessionRequest) GetBqResultsTable() string {
	if x, ok := m.GetResultDestination().(*StartTestSessionRequest_BqResultsTable); ok {
		return x.BqResultsTable
	}
	return ""
}

func (m *StartTestSessionRequest) GetGcsResultsFileUri() string {
	if x, ok := m.GetResultDestination().(*StartTestSessionRequest_GcsResultsFileUri); ok {
		return x.GcsResultsFileUri
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*StartTestSessionRequest) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*StartTestSessionRequest_BqResultsTable)(nil),
		(*StartTestSessionRequest_GcsResultsFileUri)(nil),
	}
}

type TestSessionResult struct {
	// The logs of the driver, which contain legible descriptions of the errors
	// and any session results.
	DriverLogs []byte `protobuf:"bytes,2,opt,name=driver_logs,json=driverLogs,proto3" json:"driver_logs,omitempty"`
	// The amount of time that this session lived, including all actions from
	// scheduling to termination.
	TimeElapsed          *duration.Duration `protobuf:"bytes,3,opt,name=time_elapsed,json=timeElapsed,proto3" json:"time_elapsed,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *TestSessionResult) Reset()         { *m = TestSessionResult{} }
func (m *TestSessionResult) String() string { return proto.CompactTextString(m) }
func (*TestSessionResult) ProtoMessage()    {}
func (*TestSessionResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbfd1c3e45490e21, []int{2}
}

func (m *TestSessionResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestSessionResult.Unmarshal(m, b)
}
func (m *TestSessionResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestSessionResult.Marshal(b, m, deterministic)
}
func (m *TestSessionResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestSessionResult.Merge(m, src)
}
func (m *TestSessionResult) XXX_Size() int {
	return xxx_messageInfo_TestSessionResult.Size(m)
}
func (m *TestSessionResult) XXX_DiscardUnknown() {
	xxx_messageInfo_TestSessionResult.DiscardUnknown(m)
}

var xxx_messageInfo_TestSessionResult proto.InternalMessageInfo

func (m *TestSessionResult) GetDriverLogs() []byte {
	if m != nil {
		return m.DriverLogs
	}
	return nil
}

func (m *TestSessionResult) GetTimeElapsed() *duration.Duration {
	if m != nil {
		return m.TimeElapsed
	}
	return nil
}

type Event struct {
	// The name of the subject of the event. This may be the name of a test
	// session or the name of one of its components.
	Subject string `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	// The kind of event.
	Kind Event_Kind `protobuf:"varint,2,opt,name=kind,proto3,enum=grpc.testing.benchmarking.scheduling.v1.Event_Kind" json:"kind,omitempty"`
	// The point in time when the event was noticed.
	Time *timestamp.Timestamp `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
	// A string containing a description of the event, if applicable.
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbfd1c3e45490e21, []int{3}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Event) GetKind() Event_Kind {
	if m != nil {
		return m.Kind
	}
	return Event_INTERNAL_ERROR
}

func (m *Event) GetTime() *timestamp.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *Event) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type TestSessionMetadata struct {
	// The latest event for the test session as a whole, assuming it has not been
	// terminated.
	LatestEvent *Event `protobuf:"bytes,1,opt,name=latest_event,json=latestEvent,proto3" json:"latest_event,omitempty"`
	// The version of the service that is processing this request.
	ServiceVersion string `protobuf:"bytes,2,opt,name=service_version,json=serviceVersion,proto3" json:"service_version,omitempty"`
	// The time that the server acknowledged the request to create a test session.
	CreateTime           *timestamp.Timestamp `protobuf:"bytes,3,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *TestSessionMetadata) Reset()         { *m = TestSessionMetadata{} }
func (m *TestSessionMetadata) String() string { return proto.CompactTextString(m) }
func (*TestSessionMetadata) ProtoMessage()    {}
func (*TestSessionMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbfd1c3e45490e21, []int{4}
}

func (m *TestSessionMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestSessionMetadata.Unmarshal(m, b)
}
func (m *TestSessionMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestSessionMetadata.Marshal(b, m, deterministic)
}
func (m *TestSessionMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestSessionMetadata.Merge(m, src)
}
func (m *TestSessionMetadata) XXX_Size() int {
	return xxx_messageInfo_TestSessionMetadata.Size(m)
}
func (m *TestSessionMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_TestSessionMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_TestSessionMetadata proto.InternalMessageInfo

func (m *TestSessionMetadata) GetLatestEvent() *Event {
	if m != nil {
		return m.LatestEvent
	}
	return nil
}

func (m *TestSessionMetadata) GetServiceVersion() string {
	if m != nil {
		return m.ServiceVersion
	}
	return ""
}

func (m *TestSessionMetadata) GetCreateTime() *timestamp.Timestamp {
	if m != nil {
		return m.CreateTime
	}
	return nil
}

func init() {
	proto.RegisterEnum("grpc.testing.benchmarking.scheduling.v1.Component_Kind", Component_Kind_name, Component_Kind_value)
	proto.RegisterEnum("grpc.testing.benchmarking.scheduling.v1.Event_Kind", Event_Kind_name, Event_Kind_value)
	proto.RegisterType((*Component)(nil), "grpc.testing.benchmarking.scheduling.v1.Component")
	proto.RegisterType((*StartTestSessionRequest)(nil), "grpc.testing.benchmarking.scheduling.v1.StartTestSessionRequest")
	proto.RegisterType((*TestSessionResult)(nil), "grpc.testing.benchmarking.scheduling.v1.TestSessionResult")
	proto.RegisterType((*Event)(nil), "grpc.testing.benchmarking.scheduling.v1.Event")
	proto.RegisterType((*TestSessionMetadata)(nil), "grpc.testing.benchmarking.scheduling.v1.TestSessionMetadata")
}

func init() {
	proto.RegisterFile("scheduling/v1/scheduling_service.proto", fileDescriptor_cbfd1c3e45490e21)
}

var fileDescriptor_cbfd1c3e45490e21 = []byte{
	// 789 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x4f, 0x73, 0xe3, 0x34,
	0x14, 0xaf, 0x93, 0xb4, 0xdd, 0x3c, 0x97, 0xd4, 0x15, 0x3b, 0x10, 0x32, 0x03, 0xdb, 0x31, 0x33,
	0xb4, 0x03, 0x33, 0xf6, 0x34, 0x3d, 0x70, 0x58, 0x0e, 0x74, 0x53, 0x2f, 0x98, 0x0d, 0x49, 0x57,
	0x4e, 0x7a, 0xe0, 0x80, 0xc7, 0x7f, 0xb4, 0xae, 0xa8, 0x63, 0xb9, 0x92, 0x1c, 0xbe, 0x0c, 0x07,
	0x3e, 0x05, 0x1f, 0x82, 0x03, 0x07, 0xce, 0x7c, 0x18, 0x46, 0x96, 0x93, 0x4d, 0x96, 0x65, 0x26,
	0xbb, 0x37, 0xe9, 0xbd, 0x9f, 0x7e, 0x7a, 0xef, 0xa7, 0x9f, 0x9f, 0xe1, 0x0b, 0x91, 0xdc, 0x91,
	0xb4, 0xca, 0x69, 0x91, 0xb9, 0xcb, 0x0b, 0xf7, 0xf5, 0x2e, 0x14, 0x84, 0x2f, 0x69, 0x42, 0x9c,
	0x92, 0x33, 0xc9, 0xd0, 0x59, 0xc6, 0xcb, 0xc4, 0x91, 0x44, 0x48, 0x5a, 0x64, 0x4e, 0x4c, 0x8a,
	0xe4, 0x6e, 0x11, 0xf1, 0x7b, 0xb5, 0x79, 0x7d, 0xc6, 0x59, 0x5e, 0x0c, 0x3e, 0xcf, 0x18, 0xcb,
	0x72, 0xe2, 0xe6, 0xac, 0xc8, 0x78, 0x55, 0x14, 0x8a, 0x98, 0x95, 0x84, 0x47, 0x92, 0xb2, 0x42,
	0x68, 0xb6, 0xc1, 0x67, 0x0d, 0xa8, 0xde, 0xc5, 0xd5, 0x2b, 0x37, 0xad, 0x34, 0xa0, 0xc9, 0x3f,
	0x79, 0x33, 0x2f, 0xe9, 0x82, 0x08, 0x19, 0x2d, 0xca, 0x06, 0x30, 0x50, 0xe5, 0xb8, 0x4d, 0x39,
	0x6e, 0xc2, 0x0a, 0xc9, 0x59, 0xae, 0x73, 0xf6, 0xdf, 0x06, 0x74, 0x47, 0x6c, 0x51, 0xb2, 0x82,
	0x14, 0x12, 0x9d, 0xc1, 0xb1, 0x4a, 0x47, 0xb4, 0x20, 0x3c, 0xa4, 0x8b, 0x28, 0x23, 0x7d, 0xe3,
	0xd4, 0x38, 0xef, 0xe2, 0xde, 0x3a, 0xec, 0xab, 0x28, 0x7a, 0x01, 0x9d, 0x7b, 0x5a, 0xa4, 0xfd,
	0xd6, 0xa9, 0x71, 0xde, 0x1b, 0x7e, 0xed, 0xec, 0xd8, 0xb0, 0xb3, 0xbe, 0xca, 0x79, 0x41, 0x8b,
	0x14, 0xd7, 0x24, 0x08, 0x41, 0xa7, 0x64, 0x2c, 0xef, 0xb7, 0xeb, 0xab, 0xea, 0xb5, 0xfd, 0x14,
	0x3a, 0x0a, 0x81, 0x8e, 0xc1, 0x9c, 0x4f, 0x82, 0x1b, 0x6f, 0xe4, 0x3f, 0xf7, 0xbd, 0x6b, 0x6b,
	0x0f, 0x01, 0x1c, 0x5c, 0x63, 0xff, 0xd6, 0xc3, 0x96, 0xa1, 0xd6, 0x81, 0x87, 0xd5, 0xba, 0xa5,
	0xd6, 0xa3, 0xb1, 0xef, 0x4d, 0x66, 0x56, 0xdb, 0xfe, 0xa7, 0x05, 0x1f, 0x07, 0x32, 0xe2, 0x72,
	0x46, 0x84, 0x0c, 0x88, 0x10, 0x94, 0x15, 0x98, 0x3c, 0x54, 0x44, 0x48, 0x34, 0x84, 0x47, 0x22,
	0x21, 0x45, 0xc4, 0x29, 0xab, 0x7b, 0x33, 0x87, 0x1f, 0x6d, 0x57, 0x1f, 0x34, 0x59, 0xbc, 0xc6,
	0xa1, 0x1f, 0xe0, 0x20, 0xe5, 0x74, 0x49, 0x78, 0xdd, 0xaf, 0x39, 0x1c, 0xbe, 0x7b, 0xbf, 0xb8,
	0x61, 0x40, 0x63, 0x38, 0xfc, 0x95, 0xf1, 0x7b, 0xc2, 0x45, 0xbf, 0x7d, 0xda, 0x7e, 0x4f, 0xb2,
	0x15, 0x05, 0xfa, 0x12, 0xac, 0xf8, 0x21, 0xe4, 0x44, 0x54, 0xb9, 0x14, 0xa1, 0x8c, 0xe2, 0x9c,
	0xf4, 0x3b, 0x4a, 0xc6, 0xef, 0xf7, 0x70, 0x2f, 0x7e, 0xc0, 0x3a, 0x31, 0x53, 0x71, 0x74, 0x01,
	0x8f, 0xb3, 0x44, 0xac, 0xc1, 0xaf, 0x68, 0x4e, 0xc2, 0x8a, 0xd3, 0xfe, 0x7e, 0x83, 0x3f, 0xc9,
	0x12, 0xd1, 0x1c, 0x78, 0x4e, 0x73, 0x32, 0xe7, 0xf4, 0xd9, 0x63, 0x40, 0x1a, 0x1e, 0xa6, 0x75,
	0x79, 0xb5, 0xed, 0x6c, 0x0e, 0x27, 0x5b, 0xc2, 0x2a, 0x00, 0x7a, 0x02, 0xa6, 0xee, 0x30, 0xcc,
	0x59, 0x26, 0x6a, 0xa1, 0x8e, 0x30, 0xe8, 0xd0, 0x98, 0x65, 0x02, 0x7d, 0x03, 0x47, 0xca, 0x98,
	0x21, 0xc9, 0xa3, 0x52, 0x90, 0xb4, 0x7e, 0x6d, 0x73, 0xf8, 0x89, 0xa3, 0xdd, 0xeb, 0xac, 0xdc,
	0xeb, 0x5c, 0x37, 0xee, 0xc6, 0xa6, 0x82, 0x7b, 0x1a, 0x6d, 0xff, 0xd6, 0x82, 0x7d, 0x6f, 0xa9,
	0x3c, 0xda, 0x87, 0x43, 0x51, 0xc5, 0xbf, 0x90, 0x44, 0x36, 0xde, 0x5c, 0x6d, 0xd1, 0x77, 0x5b,
	0xa6, 0xbc, 0xdc, 0x59, 0xd7, 0x9a, 0x77, 0xd3, 0x90, 0x0e, 0x74, 0xd4, 0xdd, 0x4d, 0x89, 0x83,
	0xff, 0x94, 0x38, 0x5b, 0x7d, 0x60, 0xb8, 0xc6, 0xa1, 0x53, 0x30, 0x53, 0x22, 0x12, 0x4e, 0x4b,
	0x55, 0xb8, 0x7e, 0x00, 0xbc, 0x19, 0xb2, 0x7f, 0x6e, 0xec, 0x8c, 0xa0, 0xe7, 0x4f, 0x66, 0x1e,
	0x9e, 0x5c, 0x8d, 0x43, 0x0f, 0xe3, 0x29, 0xb6, 0xf6, 0x50, 0x17, 0xf6, 0x5f, 0xce, 0xbd, 0xb9,
	0xa7, 0x0d, 0x7d, 0x35, 0x1a, 0x79, 0x37, 0x33, 0xab, 0x85, 0x3e, 0x80, 0xee, 0x0d, 0x9e, 0xde,
	0xfa, 0x81, 0x3f, 0x9d, 0x58, 0x6d, 0x74, 0x08, 0x6d, 0x3c, 0x9f, 0x58, 0x1d, 0xf4, 0x08, 0x3a,
	0xd7, 0xd3, 0x89, 0x67, 0xed, 0xab, 0x83, 0x9a, 0xe3, 0xc0, 0xfe, 0xcb, 0x80, 0x0f, 0x37, 0xde,
	0xe4, 0x47, 0x22, 0xa3, 0x34, 0x92, 0x11, 0x7a, 0x09, 0x47, 0x79, 0xa4, 0x24, 0x08, 0x89, 0x6a,
	0xb2, 0x71, 0xbc, 0xf3, 0x6e, 0xd2, 0x60, 0x53, 0x73, 0x68, 0xfd, 0xcf, 0xe0, 0xb8, 0x99, 0x76,
	0xe1, 0x92, 0x70, 0x75, 0x5b, 0x2d, 0x78, 0x17, 0xf7, 0x9a, 0xf0, 0xad, 0x8e, 0xa2, 0xa7, 0x60,
	0x26, 0x9c, 0x44, 0x92, 0x84, 0x3b, 0x8a, 0x09, 0x1a, 0xae, 0x02, 0xc3, 0x3f, 0x0c, 0x38, 0x09,
	0xd6, 0xa5, 0x04, 0x9a, 0x19, 0xfd, 0x6e, 0x80, 0xf5, 0xe6, 0x87, 0x8d, 0xbe, 0xdd, 0xb9, 0x9b,
	0xff, 0x99, 0x09, 0x83, 0x4f, 0x57, 0x45, 0x6d, 0xcc, 0x61, 0x67, 0xba, 0x9a, 0xc3, 0xf6, 0x57,
	0x7f, 0x5e, 0x9d, 0xbf, 0xcd, 0xf2, 0x6f, 0x53, 0xfc, 0xd9, 0xe5, 0x4f, 0x17, 0x19, 0x95, 0x77,
	0x55, 0xec, 0x24, 0x6c, 0xe1, 0xd6, 0x93, 0x77, 0x3d, 0x7e, 0x13, 0xc9, 0x73, 0x3d, 0xa6, 0xdd,
	0xad, 0x3f, 0x49, 0x7c, 0x50, 0x07, 0x2f, 0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xb2, 0x84, 0x9e,
	0x04, 0x61, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SchedulingServiceClient is the client API for SchedulingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SchedulingServiceClient interface {
	// Starts a test session, which encapsulates a test's scenario, required
	// components and output settings.
	//
	// The scenario contains the configuration for the test itself. The components
	// represent physical or virtual resources that must be reserved and monitored
	// when running a test. The output settings specify where the result will be
	// saved.
	//
	// This is a long-running operation that is managed by the Operations service
	// on the same server.  The long-running operation will assign a unique
	// identifier to the session.
	//
	// The unique identifier can be used to poll for the session's status and
	// result. Streaming is not supported at present.
	//
	// It can also be used to cancel while in progress, but cancellations operate
	// like killing a running process. If they occur while result are being
	// reported, the persistent storage may receive some but not all of them.
	//
	// For the specification of a google.longrunning.Operation message, see the
	// Long-running AIP at https://aip.dev/151.
	StartTestSession(ctx context.Context, in *StartTestSessionRequest, opts ...grpc.CallOption) (*longrunning.Operation, error)
}

type schedulingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSchedulingServiceClient(cc grpc.ClientConnInterface) SchedulingServiceClient {
	return &schedulingServiceClient{cc}
}

func (c *schedulingServiceClient) StartTestSession(ctx context.Context, in *StartTestSessionRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	out := new(longrunning.Operation)
	err := c.cc.Invoke(ctx, "/grpc.testing.benchmarking.scheduling.v1.SchedulingService/StartTestSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SchedulingServiceServer is the server API for SchedulingService service.
type SchedulingServiceServer interface {
	// Starts a test session, which encapsulates a test's scenario, required
	// components and output settings.
	//
	// The scenario contains the configuration for the test itself. The components
	// represent physical or virtual resources that must be reserved and monitored
	// when running a test. The output settings specify where the result will be
	// saved.
	//
	// This is a long-running operation that is managed by the Operations service
	// on the same server.  The long-running operation will assign a unique
	// identifier to the session.
	//
	// The unique identifier can be used to poll for the session's status and
	// result. Streaming is not supported at present.
	//
	// It can also be used to cancel while in progress, but cancellations operate
	// like killing a running process. If they occur while result are being
	// reported, the persistent storage may receive some but not all of them.
	//
	// For the specification of a google.longrunning.Operation message, see the
	// Long-running AIP at https://aip.dev/151.
	StartTestSession(context.Context, *StartTestSessionRequest) (*longrunning.Operation, error)
}

// UnimplementedSchedulingServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSchedulingServiceServer struct {
}

func (*UnimplementedSchedulingServiceServer) StartTestSession(ctx context.Context, req *StartTestSessionRequest) (*longrunning.Operation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTestSession not implemented")
}

func RegisterSchedulingServiceServer(s *grpc.Server, srv SchedulingServiceServer) {
	s.RegisterService(&_SchedulingService_serviceDesc, srv)
}

func _SchedulingService_StartTestSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartTestSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulingServiceServer).StartTestSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.testing.benchmarking.scheduling.v1.SchedulingService/StartTestSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulingServiceServer).StartTestSession(ctx, req.(*StartTestSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SchedulingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.testing.benchmarking.scheduling.v1.SchedulingService",
	HandlerType: (*SchedulingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartTestSession",
			Handler:    _SchedulingService_StartTestSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scheduling/v1/scheduling_service.proto",
}

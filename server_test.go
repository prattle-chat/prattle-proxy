package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	bufSize      = 1024 * 1024
	emptyHandler = func(_ context.Context, _ interface{}) (resp interface{}, err error) {
		return nil, nil
	}

	config = &Configuration{
		DomainName: "testing",
		Federations: map[string]*Federation{
			"none": {
				PSK:       "blahblahblah",
				messaging: dummyMessagingClient{},
				group:     dummyGroupsClient{},
				user:      dummyUserClient{},
			},
		},
		MaxKeys: 5,
	}

	lis       = bufconn.Listen(bufSize)
	bufDialer = func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
)

type dummyClientStream struct{}

func (dummyClientStream) Header() (metadata.MD, error) { return nil, nil }
func (dummyClientStream) Trailer() metadata.MD         { return nil }
func (dummyClientStream) CloseSend() error             { return nil }
func (dummyClientStream) Context() context.Context     { return context.Background() }
func (dummyClientStream) SendMsg(interface{}) error    { return nil }
func (dummyClientStream) RecvMsg(interface{}) error    { return nil }

type dummyPKC struct {
	run int
	grpc.ClientStream
}

func (d *dummyPKC) Recv() (*server.PublicKeyValue, error) {
	if d.run > 0 {
		return nil, io.EOF
	}

	d.run++
	return &server.PublicKeyValue{Value: "key-1"}, nil
}

type dummySC struct {
	run int
	grpc.ClientStream
}

func (d *dummySC) Recv() (*server.MessageWrapper, error) {
	if d.run > 0 {
		return nil, io.EOF
	}

	d.run++
	return &server.MessageWrapper{}, nil
}

type Messaging_SubscribeClient interface {
	Recv() (*server.MessageWrapper, error)
	grpc.ClientStream
}

type dummyMessagingClient struct{}

func (d dummyMessagingClient) Subscribe(context.Context, *emptypb.Empty, ...grpc.CallOption) (server.Messaging_SubscribeClient, error) {
	return &dummySC{0, dummyClientStream{}}, nil
}

func (d dummyMessagingClient) Send(context.Context, *server.MessageWrapper, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

type dummyGroupsClient struct{}

func (dummyGroupsClient) Create(context.Context, *server.Group, ...grpc.CallOption) (*server.Group, error) {
	return nil, nil
}

func (dummyGroupsClient) Join(context.Context, *server.JoinRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (dummyGroupsClient) Info(context.Context, *server.InfoRequest, ...grpc.CallOption) (*server.Group, error) {
	return new(server.Group), nil
}
func (dummyGroupsClient) Invite(context.Context, *server.InviteRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (dummyGroupsClient) PromoteUser(context.Context, *server.PromoteRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (dummyGroupsClient) DemoteUser(context.Context, *server.DemoteRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (dummyGroupsClient) Leave(context.Context, *server.LeaveRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

type dummyUserClient struct{}

func (d dummyUserClient) AddPublicKey(context.Context, *server.PublicKeyValue, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (d dummyUserClient) DelPublicKey(context.Context, *server.PublicKeyValue, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (d dummyUserClient) Tokens(context.Context, *emptypb.Empty, ...grpc.CallOption) (server.User_TokensClient, error) {
	return nil, nil
}
func (d dummyUserClient) DelToken(context.Context, *server.TokenValue, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (d dummyUserClient) PublicKey(context.Context, *server.PublicKeyRequest, ...grpc.CallOption) (server.User_PublicKeyClient, error) {
	return &dummyPKC{0, dummyClientStream{}}, nil
}

type key string

func (k key) Auth() (ctx context.Context) {
	ctx = context.Background()
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("bearer %s", k)})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return
}

func newTestServer(r Redis) (s Server) {
	s = Server{
		UnimplementedAuthenticationServer: server.UnimplementedAuthenticationServer{},
		UnimplementedGroupsServer:         server.UnimplementedGroupsServer{},
		UnimplementedMessagingServer:      server.UnimplementedMessagingServer{},
		UnimplementedUserServer:           server.UnimplementedUserServer{},
		redis:                             r,
		config:                            config,
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			s.FederatedEndpointStreamInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			s.FederatedEndpointUnaryInterceptor,
		)),
	)

	server.RegisterAuthenticationServer(grpcServer, s)
	server.RegisterGroupsServer(grpcServer, s)
	server.RegisterMessagingServer(grpcServer, s)
	server.RegisterUserServer(grpcServer, s)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return
}

func TestDefautMinter(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("unexpected error\n%#v", err)
		}
	}()

	_, err := minter("testing")
	if err != nil {
		t.Errorf("unexpected error\n%#v", err)
	}
}

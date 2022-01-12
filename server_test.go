package main

import (
	"context"
	"fmt"
	"net"

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
			},
		},
	}

	lis       = bufconn.Listen(bufSize)
	bufDialer = func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
)

type dummyMessagingClient struct{}

func (d dummyMessagingClient) Subscribe(context.Context, *emptypb.Empty, ...grpc.CallOption) (server.Messaging_SubscribeClient, error) {
	return nil, nil
}
func (d dummyMessagingClient) PublicKey(context.Context, *server.Auth, ...grpc.CallOption) (server.Messaging_PublicKeyClient, error) {
	return nil, nil
}
func (d dummyMessagingClient) Send(context.Context, *server.MessageWrapper, ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
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
		UnimplementedSelfServer:           server.UnimplementedSelfServer{},
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
	server.RegisterSelfServer(grpcServer, s)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return
}

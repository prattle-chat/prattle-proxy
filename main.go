package main

import (
	"crypto/tls"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/prattle-chat/prattle-proxy/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	redis, err := NewRedis(config.RedisAddr)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		panic(err)
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		panic(err)
	}

	s := Server{
		UnimplementedAuthenticationServer: server.UnimplementedAuthenticationServer{},
		UnimplementedGroupsServer:         server.UnimplementedGroupsServer{},
		UnimplementedMessagingServer:      server.UnimplementedMessagingServer{},
		UnimplementedUserServer:           server.UnimplementedUserServer{},
		redis:                             redis,
		config:                            config,
	}

	logger, err := zap.NewProduction(zap.AddCallerSkip(6))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
			s.FederatedEndpointStreamInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
			s.FederatedEndpointUnaryInterceptor,
		)),
	)

	reflection.Register(grpcServer)

	server.RegisterAuthenticationServer(grpcServer, s)
	server.RegisterGroupsServer(grpcServer, s)
	server.RegisterMessagingServer(grpcServer, s)
	server.RegisterUserServer(grpcServer, s)

	log.Printf("Starting server on %s", config.ListenAddr)

	panic(grpcServer.Serve(lis))
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("certs/server-cert.pem", "certs/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	return credentials.NewTLS(config), nil
}

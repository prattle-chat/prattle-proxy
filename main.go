package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	RedisAddr = os.Getenv("REDIS_ADDR")
)

func main() {
	redis, err := NewRedis(RedisAddr)
	if err != nil {
		panic(err)
	}

	// #nosec
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)

	reflection.Register(grpcServer)

	s := Server{
		UnimplementedAuthenticationServer: server.UnimplementedAuthenticationServer{},
		UnimplementedGroupsServer:         server.UnimplementedGroupsServer{},
		UnimplementedMessagingServer:      server.UnimplementedMessagingServer{},
		UnimplementedSelfServer:           server.UnimplementedSelfServer{},
		redis:                             redis,
	}

	server.RegisterAuthenticationServer(grpcServer, s)
	server.RegisterGroupsServer(grpcServer, s)
	server.RegisterMessagingServer(grpcServer, s)
	server.RegisterSelfServer(grpcServer, s)

	log.Print("Starting server")

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

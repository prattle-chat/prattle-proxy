package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	inaccessiblePeerError = status.Error(codes.Unavailable, "unable to connect to peer")
)

// Federation holds per federated-connection details,
// such as PSKs and domain names
type Federation struct {
	// ConnectionString is a gRPC connection address for
	// another server which talks prattle
	ConnectionString string `mapstructure:"connection_string"`

	// PSK is used as a bearer token when proxying
	// connections to and from this Federation instance
	PSK string `mapstructure:"psk"`

	// messaging holds a gRPC connection to this federated prattle
	// instance on the messaging namespace
	messaging server.MessagingClient

	// group holds a gRPC connection to this federated prattle
	// instance on the groups namespace
	group server.GroupsClient

	// user holds a gRPC connection to this federated prattle
	// instance on the user namespace
	user server.UserClient
}

func (f *Federation) connect() (err error) {
	// #nosec
	conn, err := grpc.Dial(f.ConnectionString, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		return inaccessiblePeerError
	}

	f.messaging = server.NewMessagingClient(conn)
	f.group = server.NewGroupsClient(conn)
	f.user = server.NewUserClient(conn)

	return
}

func (f Federation) auth(id string) (ctx context.Context) {
	ctx = context.Background()

	md := metadata.New(map[string]string{
		"authorization":  fmt.Sprintf("bearer %s", f.PSK),
		operatorIDHeader: id,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return
}

// Send proxies a message to a federated connection
func (f Federation) Send(id string, mw *server.MessageWrapper) (err error) {
	_, err = f.messaging.Send(f.auth(id), mw)

	return cleanError(err)
}

// PubicKey proxies a PublicKey request to a peered prattle
func (f Federation) PublicKey(in *server.PublicKeyRequest, pks server.User_PublicKeyServer) (err error) {
	kc, err := f.user.PublicKey(f.auth("skipped-open-request"), in)
	if err != nil {
		return cleanError(err)
	}

	for {
		var k *server.PublicKeyValue
		k, err = kc.Recv()
		if err != nil && err == io.EOF {
			err = nil

			break
		}

		if err != nil {
			break
		}

		err = pks.Send(k)
		if err != nil {
			break
		}
	}

	return cleanError(err)
}

func (f Federation) JoinGroup(id string, in *server.JoinRequest) (err error) {
	_, err = f.group.Join(f.auth(id), in)

	return cleanError(err)
}

func (f Federation) GroupInfo(id string, in *server.InfoRequest) (out *server.Group, err error) {
	return f.group.Info(f.auth(id), in)
}

func (f Federation) InviteToGroup(id string, in *server.InviteRequest) (err error) {
	_, err = f.group.Invite(f.auth(id), in)

	return cleanError(err)
}

func (f Federation) PromoteUser(id string, in *server.PromoteRequest) (err error) {
	_, err = f.group.PromoteUser(f.auth(id), in)

	return cleanError(err)
}

func (f Federation) DemoteUser(id string, in *server.DemoteRequest) (err error) {
	_, err = f.group.DemoteUser(f.auth(id), in)

	return cleanError(err)
}

func (f Federation) LeaveGroup(id string, in *server.LeaveRequest) (err error) {
	_, err = f.group.Leave(f.auth(id), in)

	return cleanError(err)
}

func cleanError(err error) error {
	if err == nil {
		return err
	}

	switch {
	case strings.HasSuffix(err.Error(), "connect: connection refused\""):
		return inaccessiblePeerError
	}

	return err
}

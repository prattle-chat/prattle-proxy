package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
}

func (f *Federation) connect() (err error) {
	conn, err := grpc.Dial(f.ConnectionString, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		return
	}

	f.messaging = server.NewMessagingClient(conn)
	f.group = server.NewGroupsClient(conn)

	return
}

func (f Federation) auth() (ctx context.Context) {
	ctx = context.Background()

	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("bearer %s", f.PSK)})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return
}

// Send proxies a message to a federated connection
func (f Federation) Send(mw *server.MessageWrapper) (err error) {
	_, err = f.messaging.Send(f.auth(), mw)

	return
}

// PubicKey proxies a PublicKey request to a peered prattle
func (f Federation) PublicKey(in *server.Auth, pks server.Messaging_PublicKeyServer) (err error) {
	kc, err := f.messaging.PublicKey(f.auth(), in)
	if err != nil {
		return
	}

	for {
		var k *server.PublicKeyValue
		k, err = kc.Recv()
		if err != nil && err == io.EOF {
			break
		}

		if err != nil {
			return
		}

		pks.Send(k)
	}

	return
}

func (f Federation) JoinGroup(in *server.GroupUser) (err error) {
	_, err = f.group.Join(f.auth(), in)

	return
}

func (f Federation) GroupInfo(in *server.Group) (err error) {
	_, err = f.group.Info(f.auth(), in)

	return
}

func (f Federation) InviteToGroup(in *server.GroupUser) (err error) {
	_, err = f.group.Invite(f.auth(), in)

	return
}

func (f Federation) PromoteUser(in *server.GroupUser) (err error) {
	_, err = f.group.PromoteUser(f.auth(), in)

	return
}

func (f Federation) DemoteUser(in *server.GroupUser) (err error) {
	_, err = f.group.DemoteUser(f.auth(), in)

	return
}

func (f Federation) LeaveGroup(in *server.GroupUser) (err error) {
	_, err = f.group.DemoteUser(f.auth(), in)

	return
}

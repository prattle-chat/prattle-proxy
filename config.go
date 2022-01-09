package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/prattle-chat/prattle-proxy/server"
	"github.com/spf13/viper"
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

	// conn holds a gRPC connection to this federated prattle
	// instance on the messaging namespace
	conn server.MessagingClient
}

func (f *Federation) connect() (err error) {
	conn, err := grpc.Dial(f.ConnectionString, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		return
	}

	f.conn = server.NewMessagingClient(conn)

	return
}

// Send proxies a message to a federated connection
func (f Federation) Send(mw *server.MessageWrapper) (err error) {
	ctx := context.Background()

	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("bearer %s", f.PSK)})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err = f.conn.Send(ctx, mw)

	return
}

// PubicKey proxies a PublicKey request to a peered prattle
func (f Federation) PublicKey(in *server.Auth, pks server.Messaging_PublicKeyServer) (err error) {
	ctx := context.Background()

	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("bearer %s", f.PSK)})
	ctx = metadata.NewOutgoingContext(ctx, md)

	kc, err := f.conn.PublicKey(ctx, in)
	if err != nil {
		return
	}

	for {
		k, err := kc.Recv()
		if err != nil && err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		pks.Send(k)
	}

	return
}

// Configuration holds configurables for the prattle-proxy
type Configuration struct {
	// DomainName contains the domain name associated with this
	// prattle
	DomainName string `mapstructure:"domain_name"`

	// ListenAddr is a host/port combination on which the prattle
	// proxy should listen
	ListenAddr string `mapstructure:"listen_addr"`

	// RedisAddr points to a redis instance for this instance of
	// prattle-proxy to use
	RedisAddr string `mapstructure:"redis_addr"`

	// MaxTokens is the maximum amount of tokens an account can
	// have.
	//
	// This essentially means the number of active connections a
	// user can have.
	MaxTokens int `mapstructure:"max_tokens"`

	// MaxKeys is the maximum amount of keys a user can have.
	//
	// It essentially puts a limit on the number of connections
	// (active or inactive) a user can have
	MaxKeys int `mapstructure:"max_keys"`

	// RevalidateFrequency governs how often the proxy checks
	// whether a token is still valid
	RevalidateFrequency time.Duration `mapstructure:"revalidate_frequency"`

	// Federations is a map of domain names to Federation structs
	// in order to proxy chats and keys and stuff
	Federations map[string]*Federation `mapstructure:"federations"`
}

func LoadConfig() (c *Configuration, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	viper.AddConfigPath("/etc/prattle")
	viper.AddConfigPath(filepath.Join(home, ".config/prattle"))
	viper.AddConfigPath(".")

	viper.SetDefault("ListenAddr", "0.0.0.0:8080")
	viper.SetDefault("RedisAddr", "localhot:6379")
	viper.SetDefault("MaxTokens", 5)
	viper.SetDefault("MaxKeys", 10)
	viper.SetDefault("RevalidateFrequence", time.Second)

	viper.SetConfigName("proxy")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	c = new(Configuration)
	err = viper.Unmarshal(c)
	if err != nil {
		return
	}

	for _, f := range c.Federations {
		err = f.connect()
		if err != nil {
			return
		}
	}

	// Ensure each federation has a unique PSK
	psk := make(map[string]int)
	for _, f := range c.Federations {
		psk[f.PSK] = 1
	}

	if len(psk) != len(c.Federations) {
		err = fmt.Errorf("each federation must have a unique key")
	}

	return
}

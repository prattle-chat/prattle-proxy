package main

import (
	"context"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	maxKeys = 10
)

func (s Server) AddPublicKey(ctx context.Context, pkv *server.PublicKeyValue) (_ *emptypb.Empty, err error) {
	u, err := s.userFromContext(ctx)
	if err != nil {
		return
	}

	u.PublicKeys = append(u.PublicKeys, pkv.Value)
	if len(u.PublicKeys) >= maxKeys {
		u.PublicKeys = u.PublicKeys[len(u.PublicKeys)-maxKeys : len(u.PublicKeys)]
	}

	err = s.redis.saveUser(u)
	if err != nil {
		err = generalError

		return
	}

	return new(emptypb.Empty), err
}

func (s Server) PublicKey(in *server.Auth, pks server.Proxy_PublicKeyServer) (err error) {
	_, err = s.userFromContext(pks.Context())
	if err != nil {
		return
	}

	// todo: this is the point at which we need to test the hostname of
	// in.UserId - if the hostname is our host then we can lookup in redis.
	//
	// else: we need to query the remote host associated with this hostname
	// and get public keys from there

	keys, err := s.redis.GetPublicKeys(in.UserId)
	if err != nil {
		err = generalError

		return
	}

	for _, key := range keys {
		pks.Send(&server.PublicKeyValue{
			Value: key,
		})
	}

	return
}
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
	u, err := s.userFromContext(pks.Context())
	if err != nil {
		return
	}

	// todo: this is the point at which we need to test the hostname of
	// in.UserId - if the hostname is our host then we can lookup in redis.
	//
	// else: we need to query the remote host associated with this hostname
	// and get public keys from there

	lookupIds := make([]string, 0)
	if groupRegexp.Match([]byte(in.UserId)) {
		g, err := s.redis.Group(in.UserId)
		if err != nil || !HasPermission(u, g, groupRead) {
			return badGroupError
		}

		for _, u := range g.Members {
			lookupIds = append(lookupIds, u)
		}
	} else {
		lookupIds = append(lookupIds, in.UserId)
	}

	for _, id := range lookupIds {
		keys, err := s.redis.GetPublicKeys(id)
		if err != nil {
			return generalError
		}

		for _, key := range keys {
			pks.Send(&server.PublicKeyValue{
				Value: key,
			})
		}
	}

	return
}

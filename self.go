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

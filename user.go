package main

import (
	"context"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Server) AddPublicKey(ctx context.Context, pkv *server.PublicKeyValue) (_ *emptypb.Empty, err error) {
	u, err := s.userFromContext(ctx)
	if err != nil {
		return
	}

	u.PublicKeys = append(u.PublicKeys, pkv.Value)
	if len(u.PublicKeys) >= s.config.MaxKeys {
		u.PublicKeys = u.PublicKeys[len(u.PublicKeys)-s.config.MaxKeys : len(u.PublicKeys)]
	}

	err = s.redis.saveUser(u)
	if err != nil {
		err = generalError

		return
	}

	return new(emptypb.Empty), err
}

// PublicKey returns the public keys associated with a user.
//
// Trying to load keys associated with a group will return an error; clients
// should maintain a store of individual user keys
//
// This function accepts either a local user, or a user for a peered
// server.
func (s Server) PublicKey(in *server.PublicKeyRequest, pks server.User_PublicKeyServer) (err error) {
	if groupRegexp.Match([]byte(in.Owner.Id)) {
		return inputError
	}

	d, err := domain(in.Owner.Id)
	if err != nil {
		return inputError
	}

	if d == s.config.DomainName {
		var u User

		u, err = s.redis.loadUser(in.Owner.Id)
		if err != nil {
			return
		}

		for _, k := range u.PublicKeys {
			err = pks.Send(&server.PublicKeyValue{
				Value: k,
			})

			if err != nil {
				err = generalError

				return
			}
		}

		return
	}

	peer, ok := s.config.Federations[d]
	if !ok {
		return notPeeredError
	}

	return peer.PublicKey(in, pks)
}

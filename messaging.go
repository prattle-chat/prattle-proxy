package main

import (
	"context"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Server) Subscribe(_ *emptypb.Empty, ss server.Messaging_SubscribeServer) (err error) {
	u, err := s.userFromContext(ss.Context())
	if err != nil {
		return
	}

	for m := range s.redis.Messages(u.Id) {
		ss.Send(&server.MessageWrapper{
			Encoded: []byte(m.Payload),
		})
	}
	return nil
}

func (s Server) Send(ctx context.Context, in *server.MessageWrapper) (out *emptypb.Empty, err error) {
	u, err := s.userFromContext(ctx)
	if err != nil {
		return
	}

	out = new(emptypb.Empty)

	if groupRegexp.Match([]byte(in.Recipient)) {
		g, err := s.redis.Group(in.Recipient)
		if err != nil {
			err = badGroupError

			return out, err
		}

		if HasPermission(u, g, groupPost) {
			for _, m := range g.Members {
				err = s.redis.WriteMessage(m, in.Encoded)
				if err != nil {
					return out, err
				}
			}
		}

		return out, err
	}

	// add to redis pub key of 'wrapper.Recipient'
	err = s.redis.WriteMessage(in.Recipient, in.Encoded)

	return
}

func (s Server) PublicKey(in *server.Auth, pks server.Messaging_PublicKeyServer) (err error) {
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

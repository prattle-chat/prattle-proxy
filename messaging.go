package main

import (
	"context"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Server) Subscribe(_ *emptypb.Empty, ss server.Proxy_SubscribeServer) (err error) {
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

func contains(sl []string, s string) bool {
	for _, elem := range sl {
		if elem == s {
			return true
		}
	}

	return false
}

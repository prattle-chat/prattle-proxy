package main

import (
	"context"
	"log"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Server) Subscribe(_ *emptypb.Empty, ss server.Proxy_SubscribeServer) (err error) {
	u, err := s.userFromContext(ss.Context())
	if err != nil {
		return
	}

	for m := range s.redis.Messages(u.Id) {
		log.Printf("%#v", m)

		ss.Send(&server.MessageWrapper{
			Encoded: []byte(m.Payload),
		})
	}
	return nil
}

func (s Server) Send(ctx context.Context, in *server.MessageWrapper) (out *emptypb.Empty, err error) {
	_, err = s.userFromContext(ctx)
	if err != nil {
		return
	}

	// add to redis pub key of 'wrapper.Recipient'
	out = new(emptypb.Empty)
	err = s.redis.WriteMessage(in.Recipient, in.Encoded)

	return
}

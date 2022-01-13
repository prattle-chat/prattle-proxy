package main

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Server) Subscribe(_ *emptypb.Empty, ss server.Messaging_SubscribeServer) (err error) {
	m := ss.Context().Value(MetadataKey{}).(*Metadata)

	var (
		buf *bytes.Buffer
		out *server.MessageWrapper
	)

	for msg := range s.redis.Messages(m.Sender.Id) {
		buf = bytes.NewBuffer(msg)
		dec := gob.NewDecoder(buf)

		err = dec.Decode(&out)
		if err != nil {
			return
		}

		err = ss.Send(out)
		if err != nil {
			return generalError
		}
	}
	return nil
}

// Send accepts a message for a user and either proxies it to the correct server,
// or passes it to a user, depending on whether the user is peered or local (respectively).
//
// This endpoint will error if the user is a group; clients should encrypt and send messages
// to users directly, setting the Recipient field of the encoded Server.Message type to the
// group name.
func (s Server) Send(ctx context.Context, in *server.MessageWrapper) (out *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)

	if groupRegexp.Match([]byte(m.Recipient.Id)) {
		err = inputError

		return
	}

	out = new(emptypb.Empty)
	if !m.Recipient.IsLocal {
		in.Sender = m.Sender.Id
		err = m.Recipient.PeerConnection.Send(in)

		return
	}

	// if in.For is set, does Sender have permission to post
	// to this group?
	if in.For != "" {
		if !groupRegexp.Match([]byte(in.For)) {
			// in.For can only be used for groups
			err = inputError

			return
		}

		if !HasPermission(User{Id: m.Recipient.Id}, Group{Id: in.For}, groupPost) {
			err = badGroupError

			return
		}
	}

	msg := mwToBytes(in)
	err = s.redis.WriteMessage(m.Recipient.Id, msg)

	return
}

// PublicKey returns the public keys associated with a user.
//
// Trying to load keys associated with a group will return an error; clients
// should maintain a store of individual user keys
//
// This function accepts either a local user, or a user for a peered
// server.
func (s Server) PublicKey(in *server.Auth, pks server.Messaging_PublicKeyServer) (err error) {
	if groupRegexp.Match([]byte(in.UserId)) {
		return inputError
	}

	d, err := domain(in.UserId)
	if err != nil {
		return inputError
	}

	if d == s.config.DomainName {
		var u User

		u, err = s.redis.loadUser(in.UserId)
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

func mwToBytes(mw *server.MessageWrapper) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// #nosec
	enc.Encode(mw)

	return buf.Bytes()
}

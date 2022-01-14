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

	for msg := range s.redis.Messages(m.Operator.Id) {
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

// Send accepts a message for a user and either relays it to the correct server,
// or passes it to a user, depending on whether the user is peered or local (respectively).
//
// This endpoint will error if the user is a group; clients should encrypt and send messages
// to users directly, setting the Recipient field of the encoded Server.Message type to the
// group name.
func (s Server) Send(ctx context.Context, in *server.MessageWrapper) (out *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)
	recipient := m.Operand

	// Sending directly to a group wont work, since there's nothing
	// subscribed to that channel
	if groupRegexp.Match([]byte(recipient.Id)) {
		err = inputError

		return
	}

	out = new(emptypb.Empty)

	// If the recipient is remote, send to that proxy
	if !recipient.IsLocal {
		err = recipient.PeerConnection.Send(m.Operator.Id, in)

		return
	}

	// If we're sending on behalf of a group, have we permission to do so?
	if in.Recipient.GroupId != "" {
		if !groupRegexp.Match([]byte(in.Recipient.GroupId)) {
			// check that it *is* a group
			err = inputError

			return
		}

		// Load group (including from remote)
		var g Group
		g, err = s.loadGroup(m.Operator.Id, in.Recipient.GroupId)
		if err != nil {
			return
		}

		if !HasPermission(User{Id: m.Operator.Id}, g, groupPost) {
			err = badGroupError

			return
		}
	}

	msg := mwToBytes(in)
	err = s.redis.WriteMessage(recipient.Id, msg)

	return
}

func mwToBytes(mw *server.MessageWrapper) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// #nosec
	enc.Encode(mw)

	return buf.Bytes()
}

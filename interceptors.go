package main

import (
	"context"
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/grpc"
)

/**
  This interceptor should:

  1. Bomb out if missing auth
  2. Set a bool if received via federation
  3. Set the sender and receiver domains
  4. Set a reference to the *Federation that this requrest needs, should that occur
  5. Other stuff /shrug
  **/

type FederationWithDomain struct {
	*Federation

	Domain string
}

type Actor struct {
	Id             string
	IsLocal        bool
	PeerConnection *FederationWithDomain
}

type Metadata struct {
	Sender    Actor
	Recipient Actor
}

type wrappedStream struct {
	grpc.ServerStream

	ctx context.Context
}

func (w wrappedStream) Context() context.Context {
	return w.ctx
}

// FederatedEndpointUnaryInterceptor will check (on certain requets paths) for valid authentication,
// adding user information to contexts, along with whether the request is from a federated
// peer.
//
// It will also set metadata around source/ destination domains, such as what it is, and whether
// it's the same domain as this instance.
func (s Server) FederatedEndpointUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	m := new(Metadata)

	switch info.FullMethod {
	// These are the only endpoints which can be federated; for anything else
	// a user *must* connect to the prattle instance associated with their user
	case "/messaging.Messaging/Send",
		"/group.Groups/Create",
		"/group.Groups/DemoteUser",
		"/group.Groups/Info",
		"/group.Groups/Invite",
		"/group.Groups/Join",
		"/group.Groups/Leave",
		"/group.Groups/PromoteUser":

		_, err = s.auth(ctx, m)
		if err != nil {
			panic(err)
		}

		err = s.validateSender(req, m)
		if err != nil {
			panic(err)
		}

		// We can't test the recipient on group creation
		// because there isn't one yet- the group name is
		// minted later on.
		//
		// Thus: skip
		if info.FullMethod == "/group.Groups/Create" {
			break
		}

		err = s.validateRecipient(req, m)
		if err != nil {
			panic(err)
		}

	case "/self.Self/AddPublicKey",
		"/self.Self/DelPublicKey",
		"/self.Self/DelToken",
		"/self.Self/Tokens":

		// There's no federation on these requests; ensure that the
		// calling user exists on our end and then done
		var token string

		token, err = s.auth(ctx, m)
		if err != nil {
			return
		}

		_, err = s.idFromToken(token)
		if err != nil {
			err = inputError

			return
		}
	}

	ctx = context.WithValue(ctx, MetadataKey{}, m)

	return handler(ctx, req)
}

// FederatedEndpointStreamInterceptor validates incoming streaming requests against
// either a federated peer's PSK, or a user's token.
//
// We don't do much more than this; these endpoints are read-only and so don't accept
// any sender information like our unary endpoints do
func (s Server) FederatedEndpointStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	// Some clients, like grpcurl, use reflection to reflect json into protobuf
	// types. Let these go right through
	if strings.Contains(info.FullMethod, "/grpc.reflection.v1alpha.") {
		return handler(srv, ss)
	}

	m := new(Metadata)
	ctx := ss.Context()

	_, err = s.auth(ctx, m)
	if err != nil {
		return
	}

	switch info.FullMethod {
	case "/messaging.Messaging/PublicKey":
		// If auth is happy at this point, then I'm happy /shrug

	case "/messaging.Messaging/Subscribe":
		// Only accept a subscribe from hosts on this server
		// - I have no idea what would happen if a federated peer
		// tried to subscribe to messages on another prattle server.
		// I *do* know, however, that it would be unexpected
		if !m.Sender.IsLocal {
			err = inputError

			return
		}
	}

	ctx = context.WithValue(ctx, MetadataKey{}, m)

	return handler(srv, wrappedStream{ss, ctx})
}

// auth checks whether a connection has either a valid peering PSK or
// user token passed through as a bearer token
func (s Server) auth(ctx context.Context, m *Metadata) (token string, err error) {
	token, err = grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil || token == "" {
		err = authError
	}

	// test whether incoming is a valid token; whether fed or user
	id, err := s.idFromToken(token)
	if err != nil {
		err = generalError
	}
	if id != "" {
		m.Sender = Actor{
			Id:      id,
			IsLocal: true,
		}

		return
	}

	err = nil

	d, ok := s.isFederated(token)
	if !ok {
		err = notPeeredError

		return
	}

	m.Sender = Actor{
		IsLocal: false,
		PeerConnection: &FederationWithDomain{
			s.config.Federations[d],
			d,
		},
	}

	return
}

func (s Server) validateSender(req interface{}, m *Metadata) (err error) {
	var sender string
	if m.Sender.IsLocal {
		sender = m.Sender.Id
	} else {
		switch v := req.(type) {
		case *server.MessageWrapper:
			sender = v.Sender

		case *server.GroupUser:
			sender = v.UserId

		// Anything else is unfederated, so trust whatever we get from
		// reading tokens
		default:
			sender = m.Sender.Id
		}
	}

	// if local, ensure the relevant sender field matches
	// the one taken from the token.
	//
	// otherwise, ensure the source domain matches the same
	// domain in the request (and trust that the peer has already
	// verified correct ownership
	if m.Sender.IsLocal {
		if sender != m.Sender.Id {
			err = inputError
		}

		return
	}

	d, err := domain(sender)
	if err != nil {
		err = inputError

		return
	}

	m.Sender.Id = sender
	if m.Sender.PeerConnection.Domain != d {
		err = inputError
	}

	return
}

func (s Server) validateRecipient(req interface{}, m *Metadata) (err error) {
	var recipient string
	switch v := req.(type) {
	case *server.MessageWrapper:
		recipient = v.Recipient

	case *server.GroupUser:
		recipient = v.GroupId

	case *server.Group:
		recipient = v.Id

	default:
		err = inputError

		return
	}

	m.Recipient = Actor{
		Id: recipient,
	}

	// if recipient is on our host, fine
	d, err := domain(recipient)
	if err != nil {
		err = inputError

		return
	}

	if d == s.config.DomainName {
		m.Recipient.IsLocal = true

		return
	}

	// Federated connection
	pc, ok := s.config.Federations[d]
	if !ok {
		err = notPeeredError

		return
	}

	m.Recipient.PeerConnection = &FederationWithDomain{
		pc,
		d,
	}
	m.Recipient.IsLocal = false

	return
}

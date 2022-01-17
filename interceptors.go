package main

import (
	"context"
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/grpc"
)

type groupIdGetter interface {
	GetGroupId() string
}

var (
	operatorIDHeader = "operator_id"
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
	// Operator is the Actor performing an operation. This is derived using one of the following:
	//
	//    1. On calls where the bearer token matches a valid local user, the valid user becomes
	//       the operator
	//    2. On calls where the bearer token matches a valid peer, the valid user is taken from
	//       the incoming context
	Operator Actor

	// Operand is the entity being Operated on, such as a user (when sending a message, say or
	// retrieving a public key), or a group (when performing group operations).
	Operand Actor
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
			return
		}

		err = s.validateOperator(req, m)
		if err != nil {
			return
		}

		// We can't test the recipient on group creation
		// because there isn't one yet- the group name is
		// minted later on.
		//
		// Thus: skip
		if info.FullMethod == "/group.Groups/Create" {
			break
		}

		err = s.validateOperand(req, m)
		if err != nil {
			return
		}

	case "/user.User/AddPublicKey",
		"/user.User/DelPublicKey",
		"/user.User/DelToken",
		"/user.User/Tokens":

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
	case "/user.User/PublicKey":
		// If auth is happy at this point, then I'm happy /shrug

	case "/messaging.Messaging/Subscribe":
		// Only accept a subscribe from hosts on this server
		// - I have no idea what would happen if a federated peer
		// tried to subscribe to messages on another prattle server.
		// I *do* know, however, that it would be unexpected
		if !m.Operator.IsLocal {
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

		return
	}

	// test whether incoming is a valid token; whether fed or user
	id, err := s.idFromToken(token)
	if err != nil {
		err = generalError

		return
	}
	if id != "" {
		var u User
		u, err = s.redis.loadUser(id)
		if err != nil || u.Id == "" {
			// If we get here then we have an out-of-date token
			// that doesn't point to a real user.
			//
			// In that case, delete the token so that the next call
			// errors correctly
			err = s.redis.DeleteToken(token)
			if err != nil {
				err = generalError

				return
			}

			err = authError

			return
		}

		if !u.Finalised {
			err = needFinaliseError

			return
		}

		m.Operator = Actor{
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

	// get operator_id from incoming context
	id = metautils.ExtractIncoming(ctx).Get(operatorIDHeader)
	if id == "" {
		err = inputError

		return
	}

	m.Operator = Actor{
		Id:      id,
		IsLocal: false,
		PeerConnection: &FederationWithDomain{
			s.config.Federations[d],
			d,
		},
	}

	return
}

func (s Server) validateOperator(req interface{}, m *Metadata) (err error) {
	// All messages must contain a sender
	if m.Operator.Id == "" {
		return inputError
	}

	// Ensure that, for messages, the sender matches the derived Operator ID
	//
	// This is the first step in validating a sender is entitled to send a message;
	// Here we compare the OperatorID with the sender field of a MessageWrapper.
	//
	// A client, then, compares the sender field of the MessageWrapper with the encoded
	// sender in the encrypted, embeded Message
	switch v := req.(type) {
	case *server.MessageWrapper:
		if v.Sender == nil {
			return inputError
		}

		if m.Operator.Id != v.Sender.Id {
			return mismatchedSenderError
		}
	}

	// Trust local operators; we've already validated them by
	// loading the token from redis
	if m.Operator.IsLocal {
		return
	}

	// If from a peer, ensure the sender domain is at least correct
	// and trust that the sender has at least validated that the sender
	// and tokens are correct
	d, err := domain(m.Operator.Id)
	if err != nil {
		err = inputError

		return
	}

	if m.Operator.PeerConnection.Domain != d {
		return mismatchedDomainError
	}

	return
}

func (s Server) validateOperand(req interface{}, m *Metadata) (err error) {
	var recipient string
	switch v := req.(type) {
	case *server.MessageWrapper:
		if v.Recipient == nil {
			return inputError
		}

		recipient = v.Recipient.Id

	case groupIdGetter:
		recipient = v.GetGroupId()

	default:
		return inputError
	}

	if recipient == "" {
		return inputError
	}

	m.Operand = Actor{
		Id: recipient,
	}

	// if recipient is on our host, fine
	d, err := domain(recipient)
	if err != nil {
		err = inputError

		return
	}

	if d == s.config.DomainName {
		m.Operand.IsLocal = true

		return
	}

	// Federated connection
	pc, ok := s.config.Federations[d]
	if !ok {
		err = notPeeredError

		return
	}

	m.Operand.PeerConnection = &FederationWithDomain{
		pc,
		d,
	}
	m.Operand.IsLocal = false

	return
}

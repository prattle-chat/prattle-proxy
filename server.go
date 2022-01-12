package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/gofrs/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/pquerna/otp/totp"
	"github.com/prattle-chat/prattle-proxy/server"
	"github.com/sethvargo/go-diceware/diceware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	passwordError     = status.Error(codes.Unauthenticated, "password must be between 16 and 64 chars")
	authError         = status.Error(codes.Unauthenticated, "missing/ unreadble bearer token")
	needFinaliseError = status.Error(codes.Unauthenticated, "must finalise signup")
	badTotpError      = status.Error(codes.Unauthenticated, "incorrect totp token")
	badPasswordError  = status.Error(codes.Unauthenticated, "incorrect password or userid")

	badGroupError  = status.Error(codes.NotFound, "group could not be found")
	notPeeredError = status.Error(codes.NotFound, "recipient is on a non-peered domain")

	generalError = status.Error(codes.Unavailable, "an internal systems error occurred")
	inputError   = status.Error(codes.Unavailable, "missing/ poorly formed input")

	mismatchedSenderError = status.Error(codes.InvalidArgument, "mismatch between sender field and owner of token")
	mismatchedDomainError = status.Error(codes.InvalidArgument, "mismatch between sender domain and the domain this federated peer belongs to")

	minter MinterFunc = func(d string) (string, error) {
		words := diceware.MustGenerate(2)

		suffix := make([]byte, 4)
		_, err := rand.Read(suffix)
		if err != nil {
			return "", generalError
		}

		return fmt.Sprintf("%s-%s-%s@%s",
			words[0],
			words[1],
			hex.EncodeToString(suffix),
			d,
		), err

	}
)

// MinterFunc accepts a domain name, and returns a new
// ID (and optional error) which can be used in both user
// and group ID generation
type MinterFunc func(string) (string, error)

type MetadataKey struct{}

type Server struct {
	server.UnimplementedAuthenticationServer
	server.UnimplementedGroupsServer
	server.UnimplementedMessagingServer
	server.UnimplementedSelfServer

	redis  Redis
	config *Configuration
}

func (s Server) mintID() (id string, err error) {
	// Try a maximum of ten times to mint an unknown ID
	for i := 0; i < 10; i++ {
		id, err = minter(s.config.DomainName)
		if err != nil {
			return
		}

		// test whether id is in use
		if !s.redis.IDExists(id) {
			return
		}
	}

	// If we get this far then we couldn't get a fresh ID in
	// ten tries, and so we have a problem somewhere
	return "", generalError
}

func (s Server) mintToken() string {
	return fmt.Sprintf("prattle-%s%s%s",
		hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()),
		hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()),
		hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()),
	)
}

func (s Server) validateTOTP(user string, token string) (valid bool, err error) {
	seed, err := s.redis.GetTOTPSeed(user)
	if err != nil {
		return
	}

	valid = totp.Validate(token, seed)

	return
}

func (s Server) validatePassword(user, password string) (valid bool, err error) {
	hash, err := s.redis.GetPasswordHash(user)
	if err != nil {
		return
	}

	return argon2id.ComparePasswordAndHash(password, hash)
}

func (s Server) userFromContext(ctx context.Context) (u User, err error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		err = authError

		return
	}

	return s.redis.UserByToken(token)
}

func (s Server) isFederated(tok string) (string, bool) {
	for d, f := range s.config.Federations {
		if f.PSK == tok {
			return d, true
		}
	}

	return "", false
}

func (s Server) idFromToken(tok string) (string, error) {
	return s.redis.UserIdByToken(tok)
}

func contains(sl []string, s string) bool {
	for _, elem := range sl {
		if elem == s {
			return true
		}
	}

	return false
}

func domain(s string) (d string, err error) {
	ss := strings.Split(s, "@")
	if len(ss) != 2 {
		err = inputError

		return
	}

	d = ss[1]

	return
}

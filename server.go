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
	"github.com/microcosm-cc/bluemonday"
	"github.com/pquerna/otp/totp"
	"github.com/prattle-chat/prattle-proxy/server"
	"github.com/sethvargo/go-diceware/diceware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	sanitiser = bluemonday.StrictPolicy()

	passwordError     = status.Error(codes.Unauthenticated, "password must be between 16 and 64 chars")
	authError         = status.Error(codes.Unauthenticated, "missing/ unreadble bearer token")
	needFinaliseError = status.Error(codes.Unauthenticated, "must finalise signup")
	badTotpError      = status.Error(codes.Unauthenticated, "incorrect totp token")
	badPasswordError  = status.Error(codes.Unauthenticated, "incorrect password or userid")

	badGroupError  = status.Error(codes.NotFound, "group could not be found")
	notMemberError = status.Error(codes.NotFound, "user is not a member of this group")
	notPeeredError = status.Error(codes.NotFound, "recipient is on a non-peered domain")

	generalError = status.Error(codes.Unavailable, "an general error occurred")
	inputError   = status.Error(codes.Unavailable, "missing/ poorly formed input")
)

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
	words := diceware.MustGenerate(2)

	suffix := make([]byte, 4)
	_, err = rand.Read(suffix)
	if err != nil {
		err = generalError

		return
	}

	id = fmt.Sprintf("%s-%s-%s@%s",
		words[0],
		words[1],
		hex.EncodeToString(suffix),
		s.config.DomainName,
	)

	// test whether id is in use
	if s.redis.IDExists(id) {
		return s.mintID()
	}

	return
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

func (s Server) isThisDomain(id string) (ok bool, err error) {
	d, err := domain(id)
	if err != nil {
		return
	}

	ok = s.config.DomainName == d

	return
}

func (s Server) idFromToken(tok string) (string, error) {
	return s.redis.UserIdByToken(tok)
}

// userConnectors accepts a slice of IDs and returns a map
// of those IDs to a Federation client that can be used to
// do stuff
func (s Server) userConnectors(ids []string) (m map[string]*Federation, err error) {
	var d string

	m = make(map[string]*Federation)

	for _, id := range ids {
		d, err = domain(id)
		if err != nil {
			return
		}

		if d != s.config.DomainName {
			fc, ok := s.config.Federations[d]
			if !ok {
				err = notPeeredError

				return
			}

			m[id] = fc
		} else {
			m[id] = nil
		}
	}

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

func domain(s string) (d string, err error) {
	ss := strings.Split(s, "@")
	if len(ss) != 2 {
		err = inputError

		return
	}

	d = ss[1]

	return
}

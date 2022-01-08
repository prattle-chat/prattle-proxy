package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

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

const (
	proxyName = "prattle"
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

	generalError = status.Error(codes.Unavailable, "an general error occurred")
	inputError   = status.Error(codes.Unavailable, "missing/ poorly formed input")
)

type Server struct {
	server.UnimplementedAuthenticationServer
	server.UnimplementedGroupsServer
	server.UnimplementedMessagingServer
	server.UnimplementedSelfServer

	redis Redis
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
		proxyName,
	)

	// test whether id is in use
	if s.redis.IDExists(id) {
		return s.mintID()
	}

	return
}

func (s Server) mintToken() string {
	return fmt.Sprintf("%s-%s%s%s",
		proxyName,
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

func contains(sl []string, s string) bool {
	for _, elem := range sl {
		if elem == s {
			return true
		}
	}

	return false
}

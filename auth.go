package main

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/pquerna/otp/totp"
	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	passwordValidation = "required,min=16,max=64"
)

func (s Server) Signup(ctx context.Context, in *server.SignupRequest) (out *server.SignupResponse, err error) {
	validate := validator.New()

	if validate.Var(in.Password, passwordValidation) != nil {
		err = passwordError

		return
	}

	out = new(server.SignupResponse)

	out.UserId, err = s.mintID()
	if err != nil {
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "prattle",
		AccountName: out.UserId,
	})

	out.TotpSeed = key.Secret()

	hash, err := createHash(in.Password)
	if err != nil {
		err = passwordError

		return
	}

	err = s.redis.AddUser(out.UserId, out.TotpSeed, hash)
	if err != nil {
		err = generalError

		return
	}

	return
}

func (s Server) Finalise(ctx context.Context, in *server.Auth) (_ *emptypb.Empty, err error) {
	err = s.validateAuth(in)
	if err != nil {
		return
	}

	err = s.redis.MarkFinalised(in.UserId)

	return &emptypb.Empty{}, nil
}

func (s Server) Token(ctx context.Context, in *server.Auth) (out *server.TokenValue, err error) {
	err = s.validateAuth(in)
	if err != nil {
		return
	}

	out = &server.TokenValue{
		Value: s.mintToken(),
	}

	err = s.redis.AddToken(in.UserId, out.Value)

	return
}

func (s Server) validateAuth(in *server.Auth) (err error) {
	if in == nil || in.UserId == "" || in.Password == "" || in.Totp == "" {
		err = inputError

		return
	}

	if !s.redis.IDExists(in.UserId) {
		err = badPasswordError

		return
	}

	valid, err := s.validatePassword(in.UserId, in.Password)
	if err != nil || !valid {
		err = badPasswordError

		return
	}

	valid, err = s.validateTOTP(in.UserId, in.Totp)
	if err != nil || !valid {
		err = badTotpError

		return
	}

	return nil
}

func createHash(p string) (s string, err error) {
	return argon2id.CreateHash(p, argon2id.DefaultParams)
}

package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/prattle-chat/prattle-proxy/server"
	"github.com/rafaeljusto/redigomock/v3"
	"google.golang.org/grpc"
)

func dummyMinter(s string) (string, error) {
	return fmt.Sprintf("some-user@%s", s), nil
}

func newTestAuthClient() (c server.AuthenticationClient) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return server.NewAuthenticationClient(conn)
}

func TestServer_Signup(t *testing.T) {
	oldMinter := minter
	defer func() {
		minter = oldMinter
	}()

	minter = dummyMinter

	longPW := ""
	for i := 0; i < 128; i++ {
		longPW += "a"
	}

	goodPW := ""
	for i := 0; i < 16; i++ {
		goodPW += "a"
	}

	for _, test := range []struct {
		name        string
		password    string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"Missing password errors", "", noRedisCall, true},
		{"Too short password errors", "password", noRedisCall, true},
		{"Paradoxically, too long a password errors", longPW, noRedisCall, true},
		{"Inability to generate new ID fails in time", goodPW, idFound, true},
		{"Error on creating user returns error", goodPW, userCreateError, true},
		{"Well formed password creates a user", goodPW, validPasswordCreateUser, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestAuthClient()

			_, err := client.Signup(context.Background(), &server.SignupRequest{Password: test.password})
			if test.expectError && err == nil {
				t.Error("expected error")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			met := conn.ExpectationsWereMet()
			if met != nil {
				t.Errorf("redis expectations were not met\n%v", met)
			}
		})
	}
}

func TestServer_Finalise(t *testing.T) {
	goodPW := ""
	wrongPW := ""
	for i := 0; i < 16; i++ {
		goodPW += "a"
		wrongPW += "b"
	}

	validTotp, err := totp.GenerateCode("PAZXXJTBLKTDLHWTDPEGGMYNVR2LMJGC", time.Now())
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name        string
		id          string
		password    string
		totp        string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"Missing ID returns error", "", goodPW, "123456", noRedisCall, true},
		{"Missing Password returns error", "some-user@testing", "", "123456", noRedisCall, true},
		{"Missing totp code returns error", "some-user@testing", goodPW, "", noRedisCall, true},
		{"User does not exist", "some-user@testing", goodPW, "123456", noSuchUser, true},
		{"Password incorrect returns error", "some-user@testing", wrongPW, "123456", validUser, true},
		{"TOTP code wrong returns error", "some-user@testing", goodPW, "123456", validUser, true},
		{"Valid auth updates finalised", "some-user@testing", goodPW, validTotp, setFinalised, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestAuthClient()

			_, err := client.Finalise(context.Background(), &server.Auth{UserId: test.id, Password: test.password, Totp: test.totp})
			if test.expectError && err == nil {
				t.Error("expected error")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			met := conn.ExpectationsWereMet()
			if met != nil {
				t.Errorf("redis expectations were not met\n%v", met)
			}
		})
	}
}

func TestServer_Token(t *testing.T) {
	goodPW := ""
	wrongPW := ""
	for i := 0; i < 16; i++ {
		goodPW += "a"
		wrongPW += "b"
	}

	validTotp, err := totp.GenerateCode("PAZXXJTBLKTDLHWTDPEGGMYNVR2LMJGC", time.Now())
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name        string
		id          string
		password    string
		totp        string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"Missing ID returns error", "", goodPW, "123456", noRedisCall, true},
		{"Missing Password returns error", "some-user@testing", "", "123456", noRedisCall, true},
		{"Missing totp code returns error", "some-user@testing", goodPW, "", noRedisCall, true},
		{"User does not exist", "some-user@testing", goodPW, "123456", noSuchUser, true},
		{"Password incorrect returns error", "some-user@testing", wrongPW, "123456", validUser, true},
		{"TOTP code wrong returns error", "some-user@testing", goodPW, "123456", validUser, true},
		{"Valid auth mints and returns a token", "some-user@testing", goodPW, validTotp, addToken, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestAuthClient()

			_, err := client.Token(context.Background(), &server.Auth{UserId: test.id, Password: test.password, Totp: test.totp})
			if test.expectError && err == nil {
				t.Error("expected error")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			met := conn.ExpectationsWereMet()
			if met != nil {
				t.Errorf("redis expectations were not met\n%v", met)
			}
		})
	}
}

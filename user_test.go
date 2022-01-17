package main

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/prattle-chat/prattle-proxy/server"
	"github.com/rafaeljusto/redigomock/v3"
	"google.golang.org/grpc"
)

func newTestUserClient() (c server.UserClient) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return server.NewUserClient(conn)
}

func TestServer_AddPublicKey(t *testing.T) {
	for _, test := range []struct {
		name        string
		key         string
		user        string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"Remote user fails", "blahblahblah", "some-user@none", validPeeredToPeered, true},
		{"User with too many keys drops older keys", "foo", "some-user@none", addTokenSuccess, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestUserClient()

			_, err := client.AddPublicKey(addOperatorHeader(key(test.key).Auth(), test.user), &server.PublicKeyValue{
				Value: "some-key",
			})

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

func TestServer_PublicKeys(t *testing.T) {
	// These tests need to test the format of usernames, because we
	// can't do this in a stream interceptor
	for _, test := range []struct {
		name        string
		key         string
		user        string
		mocks       func(*redigomock.Conn)
		expect      []string
		expectError bool
	}{
		{"Valid local user, no keys", "foo", "some-user@testing", validTokenAndUser, []string{}, false},
		{"Valid local user, many keys", "foo", "some-user@testing", validTokenAndUserWithKeys, []string{"key-1", "key-2", "key-3", "key-4", "key-5"}, false},
		{"Valid remote user", "blahblahblah", "some-user@none", validPeeredToPeered, []string{"key-1"}, false},
		{"Bad user domain errors", "foo", "some-user", validTokenAndUser, []string{}, true},
		{"User is on a non-peered domain", "foo", "some-user@remote", validTokenAndUser, []string{}, true},
		{"Group public keys fails", "foo", "g:group@testing", validTokenAndUser, []string{}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestUserClient()

			pks, err := client.PublicKey(addOperatorHeader(key(test.key).Auth(), test.user), &server.PublicKeyRequest{
				Owner: &server.Subject{Id: test.user},
			})

			keys := make([]string, 0)
			for {
				var k *server.PublicKeyValue
				k, err = pks.Recv()
				if err != nil {
					if err == io.EOF {
						err = nil
					}

					break
				}

				keys = append(keys, k.Value)
			}

			if test.expectError && err == nil {
				t.Error("expected error")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if !reflect.DeepEqual(test.expect, keys) {
				t.Errorf("expected %#v, received %#v", test.expect, keys)
			}

			met := conn.ExpectationsWereMet()
			if met != nil {
				t.Errorf("redis expectations were not met\n%v", met)
			}
		})
	}
}

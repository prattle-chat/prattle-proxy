package main

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/prattle-chat/prattle-proxy/server"
	"github.com/rafaeljusto/redigomock/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func newTestMessageClient() (c server.MessagingClient) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return server.NewMessagingClient(conn)
}

func TestServer_Send(t *testing.T) {
	// These tests can ignore auth and so on; we cover that
	// in depth in our interceptor tests
	for _, test := range []struct {
		name        string
		key         string
		sender      string
		recipient   string
		behalfOf    string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"Valid local user can send to valid local user", "foo", "some-user@testing", "recipient@testing", "", validTokenUserRecipientAndPublish, false},
		{"Valid local user can send to valid remote user", "foo", "some-user@testing", "some-user@none", "", validTokenAndUser, false},
		{"Valid remote user can send to valid local user", "blahblahblah", "some-user@none", "recipient@testing", "", validPeeredUser, false},
		{"Valid remote user can send to valid remote user", "blahblahblah", "some-user@none", "some-user@none", "", validPeeredToPeered, false},
		{"Valid remote cannot sent on behalf of user", "blahblahblah", "some-user@none", "some-user@testing", "admin@testing", validPeeredToPeered, true},
		{"Valid remote cannot sent on behalf of group the user has no permission with", "blahblahblah", "some-user@none", "some-user@testing", "g:open@testing", validPeeredToPeered, true},

		{"Sending to groups fails", "foo", "some-user@testing", "g:group@testing", "", validTokenAndUser, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestMessageClient()
			_, err := client.Send(key(test.key).Auth(), &server.MessageWrapper{
				Sender:    test.sender,
				Recipient: test.recipient,
				For:       test.behalfOf,
			})

			if test.expectError && err == nil {
				t.Error("expected error")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %s", err)
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
		{"Valid remote user", "blahblahblah", "some-user@none", validPeeredUser, []string{"key-1"}, false},
		{"Bad user domain errors", "foo", "some-user", validTokenAndUser, []string{}, true},
		{"User is on a non-peered domain", "foo", "some-user@remote", validTokenAndUser, []string{}, true},
		{"Group public keys fails", "foo", "g:group@testing", validTokenAndUser, []string{}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestMessageClient()

			pks, err := client.PublicKey(key(test.key).Auth(), &server.Auth{UserId: test.user})

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
		})
	}
}

func TestServer_Subscribe(t *testing.T) {
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
		{"Valid user with messages", "foo", "some-user@testing", validPeerReceiveMessage, []string{}, false},
		{"Remote users cannot access messages", "foo", "some-user@none", validPeeredToPeered, []string{}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestMessageClient()

			sc, err := client.Subscribe(key(test.key).Auth(), new(emptypb.Empty))

			for {
				_, err = sc.Recv()

				break
			}

			if test.expectError && err == nil {
				t.Error("expected error")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %s", err)
			}
		})
	}
}

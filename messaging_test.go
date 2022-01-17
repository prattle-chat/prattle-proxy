package main

import (
	"context"
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
		recipient   *server.Subject
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"Valid local user can send to valid local user", "foo", "some-user@testing", &server.Subject{Id: "recipient@testing"}, validTokenUserRecipientAndPublish, false},
		{"Valid local user can send to valid remote user", "foo", "some-user@testing", &server.Subject{Id: "some-user@none"}, validTokenAndUser, false},
		{"Valid remote user can send to valid local user", "blahblahblah", "some-user@none", &server.Subject{Id: "recipient@testing"}, validPeeredUser, false},
		{"Valid remote user can send to valid remote user", "blahblahblah", "some-user@none", &server.Subject{Id: "some-user@none"}, validPeeredToPeered, false},
		{"Valid remote cannot send on behalf of user", "blahblahblah", "some-user@none", &server.Subject{Id: "some-user@testing", GroupId: "some-user@none"}, validPeeredToPeered, true},
		{"Valid remote cannot send on behalf of group the user has no permission with", "blahblahblah", "some-user@none", &server.Subject{Id: "some-user@testing", GroupId: "g:open@testing"}, validPeeredNoPermission, true},
		{"Valid remote cannot send on behalf of remote group the user has no permission with", "blahblahblah", "some-user@none", &server.Subject{Id: "some-user@testing", GroupId: "g:open@none"}, validPeeredToPeered, true},
		{"Valid remote cannot send on behalf of remote group on non-peered domain", "blahblahblah", "some-user@none", &server.Subject{Id: "some-user@testing", GroupId: "g:open@unknown"}, validPeeredToPeered, true},

		{"Sending to groups fails", "foo", "some-user@testing", &server.Subject{Id: "g:group@testing"}, validTokenAndUser, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestMessageClient()
			_, err := client.Send(addOperatorHeader(key(test.key).Auth(), test.sender), &server.MessageWrapper{
				Sender:    &server.Subject{Id: test.sender},
				Recipient: test.recipient,
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
		{"Remote users cannot access messages", "blahblahblah", "some-user@none", validPeeredToPeered, []string{}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestMessageClient()

			sc, err := client.Subscribe(addOperatorHeader(key(test.key).Auth(), test.user), new(emptypb.Empty))

			for {
				_, err = sc.Recv()

				break
			}

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

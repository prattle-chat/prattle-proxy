package main

import (
	"context"
	"testing"

	"github.com/prattle-chat/prattle-proxy/server"
	"github.com/rafaeljusto/redigomock/v3"
	"google.golang.org/grpc/metadata"
)

func addOperatorHeader(ctx context.Context, s string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, operatorIDHeader, s)
}

// TestServer_FederatedEndpoints_SendMessage tests the unary interceptor by trying to send some messages
// with varying (expected) success
func TestServer_FederatedEndpoints_SendMessage(t *testing.T) {
	for _, test := range []struct {
		name        string
		mocks       func(*redigomock.Conn)
		ctx         context.Context
		request     *server.MessageWrapper
		expectError string
	}{
		{"missing auth", noRedisCall, context.Background(), &server.MessageWrapper{}, "rpc error: code = Unauthenticated desc = missing/ unreadble bearer token"},
		{"redis breaks on auth lookup", redisErrorOnAuth, key("foo").Auth(), &server.MessageWrapper{}, "rpc error: code = Unavailable desc = an internal systems error occurred"},
		{"token exists but user does not", validTokenInvalidUser, key("foo").Auth(), &server.MessageWrapper{}, "rpc error: code = Unauthenticated desc = missing/ unreadble bearer token"},
		{"token and user exist, but user has not finished signing up", validTokenNonFinaliseduser, key("foo").Auth(), &server.MessageWrapper{}, "rpc error: code = Unauthenticated desc = must finalise signup"},
		{"valid user missing recipient", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Sender: &server.Subject{Id: "some-user@testing"}}, "rpc error: code = Unavailable desc = missing/ poorly formed input"},
		{"valid user, dodgy recipient", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient"}, Sender: &server.Subject{Id: "some-user@testing"}}, "rpc error: code = Unavailable desc = missing/ poorly formed input"},
		{"valid user, local recipient", validTokenUserRecipientAndPublish, key("foo").Auth(), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient@testing"}, Sender: &server.Subject{Id: "some-user@testing"}}, ""},
		{"valid user, recipient on non-peered domain", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient@remote"}, Sender: &server.Subject{Id: "some-user@testing"}}, "rpc error: code = NotFound desc = recipient is on a non-peered domain"},
		{"valid user, recipient on peered domain", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient@none"}, Sender: &server.Subject{Id: "some-user@testing"}}, ""},
		{"peered user, recipient on local domain", validPeeredUser, addOperatorHeader(key("blahblahblah").Auth(), "some-user@none"), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient@testing"}, Sender: &server.Subject{Id: "some-user@none"}}, ""},
		{"non-peered user, recipient on local domain", invalidPeeredUser, key("foobarbaz").Auth(), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient@testing"}, Sender: &server.Subject{Id: "some-user@none"}}, "rpc error: code = NotFound desc = recipient is on a non-peered domain"},
		{"peered user, peered recipient", validPeeredToPeered, addOperatorHeader(key("blahblahblah").Auth(), "some-user@none"), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient@none"}, Sender: &server.Subject{Id: "some-user@none"}}, ""},
		{"peered user tries to spoof sender", validPeeredToPeered, addOperatorHeader(key("blahblahblah").Auth(), "some-user@none"), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient@none"}, Sender: &server.Subject{Id: "admin@third-party"}}, "rpc error: code = InvalidArgument desc = mismatch between sender field and owner of token"},
		{"peered user missing domain", validPeeredToPeered, addOperatorHeader(key("blahblahblah").Auth(), "admin"), &server.MessageWrapper{Recipient: &server.Subject{Id: "recipient@none"}, Sender: &server.Subject{Id: "admin"}}, "rpc error: code = Unavailable desc = missing/ poorly formed input"},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestMessageClient()

			_, err := client.Send(test.ctx, test.request)
			if test.expectError == "" && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if test.expectError != "" && err == nil {
				t.Error("expected error")
			} else if test.expectError != "" && test.expectError != err.Error() {
				t.Errorf("expected %q, received %q", test.expectError, err.Error())
			}

			met := conn.ExpectationsWereMet()
			if met != nil {
				t.Errorf("redis expectations were not met\n%v", met)
			}
		})
	}
}

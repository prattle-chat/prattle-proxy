package main

import (
	"context"
	"testing"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/grpc"
)

// TestServer_FederatedEndpoints_SendMessage tests the unary interceptor by trying to send some messages
// with varying (expected) success
func TestServer_FederatedEndpoints_SendMessage(t *testing.T) {
	for _, test := range []struct {
		name        string
		mocks       func()
		ctx         context.Context
		request     *server.MessageWrapper
		expectError string
	}{
		{"missing auth", noRedisCall, context.Background(), &server.MessageWrapper{}, "rpc error: code = Unauthenticated desc = missing/ unreadble bearer token"},
		{"redis breaks on auth lookup", redisErrorOnAuth, key("foo").Auth(), &server.MessageWrapper{}, "rpc error: code = Unavailable desc = an internal systems error occurred"},
		{"token exists but user does not", validTokenInvalidUser, key("foo").Auth(), &server.MessageWrapper{}, "rpc error: code = Unauthenticated desc = missing/ unreadble bearer token"},
		{"token and user exist, but user has not finished signing up", validTokenNonFinaliseduser, key("foo").Auth(), &server.MessageWrapper{}, "rpc error: code = Unauthenticated desc = must finalise signup"},
		{"valid user missing sender", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Recipient: "some-user@testing"}, "rpc error: code = Unavailable desc = missing/ poorly formed input"},
		{"valid user missing recipient", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Sender: "some-user@testing"}, "rpc error: code = Unavailable desc = missing/ poorly formed input"},
		{"valid user tries to spoof sender", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Recipient: "recipient@testing", Sender: "admin@prattle"}, "rpc error: code = InvalidArgument desc = mismatch between sender field and owner of token"},
		{"valid user, dodgy recipient", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Recipient: "recipient", Sender: "some-user@testing"}, "rpc error: code = Unavailable desc = missing/ poorly formed input"},
		{"valid user, local recipient", validTokenUserRecipientAndPublish, key("foo").Auth(), &server.MessageWrapper{Recipient: "recipient@testing", Sender: "some-user@testing"}, ""},
		{"valid user, recipient on non-peered domain", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Recipient: "recipient@remote", Sender: "some-user@testing"}, "rpc error: code = NotFound desc = recipient is on a non-peered domain"},
		{"valid user, recipient on peered domain", validTokenAndUser, key("foo").Auth(), &server.MessageWrapper{Recipient: "recipient@none", Sender: "some-user@testing"}, ""},
		{"peered user, recipient on local domain", validPeeredUser, key("blahblahblah").Auth(), &server.MessageWrapper{Recipient: "recipient@testing", Sender: "some-user@none"}, ""},
		{"non-peered user, recipient on local domain", invalidPeeredUser, key("foobarbaz").Auth(), &server.MessageWrapper{Recipient: "recipient@testing", Sender: "some-user@none"}, "rpc error: code = NotFound desc = recipient is on a non-peered domain"},
		{"peered user, peered recipient", validPeeredToPeered, key("blahblahblah").Auth(), &server.MessageWrapper{Recipient: "recipient@none", Sender: "some-user@none"}, ""},
		{"peered user tries to spoof sender domain", validPeeredToPeered, key("blahblahblah").Auth(), &server.MessageWrapper{Recipient: "recipient@none", Sender: "admin@third-party"}, "rpc error: code = InvalidArgument desc = mismatch between sender domain and the domain this federated peer belongs to"},
		{"peered user missing domain", validPeeredToPeered, key("blahblahblah").Auth(), &server.MessageWrapper{Recipient: "recipient@none", Sender: "admin"}, "rpc error: code = Unavailable desc = missing/ poorly formed input"},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks()
			newTestServer(NewDummyRedis(redigoMockConn))

			conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()

			client := server.NewMessagingClient(conn)

			_, err = client.Send(test.ctx, test.request)
			if test.expectError == "" && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if test.expectError != "" && err == nil {
				t.Error("expected error")
			} else if test.expectError != "" && test.expectError != err.Error() {
				t.Errorf("expected %q, received %q", test.expectError, err.Error())
			}

			met := redigoMockConn.ExpectationsWereMet()
			if met != nil {
				t.Errorf("redis expectations were not met\n%v", met)
			}
		})
	}
}

package main

import (
	"context"
	"testing"

	"github.com/prattle-chat/prattle-proxy/server"
	"github.com/rafaeljusto/redigomock/v3"
	"google.golang.org/grpc"
)

func newTestGroupsClient() (c server.GroupsClient) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return server.NewGroupsClient(conn)
}

func TestServer_Create(t *testing.T) {
	oldMinter := minter
	defer func() {
		minter = oldMinter
	}()

	minter = dummyMinter

	for _, test := range []struct {
		name        string
		key         string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"Failure to mint a group ID returns error in time", "foo", groupIdFound, true},
		{"Failure to validate a group ID returns error", "foo", groupIdLookupErrors, true},
		{"Failure to store a group returns error", "foo", groupAddError, true},
		{"Happy path", "foo", groupAddSuccess, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestGroupsClient()

			_, err := client.Create(key(test.key).Auth(), new(server.Group))
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

func TestServer_Join(t *testing.T) {
	for _, test := range []struct {
		name        string
		key         string
		group       string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"request with missing groupID returns error", "foo", "", validTokenAndUser, true},
		{"trying to join a closed groupreturns error", "foo", "g:closed@testing", closedGroup, true},
		{"trying to join a missing group returns error", "foo", "g:closed@testing", missingGroup, true},
		{"errors when trying to join returns an error", "foo", "g:open@testing", groupJoinOnErroringGroup, true},
		{"happy path", "foo", "g:open@testing", groupJoinSuccess, false},
		{"external group", "foo", "g:group@none", validTokenAndUser, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestGroupsClient()

			_, err := client.Join(key(test.key).Auth(), &server.JoinRequest{
				GroupId: test.group,
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

func TestServer_Invite(t *testing.T) {
	for _, test := range []struct {
		name        string
		key         string
		user        string
		group       string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"inviting to a group without being an owner fails", "foo", "some-user@testing", "g:closed@testing", closedGroup, true},
		{"inviting a non-existent user to a group fails", "foo", "another-user@testing", "g:open@testing", groupInviteeNotExist, true},
		{"inviting to a non-existent group fails", "foo", "some-user@testing", "g:closed@testing", missingGroup, true},
		{"inviting to a group with permission succedes", "foo", "some-user@testing", "g:open@testing", validGroupWithMember, false},
		{"external group", "foo", "some-user@testing", "g:group@none", validTokenAndUser, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestGroupsClient()

			_, err := client.Invite(key(test.key).Auth(), &server.InviteRequest{
				Invitee: test.user,
				GroupId: test.group,
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

func TestServer_PromoteUser(t *testing.T) {
	for _, test := range []struct {
		name        string
		key         string
		user        string
		group       string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"local group, caller has no permission", "foo", "some-user@testing", "g:open@testing", validGroup, true},
		{"local group, caller not even a member", "foo", "some-user@testing", "g:open@testing", validGroupNotMember, true},
		{"local group, caller is owner, promotee is not a member", "foo", "another@testing", "g:open@testing", validGroupWithMember, true},
		{"local group, caller is owner, promotee is already owner", "foo", "some-user@testing", "g:open@testing", validGroupWithMember, false},
		{"local group, caller is owner, promotee is just a member", "foo", "another@testing", "g:open@testing", validGroupAndMemberWithOwner, false},
		{"external group calls federated client", "foo", "some-user@testing", "g:group@none", validTokenAndUser, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestGroupsClient()

			_, err := client.PromoteUser(key(test.key).Auth(), &server.PromoteRequest{
				Promotee: test.user,
				GroupId:  test.group,
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

func TestServer_DemoteUser(t *testing.T) {
	for _, test := range []struct {
		name        string
		key         string
		user        string
		group       string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"local group, caller has no permission", "foo", "some-user@testing", "g:open@testing", validGroup, true},
		{"local group, caller not even a member", "foo", "some-user@testing", "g:open@testing", validGroupNotMember, true},
		{"local group, caller is owner, demotee is not a member", "foo", "another@testing", "g:open@testing", validGroupWithMember, true},
		{"local group, caller is owner, demotee is owner", "foo", "some-user@testing", "g:open@testing", validGroupWithMemberDemote, false},
		{"local group, caller is owner, demotee is just a member", "foo", "some-user@testing", "g:open@testing", validGroupWithMemberDemote, false},
		{"external group calls federated client", "foo", "some-user@testing", "g:group@none", validTokenAndUser, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestGroupsClient()

			_, err := client.DemoteUser(key(test.key).Auth(), &server.DemoteRequest{
				Demotee: test.user,
				GroupId: test.group,
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

func TestServer_Info(t *testing.T) {
	for _, test := range []struct {
		name        string
		key         string
		group       string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"accessing a group without being a member fails", "foo", "g:closed@testing", closedGroup, true},
		{"accessing a non-existant group fails", "foo", "g:closed@testing", missingGroup, true},
		{"accessing a group with permission succedes", "foo", "g:open@testing", validGroup, false},
		{"external group", "foo", "g:group@none", validTokenAndUser, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestGroupsClient()

			_, err := client.Info(key(test.key).Auth(), &server.InfoRequest{
				GroupId: test.group,
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

func TestServer_Leave(t *testing.T) {
	for _, test := range []struct {
		name        string
		key         string
		group       string
		mocks       func(*redigomock.Conn)
		expectError bool
	}{
		{"leaving a group without being a member fails", "foo", "g:closed@testing", closedGroup, true},
		{"leaving a non-existant group fails", "foo", "g:closed@testing", missingGroup, true},
		{"leaving a group with permission succedes", "foo", "g:open@testing", validGroupWithMemberDemote, false},
		{"external group", "foo", "g:group@none", validTokenAndUser, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mocks(conn)

			newTestServer(NewDummyRedis(conn))

			client := newTestGroupsClient()

			_, err := client.Leave(key(test.key).Auth(), &server.LeaveRequest{
				GroupId: test.group,
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

func TestHasPermission(t *testing.T) {
	closedGroup := Group{
		Owners:      []string{"owner"},
		Members:     []string{"owner", "member1", "member2"},
		IsOpen:      false,
		IsBroadcast: false,
	}

	openGroup := Group{
		Owners:      []string{"owner"},
		Members:     []string{"owner", "member1", "member2"},
		IsOpen:      true,
		IsBroadcast: false,
	}

	broadcastClosedGroup := Group{
		Owners:      []string{"owner"},
		Members:     []string{"owner", "member1", "member2"},
		IsOpen:      false,
		IsBroadcast: true,
	}

	broadcastOpenGroup := Group{
		Owners:      []string{"owner"},
		Members:     []string{"owner", "member1", "member2"},
		IsOpen:      true,
		IsBroadcast: true,
	}

	owner := User{Id: "owner"}
	member1 := User{Id: "member1"}
	nonMember := User{Id: "non-member"}

	for _, test := range []struct {
		name   string
		u      User
		g      Group
		op     GroupOperation
		expect bool
	}{
		// Reading groups
		{"Owner can read closed group", owner, closedGroup, groupRead, true},
		{"Owner can read open group", owner, openGroup, groupRead, true},
		{"Owner can read closed broadcast group", owner, broadcastClosedGroup, groupRead, true},
		{"Owner can read open broadcast group", owner, broadcastOpenGroup, groupRead, true},

		{"Member can read closed group", member1, closedGroup, groupRead, true},
		{"Member can read open group", member1, openGroup, groupRead, true},
		{"Member can not read closed broadcast group", member1, broadcastClosedGroup, groupRead, false},
		{"Member can not read open broadcast group", member1, broadcastOpenGroup, groupRead, false},

		// Posting to groups
		{"Owner can post to closed group", owner, closedGroup, groupPost, true},
		{"Owner can post to open group", owner, openGroup, groupPost, true},
		{"Owner can post to closed broadcast group", owner, broadcastClosedGroup, groupPost, true},
		{"Owner can post to open broadcast group", owner, broadcastOpenGroup, groupPost, true},

		{"Member can post to closed group", member1, closedGroup, groupPost, true},
		{"Member can post to open group", member1, openGroup, groupPost, true},
		{"Member can not post to closed broadcast group", member1, broadcastClosedGroup, groupPost, false},
		{"Member can not post to open broadcast group", member1, broadcastOpenGroup, groupPost, false},

		// Joining groups
		{"Non-member can not join closed group", nonMember, closedGroup, groupJoin, false},
		{"Non-member can join open group", nonMember, openGroup, groupJoin, true},
		{"Non-member can not join closed broadcast group", nonMember, broadcastClosedGroup, groupJoin, false},
		{"Non-member can join open broadcast group", nonMember, broadcastOpenGroup, groupJoin, true},

		// Modifying groups
		{"Owner can modify closed group", owner, closedGroup, groupModify, true},
		{"Owner can modify open group", owner, openGroup, groupModify, true},
		{"Owner can modify closed broadcast group", owner, broadcastClosedGroup, groupModify, true},
		{"Owner can modify open broadcast group", owner, broadcastOpenGroup, groupModify, true},

		{"Member can not modify closed group", member1, closedGroup, groupModify, false},
		{"Member can not modify open group", member1, openGroup, groupModify, false},
		{"Member can not modify closed broadcast group", member1, broadcastClosedGroup, groupModify, false},
		{"Member can not modify open broadcast group", member1, broadcastOpenGroup, groupModify, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := HasPermission(test.u, test.g, test.op)
			if test.expect != received {
				t.Errorf("expected %v, received %v", test.expect, received)
			}
		})
	}
}

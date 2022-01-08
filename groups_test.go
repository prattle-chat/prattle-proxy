package main

import (
	"testing"
)

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

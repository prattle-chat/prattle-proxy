package main

import (
	"context"
	"regexp"

	"github.com/prattle-chat/prattle-proxy/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

/**
  On confusing/ ambiguous method calls:

  Our gRPC calls are all namespaced to things like, in this instance, 'Groups'
  This allows us to call 'Groups.Create', which makes for a nicer client experience.

  However... we use a single Server which implements all of these namespaces- we have
  lots of reusable functions and calls to redis that make sense in a single Server; rather
  than copying references or even structs about.

  The issue becomes, then, implementing functions like the above 'Groups.Create' where
  Server needs to implement a function ambiguously called 'Create'. There's no real
  hint that this function is to create a group.

  So, let this long winded comment (plus comments above each function) go some way to
  provide that hint.
  **/

type GroupOperation uint8

const (
	groupRead GroupOperation = iota
	groupPost
	groupJoin
	groupModify
	groupLeave
)

var (
	groupFormat = "g:%s"
	groupRegexp = regexp.MustCompile("^g:.+")
)

// Create is called byt Group.Create
func (s Server) Create(ctx context.Context, in *server.Group) (out *server.Group, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)

	// It's hard to see how this would happen, beyond a badly
	// configured server, but it's worth guarding against
	if !m.Sender.IsLocal {
		err = inputError

		return
	}

	gid, err := s.mintGroupID()
	if err != nil {
		err = generalError

		return
	}

	err = s.redis.AddGroup(gid, m.Sender.Id, in.IsOpen, in.IsBroadcast)
	if err != nil {
		err = generalError

		return
	}

	out = in
	out.Id = gid
	out.Members = []string{m.Sender.Id}
	out.Owners = []string{m.Sender.Id}

	return
}

// Join is called by Group.Join
func (s Server) Join(ctx context.Context, in *server.GroupUser) (out *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)

	out = new(emptypb.Empty)

	if m.Recipient.IsLocal {
		_, err = s.groupPermitted(in.GroupId, m.Sender.Id, groupJoin)
		if err != nil {
			return
		}

		if m.Sender.IsLocal {
			err = s.redis.JoinGroup(in.GroupId, m.Sender.Id)
			if err != nil {
				err = inputError
			}
		}

		return
	}

	// External group
	in.UserId = m.Sender.Id

	err = m.Recipient.PeerConnection.JoinGroup(in)

	return
}

// Info is called by Group.Info
func (s Server) Info(ctx context.Context, in *server.GroupUser) (out *server.Group, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)

	if m.Recipient.IsLocal {
		var g Group

		g, err = s.groupPermitted(in.GroupId, m.Sender.Id, groupRead)
		if err != nil {
			return
		}

		out = &server.Group{
			Id:          g.Id,
			Owners:      g.Owners,
			Members:     g.Members,
			IsOpen:      g.IsOpen,
			IsBroadcast: g.IsBroadcast,
		}
		return
	}

	// External group
	in.UserId = m.Sender.Id
	return m.Recipient.PeerConnection.GroupInfo(in)
}

// Invite is called by Group.Invite.
//
// In this, the values of GroupUser mean:
//
//  1. GroupID - The group we're inviting a user to
//  2. UserID  - The user we're inviting to the specified group
//  3. For     - The user who is doing the inviting
//
// Thus, Metadata.Sender.Id should be the same as 'For' on incoming calls from
// peers, and should be set on calls *to* peers
func (s Server) Invite(ctx context.Context, in *server.GroupUser) (_ *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)

	if m.Recipient.IsLocal {
		return new(emptypb.Empty), nil
	}

	// External group
	in.For = m.Sender.Id
	return new(emptypb.Empty), m.Recipient.PeerConnection.InviteToGroup(in)
}

func (s Server) PromoteUser(ctx context.Context, in *server.GroupUser) (_ *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)

	if m.Recipient.IsLocal {
		return new(emptypb.Empty), nil
	}

	// External group
	return new(emptypb.Empty), m.Recipient.PeerConnection.PromoteUser(in)
}

func (s Server) DemoteUser(ctx context.Context, in *server.GroupUser) (_ *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)

	if m.Recipient.IsLocal {
		return new(emptypb.Empty), nil
	}

	// External group
	return new(emptypb.Empty), m.Recipient.PeerConnection.DemoteUser(in)
}

func (s Server) LeaveGroup(ctx context.Context, in *server.GroupUser) (_ *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)

	if m.Recipient.IsLocal {
		return new(emptypb.Empty), nil
	}

	// External group
	in.UserId = m.Sender.Id
	return new(emptypb.Empty), m.Recipient.PeerConnection.LeaveGroup(in)
}

func (s Server) groupPermitted(groupId, actorId string, op GroupOperation) (g Group, err error) {
	g, err = s.redis.Group(groupId)
	if err != nil {
		err = badGroupError

		return
	}

	if !HasPermission(User{Id: actorId}, g, op) {
		err = badGroupError
	}

	return
}

func HasPermission(u User, g Group, o GroupOperation) bool {
	// any user has permission to leave any group
	if o == groupLeave {
		return true
	}

	// if u in group owners, then always true
	if contains(g.Owners, u.Id) {
		return true
	}

	if !contains(g.Members, u.Id) && o != groupJoin {
		// join is the only thing a non-member can do;
		// if not a member and not looking to join then
		// automatically that user has no permission to
		// do whatever it is they're doing.
		//
		// Note: a user *still* might not be allowed to join
		// based on other permissions, which is why we return
		// false here, rather than assuming permission is
		// granted
		return false
	}

	// if g is open and u wants to join, grant permission
	if o == groupJoin && g.IsOpen {
		return true
	}

	// A non-owner can read and/or post to a non-broadcast group
	// they're a member of
	if !g.IsBroadcast && (o == groupRead || o == groupPost) {
		return true
	}

	// If we get this far then operation is not permitted
	return false
}

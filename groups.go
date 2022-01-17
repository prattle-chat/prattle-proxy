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

	gid, err := s.mintGroupID()
	if err != nil {
		err = generalError

		return
	}

	err = s.redis.AddGroup(gid, m.Operator.Id, in.IsOpen, in.IsBroadcast)
	if err != nil {
		err = generalError

		return
	}

	out = in
	out.Id = gid
	out.Members = []string{m.Operator.Id}
	out.Owners = []string{m.Operator.Id}

	return
}

// Join is called by Group.Join
func (s Server) Join(ctx context.Context, in *server.JoinRequest) (out *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)
	group := m.Operand

	out = new(emptypb.Empty)

	if group.IsLocal {
		// Can we join this group?
		_, err = s.groupPermitted(in.GroupId, m.Operator.Id, groupJoin)
		if err != nil {
			return
		}

		err = s.redis.JoinGroup(in.GroupId, m.Operator.Id)
		if err != nil {
			err = inputError
		}

		return
	}

	err = group.PeerConnection.JoinGroup(m.Operator.Id, in)

	return
}

// Info is called by Group.Info
func (s Server) Info(ctx context.Context, in *server.InfoRequest) (out *server.Group, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)
	group := m.Operand

	if group.IsLocal {
		var g Group

		g, err = s.groupPermitted(group.Id, m.Operator.Id, groupRead)
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

	return group.PeerConnection.GroupInfo(m.Operator.Id, in)
}

// Invite is called by Group.Invite.
func (s Server) Invite(ctx context.Context, in *server.InviteRequest) (_ *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)
	group := m.Operand

	// if invitee is local, check exists before doing anything else
	d, err := domain(in.Invitee)
	if err != nil {
		return
	}

	if d == s.config.DomainName {
		var u User

		u, err = s.redis.loadUser(in.Invitee)
		if err != nil || u.Id == "" {
			err = badUserError

			return
		}
	}

	if group.IsLocal {
		var g Group

		g, err = s.groupPermitted(group.Id, m.Operator.Id, groupModify)
		if err != nil {
			return
		}

		if !contains(g.Members, in.Invitee) {
			g.Members = append(g.Members, in.Invitee)

			err = s.redis.saveGroup(g)
		}
		return new(emptypb.Empty), err
	}

	// External group
	return new(emptypb.Empty), group.PeerConnection.InviteToGroup(m.Operator.Id, in)
}

func (s Server) PromoteUser(ctx context.Context, in *server.PromoteRequest) (_ *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)
	group := m.Operand

	if group.IsLocal {
		var g Group

		g, err = s.groupPermitted(group.Id, m.Operator.Id, groupModify)
		if err != nil {
			return
		}

		switch {
		case contains(g.Owners, in.Promotee):
			// noop

		case contains(g.Members, in.Promotee):
			g.Owners = append(g.Owners, in.Promotee)
			err = s.redis.saveGroup(g)

		default:
			err = badUserError
		}

		return new(emptypb.Empty), err
	}

	// External group
	return new(emptypb.Empty), group.PeerConnection.PromoteUser(m.Operator.Id, in)
}

func (s Server) DemoteUser(ctx context.Context, in *server.DemoteRequest) (_ *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)
	group := m.Operand

	if group.IsLocal {
		var g Group

		g, err = s.groupPermitted(group.Id, m.Operator.Id, groupModify)
		if err != nil {
			return
		}

		switch {
		case contains(g.Owners, in.Demotee):
			g.Owners = remove(g.Owners, in.Demotee)

		case contains(g.Members, in.Demotee):
			g.Members = remove(g.Members, in.Demotee)

		default:
			err = badUserError

			return
		}

		err = s.redis.saveGroup(g)

		return new(emptypb.Empty), err
	}

	// External group
	return new(emptypb.Empty), group.PeerConnection.DemoteUser(m.Operator.Id, in)
}

func (s Server) Leave(ctx context.Context, in *server.LeaveRequest) (out *emptypb.Empty, err error) {
	m := ctx.Value(MetadataKey{}).(*Metadata)
	group := m.Operand

	out = new(emptypb.Empty)

	if group.IsLocal {
		var g Group

		// Can we leave this group?
		g, err = s.groupPermitted(in.GroupId, m.Operator.Id, groupLeave)
		if err != nil {
			return
		}

		if !contains(g.Owners, m.Operator.Id) && !contains(g.Members, m.Operator.Id) {
			err = badGroupError

			return
		}

		err = s.redis.RemoveFromGroup(in.GroupId, m.Operator.Id)
		if err != nil {
			err = inputError
		}

		return
	}

	err = group.PeerConnection.LeaveGroup(m.Operator.Id, in)

	return
}

func (s Server) groupPermitted(groupId, actorId string, op GroupOperation) (g Group, err error) {
	g, err = s.redis.Group(groupId)
	if err != nil {
		err = badGroupError

		return
	}

	if g.Id == "" {
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

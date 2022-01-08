package main

import (
	"context"
	"fmt"
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
	groupRegexp = regexp.MustCompile("g:.+")
)

// Create is called byt Group.Create
func (s Server) Create(ctx context.Context, in *server.Group) (out *server.Group, err error) {
	u, err := s.userFromContext(ctx)
	if err != nil {
		return
	}

	id, err := s.mintID()
	if err != nil {
		return
	}

	groupId := fmt.Sprintf(groupFormat, id)
	err = s.redis.AddGroup(groupId, u.Id, in.IsOpen, in.IsBroadcast)
	if err != nil {
		err = generalError
	}

	out = in
	out.Id = groupId
	out.Owners = []string{u.Id}
	out.Members = []string{u.Id}

	return in, err
}

// Join is called by Grouo.Join
func (s Server) Join(ctx context.Context, in *server.Group) (out *emptypb.Empty, err error) {
	u, _, err := s.UserAndGroup(ctx, in, groupJoin)
	if err != nil {
		return
	}

	if s.redis.JoinGroup(in.Id, u.Id) != nil {
		err = generalError
	}

	out = new(emptypb.Empty)
	return
}

// Info is called by Group.Info
func (s Server) Info(ctx context.Context, in *server.Group) (out *server.Group, err error) {
	_, g, err := s.UserAndGroup(ctx, in, groupRead)
	if err != nil {
		return
	}

	out = in
	out.Owners = g.Owners
	out.Members = g.Members
	out.IsBroadcast = g.IsBroadcast
	out.IsOpen = g.IsOpen

	return
}

// Invite is called by Group.Invite
func (s Server) Invite(ctx context.Context, in *server.GroupUser) (_ *emptypb.Empty, err error) {
	_, _, err = s.UserAndGroup(ctx, &server.Group{Id: in.GroupId}, groupModify)
	if err != nil {
		return
	}

	if s.redis.JoinGroup(in.GroupId, in.UserId) != nil {
		err = generalError
	}

	return new(emptypb.Empty), err
}

func (s Server) PromoteUser(ctx context.Context, in *server.GroupUser) (_ *emptypb.Empty, err error) {
	_, g, err := s.UserAndGroup(ctx, &server.Group{Id: in.GroupId}, groupModify)
	if err != nil {
		return
	}

	if !contains(g.Members, in.UserId) {
		err = notMemberError

		return
	}

	if s.redis.PromoteUser(g.Id, in.UserId) != nil {
		err = generalError
	}

	return new(emptypb.Empty), err
}

func (s Server) DemoteUser(ctx context.Context, in *server.GroupUser) (_ *emptypb.Empty, err error) {
	_, g, err := s.UserAndGroup(ctx, &server.Group{Id: in.GroupId}, groupModify)
	if err != nil {
		return
	}

	if s.redis.DemoteUser(g.Id, in.UserId) != nil {
		err = generalError
	}

	return new(emptypb.Empty), err
}

func (s Server) LeaveGroup(ctx context.Context, in *server.Group) (_ *emptypb.Empty, err error) {
	u, g, err := s.UserAndGroup(ctx, in, groupLeave)
	if err != nil {
		return
	}

	if s.redis.RemoveFromGroup(g.Id, u.Id) != nil {
		err = generalError
	}

	return new(emptypb.Empty), err
}

// UserAndGroup verifies the user token used in a request, verifies the group
// exists, and then, finally, that the user specified has permission to perform
// the specified operation against this group.
//
// The error returned by this function can be passed directly to the client
func (s Server) UserAndGroup(ctx context.Context, in *server.Group, op GroupOperation) (u User, g Group, err error) {
	u, err = s.userFromContext(ctx)
	if err != nil {
		return
	}

	g, err = s.redis.Group(in.Id)
	if err != nil {
		err = badGroupError

		return
	}

	if !HasPermission(u, g, op) {
		err = badGroupError

		return
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

/**
accesses:

  broadcast group: only owners can access group info, public keys
  normal group:    all members can access info and keys


  **/

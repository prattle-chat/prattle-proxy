package main

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/go-redis/redis/v8"
)

var (
	tokenIDHMKey = "tokensToIDs"
)

type Redis struct {
	c *redis.Client
}

type User struct {
	Id          string   `mapstructure:"id"`
	Password    string   `mapstructure:"password"`
	Seed        string   `mapstructure:"seed"`
	Finalised   bool     `mapstructure:"finalised"`
	Tokens      []string `mapstructure:"tokens"`
	Connections []string `mapstructure:"connections"`
	PublicKeys  []string `mapstructure:"public_keys"`
}

func (r Redis) loadUser(id string) (u User, err error) {
	err = r.load(id, &u)

	return
}

func (r Redis) saveUser(u User) (err error) {
	return r.save(u.Id, u)
}

type Group struct {
	Id          string
	Owners      []string
	Members     []string //owner exists in both owners and members; we use members to send messages
	IsOpen      bool
	IsBroadcast bool
}

func (r Redis) loadGroup(id string) (g Group, err error) {
	err = r.load(id, &g)

	return
}

func (r Redis) saveGroup(g Group) (err error) {
	return r.save(g.Id, g)
}

func NewRedis(addr string) (r Redis, err error) {
	r.c = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return
}

func (r Redis) AddUser(id, seed, password string) (err error) {
	u := User{
		Id:          id,
		Seed:        seed,
		Password:    password,
		Tokens:      make([]string, 0),
		Connections: make([]string, 0),
		PublicKeys:  make([]string, 0),
	}

	return r.saveUser(u)
}

func (r Redis) AddToken(id, token string) (err error) {
	u, err := r.loadUser(id)
	if err != nil {
		return
	}

	u.Tokens = append(u.Tokens, token)

	err = r.saveUser(u)
	if err != nil {
		return
	}

	return r.c.HSet(context.Background(), tokenIDHMKey, token, id).Err()
}

func (r Redis) AddPublicKey(id, pk string) (err error) {
	u, err := r.loadUser(id)
	if err != nil {
		return
	}

	u.PublicKeys = append(u.PublicKeys, pk)

	return r.saveUser(u)
}

func (r Redis) GetPublicKeys(id string) ([]string, error) {
	u, err := r.loadUser(id)
	if err != nil {
		return nil, err
	}

	return u.PublicKeys, nil
}

func (r Redis) MarkFinalised(id string) (err error) {
	u, err := r.loadUser(id)
	if err != nil {
		return
	}

	u.Finalised = true

	return r.saveUser(u)
}

func (r Redis) GetTOTPSeed(id string) (s string, err error) {
	u, err := r.loadUser(id)
	if err != nil {
		return
	}

	return u.Seed, nil
}

func (r Redis) GetPasswordHash(id string) (s string, err error) {
	u, err := r.loadUser(id)
	if err != nil {
		return
	}

	return u.Password, nil
}

func (r Redis) IDExists(id string) bool {
	u, err := r.loadUser(id)

	return err == nil && u.Id != ""
}

func (r Redis) UserByToken(token string) (u User, err error) {
	res, err := r.c.HGet(context.Background(), tokenIDHMKey, token).Result()
	if err != nil {
		return
	}

	return r.loadUser(res)
}

func (r Redis) Messages(id string) <-chan *redis.Message {
	return r.c.Subscribe(context.Background(), id).Channel()
}

func (r Redis) WriteMessage(recipient string, payload []byte) error {
	return r.c.Publish(context.Background(), recipient, payload).Err()
}

func (r Redis) AddGroup(id, owner string, open, broadcast bool) (err error) {
	g := Group{
		Id:          id,
		Owners:      []string{owner},
		Members:     []string{owner},
		IsOpen:      open,
		IsBroadcast: broadcast,
	}

	return r.saveGroup(g)
}

func (r Redis) Group(id string) (Group, error) {
	return r.loadGroup(id)
}

func (r Redis) JoinGroup(group, user string) (err error) {
	g, err := r.loadGroup(group)
	if err != nil {
		return
	}

	g.Members = append(g.Members, user)

	return r.saveGroup(g)
}

func (r Redis) PromoteUser(group, user string) (err error) {
	g, err := r.loadGroup(group)
	if err != nil {
		return
	}

	g.Owners = append(g.Members, user)

	return r.saveGroup(g)
}

func (r Redis) DemoteUser(group, user string) (err error) {
	g, err := r.loadGroup(group)
	if err != nil {
		return
	}

	if !contains(g.Owners, user) {
		g.Owners = removeElem(g.Owners, user)
	} else {
		g.Members = removeElem(g.Members, user)
	}

	return r.saveGroup(g)
}

func (r Redis) RemoveFromGroup(group, user string) (err error) {
	g, err := r.loadGroup(group)
	if err != nil {
		return
	}

	g.Owners = removeElem(g.Owners, user)
	g.Members = removeElem(g.Members, user)

	return r.saveGroup(g)
}

func (r Redis) load(id string, o interface{}) (err error) {
	res, err := r.c.Get(context.Background(), id).Result()
	if err != nil {
		return
	}

	b := bytes.NewBufferString(res)

	dec := gob.NewDecoder(b)
	err = dec.Decode(o)

	return
}

func (r Redis) save(id string, o interface{}) (err error) {
	b := bytes.Buffer{}

	enc := gob.NewEncoder(&b)
	err = enc.Encode(o)
	if err != nil {
		return
	}

	return r.c.Set(context.Background(), id, b.String(), 0).Err()
}

func removeElem(ss []string, s string) (out []string) {
	out = make([]string, 0)

	for _, elem := range ss {
		if s != elem {
			out = append(out, elem)
		}
	}

	return
}

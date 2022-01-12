package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	tokenIDHMKey = "tokensToIDs"
)

type Redis struct {
	pool *redis.Pool
}

type stringSlice []string

func (ss stringSlice) RedisArg() interface{} {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)

	//#nosec
	enc.Encode(ss)

	return buf.Bytes()
}

func (ss *stringSlice) RedisScan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("string slice: cannot convert from %T to %T", src, ss)
	}

	var buf bytes.Buffer
	if _, err := buf.Write(v); err != nil {
		return fmt.Errorf("string slice: write failed: %w", err)
	}

	dec := json.NewDecoder(&buf)
	if err := dec.Decode(ss); err != nil {
		return fmt.Errorf("string slice: decode failed: %w", err)
	}

	return nil
}

type User struct {
	Id          string      `mapstructure:"id"`
	Password    string      `mapstructure:"password"`
	Seed        string      `mapstructure:"seed"`
	Finalised   bool        `mapstructure:"finalised"`
	Tokens      stringSlice `mapstructure:"tokens"`
	Connections stringSlice `mapstructure:"connections"`
	PublicKeys  stringSlice `mapstructure:"public_keys"`
}

type Group struct {
	Id          string
	Owners      stringSlice
	Members     stringSlice //owner exists in both owners and members; we use members to send messages
	IsOpen      bool
	IsBroadcast bool
}

func NewRedis(addr string) (r Redis, err error) {
	r.pool = &redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 30 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}

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

	c := r.pool.Get()
	defer c.Close()

	_, err = c.Do("HSET", tokenIDHMKey, token, id)

	return
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
	c := r.pool.Get()
	defer c.Close()

	_, err = c.Do("HSET", id, "Finalised", true)

	return
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

func (r Redis) UserIdByToken(token string) (id string, err error) {
	c := r.pool.Get()
	defer c.Close()

	res, err := c.Do("HGET", tokenIDHMKey, token)
	if err != nil || res == nil {
		return
	}

	return string(res.([]byte)), nil
}

func (r Redis) UserByToken(token string) (u User, err error) {
	id, err := r.UserIdByToken(token)
	if err != nil {
		return
	}

	return r.loadUser(id)
}

func (r Redis) DeleteToken(token string) (err error) {
	c := r.pool.Get()
	defer c.Close()

	_, err = c.Do("HDEL", tokenIDHMKey, token)

	return
}

func (r Redis) Messages(id string) chan []byte {
	c := r.pool.Get()

	out := make(chan []byte)

	go func(c redis.Conn, id string, o chan []byte) {
		defer c.Close()
		defer close(out)

		psc := redis.PubSubConn{Conn: c}

		// #nosec
		psc.Subscribe(id)

		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				o <- v.Data
			case redis.Subscription:
				continue
			case error:
				log.Print(v)

				return
			}
		}
	}(c, id, out)

	return out
}

func (r Redis) WriteMessage(recipient string, payload []byte) (err error) {
	c := r.pool.Get()
	defer c.Close()

	_, err = c.Do("PUBLISH", recipient, payload)

	return
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

func (r Redis) loadUser(id string) (u User, err error) {
	err = r.load(id, &u)

	return
}

func (r Redis) saveUser(u User) (err error) {
	return r.save(u.Id, u)
}

func (r Redis) loadGroup(id string) (g Group, err error) {
	err = r.load(id, &g)

	return
}

func (r Redis) saveGroup(g Group) (err error) {
	return r.save(g.Id, g)
}

func (r Redis) save(id string, o interface{}) (err error) {
	c := r.pool.Get()
	defer c.Close()

	_, err = c.Do("HSET", redis.Args{}.Add(id).AddFlat(o)...)

	return
}

func (r Redis) load(id string, o interface{}) (err error) {
	c := r.pool.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", id))
	if err != nil {
		return
	}

	return redis.ScanStruct(values, o)
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

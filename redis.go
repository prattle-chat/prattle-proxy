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

type User struct {
	Id          string   `mapstructure:"id"`
	Password    string   `mapstructure:"password"`
	Seed        string   `mapstructure:"seed"`
	Finalised   bool     `mapstructure:"finalised"`
	Tokens      []string `mapstructure:"tokens"`
	Connections []string `mapstructure:"connections"`
	PublicKeys  []string `mapstructure:"public_keys"`
}

type Redis struct {
	c *redis.Client
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

func (r Redis) loadUser(id string) (u User, err error) {
	res, err := r.c.Get(context.Background(), id).Result()
	if err != nil {
		return
	}

	b := bytes.NewBufferString(res)

	dec := gob.NewDecoder(b)
	err = dec.Decode(&u)

	return
}

func (r Redis) saveUser(u User) (err error) {
	b := bytes.Buffer{}

	enc := gob.NewEncoder(&b)
	err = enc.Encode(u)
	if err != nil {
		return
	}

	return r.c.Set(context.Background(), u.Id, b.String(), 0).Err()
}

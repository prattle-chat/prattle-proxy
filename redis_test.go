package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
)

var conn = redigomock.NewConn()

func noRedisCall(conn *redigomock.Conn) {
	conn.Clear()
}

func redisErrorOnAuth(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "foo").
		ExpectError(fmt.Errorf("an error"))
}

func validTokenInvalidUser(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	conn.Command("HGETALL", "some-user@testing").
		ExpectError(fmt.Errorf("an error"))

	conn.Command("HDEL", tokenIDHMKey, "foo").
		Expect("ok")
}

func validTokenNonFinaliseduser(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	conn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "0",
			"PublicKeys": "[]\n",
		})
}

func validTokenAndUser(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	conn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": "[]\n",
		})
}

func validTokenAndUserWithKeys(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	conn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": `["key-1", "key-2", "key-3", "key-4", "key-5"]`,
		})
}

func validTokenUserRecipientAndPublish(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	conn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": "[]\n",
		})

	conn.Command("PUBLISH", "recipient@testing", []uint8{0x41, 0xff, 0x81, 0x3, 0x1, 0x1, 0xe, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x1, 0xc, 0x0, 0x1, 0x6, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xc, 0x0, 0x1, 0x7, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x1, 0xa, 0x0, 0x0, 0x0, 0x29, 0xff, 0x82, 0x1, 0x11, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x40, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1, 0x11, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x40, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x0})
}

func validPeeredUser(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "blahblahblah").
		Expect([]byte(""))

	conn.Command("PUBLISH", "recipient@testing", []uint8{0x41, 0xff, 0x81, 0x3, 0x1, 0x1, 0xe, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x1, 0xc, 0x0, 0x1, 0x6, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xc, 0x0, 0x1, 0x7, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x1, 0xa, 0x0, 0x0, 0x0, 0x26, 0xff, 0x82, 0x1, 0x11, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x40, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1, 0xe, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x40, 0x6e, 0x6f, 0x6e, 0x65, 0x0})
}

func invalidPeeredUser(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "foobarbaz").
		Expect([]byte(""))
}

func validPeeredToPeered(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "blahblahblah").
		Expect([]byte(""))
}

func createSubscriptionMessage(channel string, data []byte) []interface{} {
	values := []interface{}{}
	values = append(values, interface{}([]byte("message")))
	values = append(values, interface{}([]byte(channel)))
	values = append(values, interface{}(data))
	return values
}

func validPeerReceiveMessage(conn *redigomock.Conn) {
	values := []interface{}{}
	values = append(values, interface{}([]byte("subscribe")))
	values = append(values, interface{}([]byte("some-user@testing")))
	values = append(values, interface{}([]byte("1")))

	conn.Clear()
	conn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	conn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": `["key-1", "key-2", "key-3", "key-4", "key-5"]`,
		})

	conn.Command("SUBSCRIBE", "some-user@testing").Expect(values)
	conn.ReceiveWait = true

	conn.AddSubscriptionMessage(createSubscriptionMessage("some-user@testing", []uint8{0x41, 0xff, 0x81, 0x3, 0x1, 0x1, 0xe, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x1, 0xc, 0x0, 0x1, 0x6, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xc, 0x0, 0x1, 0x7, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x1, 0xa, 0x0, 0x0, 0x0, 0x26, 0xff, 0x82, 0x1, 0x11, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x40, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1, 0xe, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x40, 0x6e, 0x6f, 0x6e, 0x65, 0x0}))

	// You need to send messages to conn.ReceiveNow in order to get a response.
	// Sending to this channel will block until receive, so do it in a goroutine
	go func() {
		conn.ReceiveNow <- true // This unlocks the subscribe message
		conn.ReceiveNow <- true // This sends the "hello" message
	}()
}

func validPasswordCreateUser(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGETALL", "some-user@testing").
		Expect([]byte(""))

	conn.GenericCommand("HSET").
		Expect([]byte("ok"))
}

func idFound(conn *redigomock.Conn) {
	conn.Clear()

	for i := 0; i < 10; i++ {
		conn.Command("HGETALL", "some-user@testing").
			ExpectMap(map[string]string{
				"Id": "some-user@testing",
			})
	}
}

func userCreateError(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGETALL", "some-user@testing").
		Expect([]byte(""))

	conn.GenericCommand("HSET").
		ExpectError(fmt.Errorf("an error"))
}

func noSuchUser(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGETALL", "some-user@testing").
		Expect([]byte(""))
}

func validUser(conn *redigomock.Conn) {
	conn.Clear()
	conn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": `["key-1", "key-2", "key-3", "key-4", "key-5"]`,
			"Password":   "$argon2id$v=19$m=65536,t=1,p=2$yQg5GtsbXHNKF1tUyeNtAw$DudFezvG5IeGTy/ILi179Fhjm29K7ECbJahb0al2BbA", // 16 'a' chars
			"Seed":       "PAZXXJTBLKTDLHWTDPEGGMYNVR2LMJGC",
		})
}

func setFinalised(conn *redigomock.Conn) {
	validUser(conn)
	conn.Command("HSET", "some-user@testing", "Finalised", true)
}

func addToken(conn *redigomock.Conn) {
	validUser(conn)
	conn.GenericCommand("HSET").
		Expect([]byte("ok"))
}

func NewDummyRedis(c *redigomock.Conn) (r Redis) {
	return Redis{
		pool: &redis.Pool{
			MaxIdle: 10,
			Dial: func() (redis.Conn, error) {
				return c, nil
			},
		},
	}
}

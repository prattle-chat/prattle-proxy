package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

var redigoMockConn = redigomock.NewConn()

func noRedisCall() {
	redigoMockConn.Clear()
}

func redisErrorOnAuth() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "foo").
		ExpectError(fmt.Errorf("an error"))
}

func validTokenInvalidUser() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	redigoMockConn.Command("HGETALL", "some-user@testing").
		ExpectError(fmt.Errorf("an error"))

	redigoMockConn.Command("HDEL", tokenIDHMKey, "foo").
		Expect("ok")
}

func validTokenNonFinaliseduser() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	redigoMockConn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "0",
			"PublicKeys": "[]\n",
		})
}

func validTokenAndUser() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	redigoMockConn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": "[]\n",
		})
}

func validTokenAndUserWithKeys() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	redigoMockConn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": `["key-1", "key-2", "key-3", "key-4", "key-5"]`,
		})
}

func validTokenUserRecipientAndPublish() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	redigoMockConn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": "[]\n",
		})

	redigoMockConn.Command("PUBLISH", "recipient@testing", []uint8{0x41, 0xff, 0x81, 0x3, 0x1, 0x1, 0xe, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x1, 0xc, 0x0, 0x1, 0x6, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xc, 0x0, 0x1, 0x7, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x1, 0xa, 0x0, 0x0, 0x0, 0x29, 0xff, 0x82, 0x1, 0x11, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x40, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1, 0x11, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x40, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x0})
}

func validPeeredUser() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "blahblahblah").
		Expect([]byte(""))

	redigoMockConn.Command("PUBLISH", "recipient@testing", []uint8{0x41, 0xff, 0x81, 0x3, 0x1, 0x1, 0xe, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x1, 0xc, 0x0, 0x1, 0x6, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xc, 0x0, 0x1, 0x7, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x1, 0xa, 0x0, 0x0, 0x0, 0x26, 0xff, 0x82, 0x1, 0x11, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x40, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1, 0xe, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x40, 0x6e, 0x6f, 0x6e, 0x65, 0x0})
}

func invalidPeeredUser() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "foobarbaz").
		Expect([]byte(""))
}

func validPeeredToPeered() {
	redigoMockConn.Clear()
	redigoMockConn.Command("HGET", tokenIDHMKey, "blahblahblah").
		Expect([]byte(""))
}

func createSubscriptionMessage(channel string, data []byte) []interface{} {
	values := []interface{}{}
	values = append(values, interface{}([]byte("message")))
	values = append(values, interface{}([]byte(channel)))
	values = append(values, interface{}(data))
	return values
}

func validPeerReceiveMessage() {
	redigoMockConn.Clear()

	values := []interface{}{}
	values = append(values, interface{}([]byte("subscribe")))
	values = append(values, interface{}([]byte("some-user@testing")))
	values = append(values, interface{}([]byte("1")))

	redigoMockConn.Command("HGET", tokenIDHMKey, "foo").
		Expect([]byte("some-user@testing"))

	redigoMockConn.Command("HGETALL", "some-user@testing").
		ExpectMap(map[string]string{
			"Id":         "some-user@testing",
			"Finalised":  "1",
			"PublicKeys": `["key-1", "key-2", "key-3", "key-4", "key-5"]`,
		})

	redigoMockConn.Command("SUBSCRIBE", "some-user@testing").Expect(values)
	redigoMockConn.ReceiveWait = true

	redigoMockConn.AddSubscriptionMessage(createSubscriptionMessage("some-user@testing", []uint8{0x41, 0xff, 0x81, 0x3, 0x1, 0x1, 0xe, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x1, 0xc, 0x0, 0x1, 0x6, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xc, 0x0, 0x1, 0x7, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x1, 0xa, 0x0, 0x0, 0x0, 0x26, 0xff, 0x82, 0x1, 0x11, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x40, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1, 0xe, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x40, 0x6e, 0x6f, 0x6e, 0x65, 0x0}))

	// You need to send messages to conn.ReceiveNow in order to get a response.
	// Sending to this channel will block until receive, so do it in a goroutine
	go func() {
		redigoMockConn.ReceiveNow <- true // This unlocks the subscribe message
		redigoMockConn.ReceiveNow <- true // This sends the "hello" message
	}()
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
# The Prattle Protocol

This document describes the prattle protocol, including:

1. Authentication flows
1. Messaging
1. Groups and Group Operations
1. Federation with other Prattle protocols

This document uses the following meanings:

* **Prattle, Server(s), Implementation(s)** refer to servers which provide messaging, groups, and federated access to other implementations
* **Home Server(s)** refers to the server on which a client user was created, and as such has ultimate authority over whether a user is valid
* **Federation, Federated Peer, Peer(s)** refer to trusted connections between two peers, which use special bearer tokens to validate one another
* **Relay** is a verb used to connote a peer accepting a message meant for a remote server, which it then forwards on *to* that remote server
* **Reference Implementation** refers to the implementation developed at [github.com/prattle-chat/prattle-proxy](github.com/prattle-chat/prattle-proxy)
* **Client(s)** refer to both users of an implementation, and the applications they use to perform operations against a server
* **Must/ Must not** describe a case with strict adherence; prattle implementations must do these things exactly as described
* **May/ May not** describe a case with loose adherence; implementations must set sensible defaults where such cases are skipped/ missing

## Protobuf files

Prattle uses protobuf definitions to describe the messages and services which clients, servers, and peers use to communicate. These files may be found at [/protos](/protos), with the reference implementation living at the root of this project.

Where applicable, this document links directly to messages and services to illustrate how things work.

## Federation

When relaying from a user's home server to a peer, the home server must first validate the token points to a valid user, then add the following header to relayed/ federated requests:

```
operator_id: some-user@example.com
```

Where `some-user@example.com` is the ID of the owner of the calling user. Peers must use the value of this header to validate payloads, such as group permissions, and that the sender of a message has not been spoofed.

A server must fail an request from a peer where this header does not exist.

## Signup and Authentication

Prattle implementations must use a mixture of generated IDs, strong passwords, and TOTP tokens to signup users. Implementations may set maximum password lengths.

Client signup looks like:

1. A client sends a password to a server in a [SignupRequest](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/auth.proto#L26-L34) message to [auth.Signup](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/auth.proto#L12-L13)
1. A server mints a TOTP seed, and an ID which it returns in a [SignupResponse](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/auth.proto#L35-L41) message
1. A client uses the TOTP seed to generate a TOTP value which it sends back, along with password and ID in an  [Auth](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/auth.proto#L43-L58) message to to [auth.Finalise](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/auth.proto#L15-L19)
1. A server returns an error or an [empty](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/empty.proto) message

A server must generate a sufficiently randomised, unique, ID for a client to use. Servers must not allow clients to set their own ID, to allow for better anonymisation and simpler flow. A server may use anything, from UUIDs to a [diceware implementation](https://en.wikipedia.org/wiki/Diceware). The reference implementation generates IDs similar to:

```go
minter MinterFunc = func(d string) (string, error) {
    words := diceware.MustGenerate(2)

    suffix := make([]byte, 4)
    _, err := rand.Read(suffix)
    if err != nil {
        return "", generalError
    }

    return fmt.Sprintf("%s-%s-%s@%s",
        words[0],
        words[1],
        hex.EncodeToString(suffix),
        d,
    ), err
}
```

Which generates IDs that look like:

* renegade-gaining-a787bc98@example.com
* outpour-obsessed-62f903de@example.com
* ashy-tidings-7658350c@example.com
* gladiator-polygraph-d17cec0b@example.com
* slab-cassette-f20e2cd8@example.com

In the above cases, the domain name is `example.com` and is used to route messages correctly.

## Token generation

Clients must use bearer tokens for further operations. A server must accept an [Auth](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/auth.proto#L43-L58) message to initiate a session, which it uses to mint an ID and return to a client in a [TokenValue](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/auth.proto#L60-L64) message, via the [auth.Token](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/auth.proto#L21-L23) endpoint.

A server must generate a cryptographically random token and must not use anything which includes de-anonymising data, such as JWT.

The reference implementation generates bearer tokens as per:

```go
func (s Server) mintToken() string {
    return fmt.Sprintf("prattle-%s%s%s",
        hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()),
        hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()),
        hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()),
    )
}
```

Which generates tokens that look like:

* prattle-495f37b7f56d4cc2b8193ae734254a7658c317fcab114d1fbc63e5ff4f3f4a388c9bb421b3054e5ba39d0f55ab1d3008
* prattle-8e2557cf4e334d0280bec4ce9990aa0ae8d80e706ef74532b170c7d35cbf2ffb36115318901f43edad9abaeb78b32aa4
* prattle-04c73b8f4d164966bf9612c72c1d3d29fa93a61fb5994d1d9dedf540e887ef61ade01cae0c0a4f648a2ae9131c252d00
* prattle-fcf6cdbaa248465792d648b70bab11570518a370a3f74d8e97c32ae9c1eb5830f2c89e8990214fc089ce7fea6ea610af
* prattle-24cd8caf976a4932b2496b7afb27353c2bf17e0692d949f3bf1b7605bc88624b879f85fb277c4546a48ee320183f13df

From hereonin, this token is used to authenticate a user and should be sent as the header value:

```
authorization: bearer prattle-24cd8caf976a4932b2496b7afb27353c2bf17e0692d949f3bf1b7605bc88624b879f85fb277c4546a48ee320183f13df
```

(for example).

## Public key exchange

Prattle clients must use [age](https://docs.google.com/document/d/11yHom20CrsuX8KQJXBBw04s80Unjv8zCg_A7sPAX_9Y/preview#) keys to encrypt/ decrypt messages. A client must upload public keys to their home server by enclosing a public key in a [PublicKeyValue](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/user.proto#L60-L64) and sending it to [user.PublicKey](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/user.proto#L48-L50).

A client may safely ignore the value of [PublicKeyValue.user_id](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/user.proto#L63) by leaving it unset. A server must not use this value to determine the owner of a key, but must instead test the bearer token of the user.PublicKey request.

A server must not relay these messages; a client must connect to their home server to upload their public key(s).

## Messaging

A [Message](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/messaging.proto#L57-L69) is used to store a message, either when sent as a Direct Message (when a message is sent from one person to one person) or a Group message (when a message is sent from one person to many people).

When a Message represents a Direct Message, the [recipient field](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/messaging.proto#L64) must be the `user_id` associated with that user. When a Message represents a Group Message, the recipient field must be the `group_id` of that group.

A Message must be encrypted and added to the 'encoded' field of a [MessageWrapper](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/messaging.proto#L28-L50). For Direct Messages, a client may retrieve the recipient's keys each time, or it may cache keys for a time. For Group Messages, a client may retrieve the keys for every member of the group, or it make cache keys for a time.

The Recipient field of a MessageWrapper must contain the `id` of the user for whom this message should be delivered. This `id` must match the owner of the PublicKey used to encrypt the message (since only that user has the private key necessary to decrypt a message).

When a MessageWrapper wraps a Group Message, the client must set the `group_id` field of the Recipient to the ID of the Group the message is being sent to. The client must match this `group_id` with the `recipient` field of the Message, discarding any message where these do not match.

Because messages to large groups (especially where users in groups have multiple PrivateKeys) result in a deluge of messages being sent by a Client, a Prattle server may limit the number of members a group has, and the number of keys a user has.

MessageWrappers are sent to a user's Home Server on [messaging.Send](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/messaging.proto#L16-L17). A server must relay Direct Messages to only to peered servers, and only to the peer whose domain matches the domain of the recipient. For Group Messages, the home server must relay via the peer whose domain matches the `group_id`. A server may locally deliver to group recipients if the domain of the recipient's ID matches the domain for the home server.

A client must connect to its home server and call [messaging.Subscribe](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/messaging.proto#L19-L25) to receive a stream of `MessageWrapper`s meant for this user.

## Groups

A [Group](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L45-L55) can be either open (anyone can join) or closed (new members must be directly invited), and regular (any user can send messages) or broadcast (only owners can message).

A Group has two levels of membership; [Owner](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L51), and [Member](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L52). Owners are also Members.

Users in a Group can be [invited](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L32), [promoted](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L35), [demoted](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L39), and can [manually leave](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L42).

A Group is created by sending a [Group](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L49-L55) message to [groups.Create](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L22). A client must set both [`is_open`](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L53) and [`is_broadcast`](https://github.com/prattle-chat/prattle-proxy/blob/main/protos/group.proto#L54). Anything else is ignored.

A Group ID is minted in the same way as a User ID, with the addition of the prefix `g:`.

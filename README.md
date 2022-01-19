# Prattle

Prattle is a chat protocol designed around privacy, security, and trust. The aims of the project are:

1. To provide a limited link between Prattle Users and Real People
1. To provide privacy through strong end-to-end secrecy
1. To provide message secrecy through strict limits in how long encrypted messages stay on servers, at the risk of dropping messages
1. To provide strong anti-harrassment and anti-spam controls through opaque blocklists stored serverside

The protocol is described in [PROTOCOL.md](PROTOCOL.md), with the [protos/](/protos/) directory containing the raw proto files used to generate servers and clients.

## Prattle Proxy

This project contains the reference implementation of the Prattle Protocol.

It:

1. Uses redis to store user, group, and token information
1. Uses redis pub/sub to route messages (deliberately no DLQ; once a message is received, it's sent- any offline client will miss out on the message. This is a very strict interpretation of goal 3)
1. Hashes passwords using argon2id ^[1](https://en.wikipedia.org/wiki/Argon2) ^[2](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html) ^[3](https://cryptobook.nakov.com/mac-and-key-derivation/argon2) ^[4](https://tools.ietf.org/id/draft-irtf-cfrg-argon2-03.html) to minimise the chances of exfiltrated passwords being used to upload new keys and initiate new client connections
1. Requires no storage, requires no root/ setuid permissions, is statically compiled, and uses a number of binary hardening techniques in the expectation that the service can be run with reasonably secure defaults.

### Building

This project uses the default go toolchain, wrapped in `make` to generate certifcates and hardened binaries:

```bash
$ make
mkdir -p certs
(cd certs && ../scripts/gen-cert)
/home/user/src/prattle-proxy/certs
+ openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj '/C=GB/ST=GB/L=Scunthorpe/O=Prattle/OU=Engineering/CN=*/emailAddress=engineering@prattle'

### SNIP

CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o prattle-proxy
```

This project contains tests. Tests are written using the default testing package, and any pull request to change that will be rejected.

```bash
$ go test
PASS
ok      github.com/prattle-chat/prattle-proxy   0.444s
```

### Running

This project can be configured either using the file `proxy.toml` (or `proxy.yaml`/ `proxy.json`/ etc., but the only tested configuration is `toml`), or through the environment.

To get started, your best bet is to muck around with the included `proxy.toml` file. Environment variables are just capitalised versions of the keys in that file.

Once built, this project can be run as per:

```bash
$ ./prattle-proxy
2022/01/19 11:06:16 Starting server on localhost:8080
```

#### Configuring Federated Peers

Peers can be configured in `proxy.toml` like:

```toml
[federations]
[federations."example.com"]
connection_string = "prattle.example.com:443"
psk = "d9cf51e01a43d92b2f1f626f852323bbafe410c94f8a3ec622f19ed40089cd40"
```

This configures a peering connection with `example.com`, connecting on `prattle.example.com:443`. All connections are TLS by default.

This provides two things:

1. When a local user tries to send a message to a user or group on `example.com`, that request is relayed to `prattle.example.com:443` setting `d9cf51e01a43d92b2f1f626f852323bbafe410c94f8a3ec622f19ed40089cd40` as the bearer token (rather than the bearer token the user uses to authenticate to the local server
1. When a request comes in using that psk as the bearer token, the local server semi-implicitly trusts that the peer has validated the sender of a message. Instead, the local server matches up the sender domain against the stored domain for this psk (to limit spoofing) and then compares against local membership/ blocklists

For a peering connection to work, it must be configured on both ends using the same pre-shared key.


## Licence

Copyright (c) 2022 prattle-chat

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

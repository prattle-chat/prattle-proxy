# Prattle Proxy

This proxy:

1. Exposes a series of gRPC endpoints to facilitate exchanging encrypted messages between users
1. Exposes an enpoint to accept/ forward messages to other prattle proxies

This proxy doesn't:

1. Have any way of decrypting messages sent between parties
1. Store anything directly on disk

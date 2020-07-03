# PubSub

Simple Publisher/Subscriber based on WebSockets

### Build
```bash
make
```

### Server
```sh
server -h
usage: server [<flags>]

PubSub server

Flags:
  -h, --help       Show context-sensitive help (also try --help-long and --help-man).
  -p, --port=8080  Listening port.
```

### Subscriber
```sh
subscriber -h
usage: subscriber [<flags>]

Subscriber

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
  -n, --count=1  Number of subscribers.
  -a, --address="localhost:8080" Server address.
```

### Publisher
```sh
publisher -h
usage: publisher [<flags>]

Publisher

Flags:
  -h, --help             Show context-sensitive help (also try --help-long and --help-man).
  -a, --address="localhost:8080" Server address.
  -m, --message=MESSAGE  Message to publish.
  ```
  
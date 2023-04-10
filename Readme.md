# xk6-ssh
A k6 extension for decoding and encoding proto files. Built using [xk6](https://github.com/grafana/xk6).

## Build

To build this extension locally follow the next steps:

- Install [Go toolchain](https://go101.org/article/go-toolchain.html)
- Install [K6](https://k6.io/docs/get-started/installation/)
- Install **xk6**: ```go install github.com/k6io/xk6/cmd/xk6@latest```
- Build extension: ```xk6 build --with xk6-protobuf=.```

## Example

Command to run examples:
```./k6 run examples/codec.js```
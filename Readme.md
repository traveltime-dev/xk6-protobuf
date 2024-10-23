# xk6-protobuf
A k6 extension for decoding and encoding proto files. Built using [xk6](https://github.com/grafana/xk6).


## Setup and Running 

Follow these steps to set up and run the K6 benchmark with Protobuf support:

### 1. Install go

Installing [Go toolchain](https://go101.org/article/go-toolchain.html)

### 2. Export the Go PATH

Add the Go binary directory to your system's `PATH` environment variable:

```bash
export PATH=$(go env GOPATH)/bin:$PATH
```

### 2. Install xk6

Install the latest version of xk6:

```bash
go install go.k6.io/xk6/cmd/xk6@latest
```

### 3. Build

Build the K6 binary with the `xk6-protobuf` extension:

```bash
xk6 build --with github.com/traveltime-dev/xk6-protobuf@latest
```

If you want to build it locally, then clone this repository and use:
```bash
xk6 build --with xk6-protobuf=.
```

### 4. Run

Run the K6 benchmark using the generated binary and the protobuf benchmark file:

```bash
./k6 run {proto-benchmark-file}.js
```

Replace `{proto-benchmark-file}.js` with the actual file name of your protobuf benchmark script.

## Examples

Command to run examples:
```./k6 run example/codec.js```

# bpflog

bpflog is a tiny Go library and example that logs messages from an eBPF program to userspace using a perf ring buffer. The example demonstrates attaching an XDP program that records IPv4 source/destination pairs and emits human-readable log lines which the Go program consumes and prints.

## This repository contains

- `handler.go` - small library that wraps a `cilium/ebpf` perf reader and exposes Start/Stop semantics.
- `example/` - example XDP program (C), generated Go bindings and a Go `main` that attaches the XDP program and consumes logs.
- `bpf.c`, `vmlinux.h`, `logf.h` - the BPF program and helpers used by the example.

## Requirements

- Linux with kernel headers and BPF support (XDP).
- Go 1.25+ (module configured in `go.mod`).
- clang/llvm and `bpftool`/`iproute2` if you build or load the BPF program manually.
- Running the example requires CAP_NET_ADMIN / root to attach the XDP program.

## Build and run the example (quick)

1. Install Go (1.25+), clang/llvm, and [mise](https://mise.jdx.dev).
2. From the repo root run the example build and run commands.

```bash
# generate vmlinux headers
mise gen-vmlinux

# run the example
cd example
mise start
```

## Notes and tips

- The example hardcodes an interface index (2) in `example/main.go`. Change it to your interface index (e.g., `ip link` to list). Alternatively modify the code to accept an interface name and resolve its index.
- The BPF program filters out 192.168.0.0/16 in the example; adjust as needed in `example/bpf.c`.
- If `go generate` fails because `bpf2go` is not available, install it from the `cilium/ebpf` cmd package:

```bash
GO111MODULE=on go install github.com/cilium/ebpf/cmd/bpf2go@latest
```

- The library `NewHandler` in `handler.go` exposes a Start/Stop loop which reads `perf.Record` entries and calls your handler callback. The default read deadline is 100ms to allow graceful stop.

## License

This project follows the license declared in the repository (check LICENSE if present). The example BPF program declares GPL.

## Contributing

Open issues or PRs if you find bugs or want to improve the example (accepting improvements such as interface selection by name, better build scripts, or CI for kernel headers).

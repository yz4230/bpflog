package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/yz4230/bpflog"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target amd64 bpf ./bpf.c

func main() {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	var objs bpfObjects
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("failed to load objects: %v", err)
	}
	defer objs.Close()

	l, err := link.AttachXDP(link.XDPOptions{
		Interface: 2,
		Program:   objs.DoXdp,
	})
	if err != nil {
		log.Fatalf("failed to attach xdp: %v", err)
	}
	defer l.Close()

	logh := bpflog.NewHandler(objs.LogfEntries, func(r *perf.Record) {
		msg := string(r.RawSample)
		idx := strings.LastIndex(msg, "\n")
		log.Printf("log: %s", msg[:idx])
	})

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		if err := logh.Start(); err != nil {
			log.Fatalf("failed to start log handler: %v", err)
		}
	})

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	<-chSignal
	logh.Stop()
	wg.Wait()
}

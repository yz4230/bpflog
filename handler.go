package bpflog

import (
	"errors"
	"os"
	"sync/atomic"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"
)

var Deadline = 100 * time.Millisecond

type Handler struct {
	m       *ebpf.Map
	running atomic.Bool
	handler func(*perf.Record)
}

func NewHandler(m *ebpf.Map, f func(*perf.Record)) *Handler {
	return &Handler{m: m, handler: f}
}

func (h *Handler) Start() error {
	h.running.Store(true)

	r, err := perf.NewReader(h.m, os.Getpagesize())
	if err != nil {
		return err
	}
	defer r.Close()

	for h.running.Load() {
		r.SetDeadline(time.Now().Add(Deadline))
		event, err := r.Read()
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			}
			return err
		}
		h.handler(&event)
	}

	return nil
}

func (h *Handler) Stop() {
	h.running.Store(false)
}

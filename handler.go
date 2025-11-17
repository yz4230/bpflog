package bpflog

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"
)

var Deadline = 100 * time.Millisecond

type Handler struct {
	m       *ebpf.Map
	closer  io.Closer
	handler func(*perf.Record)
}

func NewHandler(m *ebpf.Map, f func(*perf.Record)) *Handler {
	return &Handler{m: m, handler: f}
}

func (h *Handler) Start() error {
	r, err := perf.NewReader(h.m, os.Getpagesize())
	if err != nil {
		return err
	}
	h.closer = r

	for {
		event, err := r.Read()
		if err != nil {
			if errors.Is(err, os.ErrClosed) {
				break
			}
			return err
		}
		h.handler(&event)
	}

	return nil
}

func (h *Handler) Stop() error {
	if h.closer != nil {
		return h.closer.Close()
	}
	return nil
}

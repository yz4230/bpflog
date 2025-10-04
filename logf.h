// clang-format off
#include "vmlinux.h" // IWYU pragma: keep
#include <bpf/bpf_helpers.h>
// clang-format on

#ifndef LOGF_H
#define LOGF_H

struct {
  __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
  __uint(max_entries, 128);
  __uint(key_size, sizeof(__u32));
  __uint(value_size, sizeof(__u32));
} logf_entries SEC(".maps");

// logf outputs logs to user space
// Note1: Arguments must be cast to u64 type
// Note2: Maximum log length is 256 bytes
#define logf(ctx, fmt, args...)                                                    \
  ({                                                                           \
    static const char _fmt[] = fmt;                                            \
    static char _buf[256];                                                     \
    u64 _args[___bpf_narg(args)];                                              \
    ___bpf_fill(_args, args);                                                  \
    int _len = bpf_snprintf(_buf, sizeof(_buf), _fmt, _args, sizeof(_args));   \
    if (_len < sizeof(_buf)) {                                                 \
      bpf_perf_event_output(ctx, &logf_entries, BPF_F_CURRENT_CPU, _buf, _len); \
    }                                                                          \
  })

#endif  // LOGF_H

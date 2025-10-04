//go:build ignore

#include "../vmlinux.h"
#include "../logf.h"

// clang-format off
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
// clang-format on

#define ETH_P_IP 0x0800

SEC("xdp")
int do_xdp(struct xdp_md* ctx)
{
    void *data = (void*)(long)ctx->data;
    void *data_end = (void*)(long)ctx->data_end;

    struct ethhdr *eth = data;
    if ((void*)(eth + 1) > data_end) return XDP_PASS;
    if (bpf_ntohs(eth->h_proto) != ETH_P_IP) return XDP_PASS;

    struct iphdr *ip = data + sizeof(*eth);
    if ((void*)(ip + 1) > data_end) return XDP_PASS;

    // if ipv4, log src and dst ip
    if (ip->version == 4) {
        // ignore 192.168.0.0/16
        if ((bpf_ntohl(ip->saddr) & 0xFFFF0000) == 0xC0A80000) { 
            return XDP_PASS;
        }
        logf(ctx, "%pI4 -> %pI4\n", (u64)&ip->saddr, (u64)&ip->daddr);
    }

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ipsets-io/ipsets"
	"github.com/ipsets-io/ipsets/internal/build"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	index, err := build.Run(ctx, ipsets.Providers(), build.Options{Dir: "docs"})

	for _, p := range index.Providers {
		for _, s := range p.Sets {
			fmt.Printf("%s/%s ipv4=%d ipv6=%d\n", p.ID, s.ID, s.IPv4.Count, s.IPv6.Count)
		}
	}
	for _, c := range index.Categories {
		for _, s := range c.Sets {
			fmt.Printf("categories/%s ipv4=%d ipv6=%d\n", c.ID, s.IPv4.Count, s.IPv6.Count)
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/masumomo/go_todo_app/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("Start with aaa %v", url)
	mux, clearnup, err := NewMux(ctx, cfg)
	if err != nil {
		return err
	}
	defer clearnup()
	s := NewServer(l, mux)

	return s.Run(ctx)
}

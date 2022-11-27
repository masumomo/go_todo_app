package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("need port number: %d\n", len(os.Args))
		os.Exit(1)
	}

	log.Printf("port number: %s\n", os.Args[1])
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Serve(l); err != nil {
			if err != http.ErrServerClosed {
				log.Printf("failed to close: %+v", err)
				return err
			}
		}

		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	return eg.Wait()
}

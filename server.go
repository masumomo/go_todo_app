package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	// cfg, err := config.New()
	// if err != nil {
	// 	return err
	// }

	// l, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", cfg.Port))
	// if err != nil {
	// 	log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	// }

	// url := fmt.Sprintf("http://%s", l.Addr().String())
	// fmt.Printf("Start with (Graceful mode) %v", url)
	// // fmt.Printf("Start with (Non graceful mode) %v", url)
	// s := &http.Server{
	// 	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		time.Sleep(5 * time.Second)
	// 		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	// 	}),
	// }
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil {
			if err != http.ErrServerClosed {
				log.Printf("failed to close: %+v", err)
				return err
			}
		}
		return nil
	})

	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	return eg.Wait()
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tonytrg/backend/internal/api"
)

func main() {
	addr := net.JoinHostPort("localhost", "8080")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go listenSignal(cancel)

	err := runServer(ctx, addr)
	if err != nil {
		log.Printf("server error: %s\n", err)
	} else {
		log.Println("application shutdown complete")
	}
}

func runServer(ctx context.Context, addr string) error {
	server := &http.Server{
		Addr:              addr,
		Handler:           api.ApiHandler(),
		ReadTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 5,
	}

	var err error
	go func() {
		if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed to start: %s\n", err)
		}
	}()
	log.Printf("server started successfully on addr: %s\n", addr)

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = server.Shutdown(ctxShutDown)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func listenSignal(cancel context.CancelFunc) {
	var captureSignal = make(chan os.Signal, 1)
	signal.Notify(captureSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	switch <-captureSignal {
	case syscall.SIGINT:
		cancel()
	case syscall.SIGTERM:
		log.Fatal("application force stop")
	default:
		fmt.Println("- unknown signal")
	}
}

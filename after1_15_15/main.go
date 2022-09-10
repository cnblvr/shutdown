package main

import (
	"context"
	"github.com/cnblvr/shutdown/sleep"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const gracefulTimeout = time.Second * 10

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	http.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		duration, err := time.ParseDuration(r.URL.Query().Get("t"))
		if err != nil {
			duration = time.Second * 5
		}
		sleep.Do(duration)
	})

	srv := http.Server{Addr: ":8080"}
	go func() {
		if err := srv.ListenAndServe(); err == http.ErrServerClosed {
			log.Printf("server is shutting down... Graceful shutdown timeout is %s", gracefulTimeout)
		} else {
			log.Printf("listen and server error: %v", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	switch err := srv.Shutdown(ctx); err {
	case nil:
		log.Printf("server is down")
	case context.DeadlineExceeded:
		log.Printf("server shut down unsuccessfully")
		return
	default:
		log.Printf("server shutdown error: %v", err)
	}
}

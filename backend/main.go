package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	gcInterval = flag.Duration("interval", 1*time.Hour, "garbage collection interval")
	expire     = flag.Duration("expire", 7*24*time.Hour, "expire duration")
	cacheDir   = flag.String("cache", "./cache", "cache directory")
	listen     = flag.String("listen", "./listen.socket", "listen unix domain socket")
	debug      = flag.Bool("debug", false, "debug flag")
	jobTimeout = flag.Duration("timeout", 5*time.Second, "job timeout")
	queueLen   = flag.Int("queue", 10, "job queue length")
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
	log.Print("Server finished")
}

func execute() error {
	flag.VisitAll(func(f *flag.Flag) {
		if v, ok := os.LookupEnv(strings.ToUpper(f.Name)); ok {
			f.Value.Set(v)
		}
	})
	flag.Parse()
	jobCh := make(chan *Job, *queueLen)
	defer close(jobCh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go GrabageCollect(ctx, *cacheDir, *gcInterval, *expire)
	go JobWorker(jobCh)
	http.Handle("/convert.cgi", &Handler{jobCh})
	return ListenAndServe()
}

func JobWorker(jobCh chan *Job) {
	for j := range jobCh {
		j.Finish <- j.Cmd.Run()
	}
}

func ListenAndServe() error {
	// listen unix domain socket
	l, err := net.Listen("unix", *listen)
	if err != nil {
		return fmt.Errorf("HTTP server Listen: %w", err)
	}
	defer l.Close()
	defer os.Remove(*listen)

	if err := os.Chmod(*listen, 0777); err != nil {
		return fmt.Errorf("HTTP socket chmod %w", err)
	}

	// start http server
	srv := http.Server{}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(l)
	}()

	// wait signal
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	// shutdown http server
	if err := srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("HTTP server Shutdown: %w", err)
	}

	// wait http server
	if err := <-errCh; err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server Serve: %w", err)
	}
	return nil
}

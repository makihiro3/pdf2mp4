package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

var (
	listen     = flag.String("listen", "./listen.socket", "listen unix domain socket")
	debug      = flag.Bool("debug", false, "debug flag")
	jobTimeout = flag.Duration("timeout", 5*time.Second, "job timeout")
)

type Handler struct {
}

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
	http.Handle("/convert.cgi", &Handler{})
	return ListenAndServe()
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Accpet POST method only
	if r.Method != http.MethodPost {
		m := http.StatusMethodNotAllowed
		w.WriteHeader(m)
		io.WriteString(w, http.StatusText(m))
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Print(err)
		m := http.StatusBadRequest
		w.WriteHeader(m)
		io.WriteString(w, http.StatusText(m))
		return
	}
	size := "0"
	switch v := r.Form.Get("size"); v {
	case "720":
		size = v
	case "1080":
		size = v
	case "1440":
		size = v
	case "2160":
		size = v
	case "original":
		size = "0"
	case "":
		size = "0"
	default:
		m := http.StatusBadRequest
		w.WriteHeader(m)
		io.WriteString(w, http.StatusText(m))
		return
	}

	interval := "1"
	switch v := r.Form.Get("interval"); v {
	case "1":
		interval = v
	case "2":
		interval = v
	case "3":
		interval = v
	default:
		m := http.StatusBadRequest
		w.WriteHeader(m)
		io.WriteString(w, http.StatusText(m))
		return
	}

	if err := h.Process(r.Context(), w, r.Body, size, interval); err != nil {
		log.Print(err)
		m := http.StatusInternalServerError
		w.WriteHeader(m)
		io.WriteString(w, http.StatusText(m))
		return
	}
	r.Body.Close()
}

func (h *Handler) Process(ctx context.Context, w io.Writer, r io.Reader, size, interval string) error {
	dir, err := os.MkdirTemp("", "example-*")
	if err != nil {
		return err
	}
	if !*debug {
		defer os.RemoveAll(dir)
	}

	inputpath := filepath.Join(dir, "input.pdf")
	input, err := os.OpenFile(inputpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	if _, err := input.ReadFrom(r); err != nil {
		return err
	}
	if err := input.Sync(); err != nil {
		return err
	}
	if err := input.Close(); err != nil {
		return err
	}

	ctx_, cancel := context.WithTimeout(ctx, *jobTimeout)
	defer cancel()
	c := exec.CommandContext(ctx_, "/run.sh", dir, size, interval)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}

	outputpath := filepath.Join(dir, "output.mp4")
	output, err := os.Open(outputpath)
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(w, output)
	return err
}

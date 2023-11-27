package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var ErrTooManyJobs = errors.New("job queue is full")

type Job struct {
	Cmd    *exec.Cmd
	Finish chan error
}

type Handler struct {
	Channel chan *Job
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
		if errors.Is(err, ErrTooManyJobs) {
			m = http.StatusTooManyRequests
		}
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
	if err := WriteFile(inputpath, r); err != nil {
		return fmt.Errorf("request read error %w", err)
	}

	ctx_, cancel := context.WithTimeout(ctx, *jobTimeout)
	defer cancel()
	c := exec.CommandContext(ctx_, "/run.sh", dir, size, interval)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := h.RunSequencial(c); err != nil {
		return err
	}

	outputpath := filepath.Join(dir, "output.mp4")
	if err := ReadFile(outputpath, w); err != nil {
		return fmt.Errorf("request write error %w", err)
	}
	return nil
}

// Run job in sequencial
func (h *Handler) RunSequencial(c *exec.Cmd) error {
	errCh := make(chan error, 1)
	defer close(errCh)
	j := &Job{Cmd: c, Finish: errCh}
	select {
	case h.Channel <- j:
	default:
		return ErrTooManyJobs
	}
	return <-errCh
}

func WriteFile(path string, r io.Reader) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.ReadFrom(r)
	return err
}

func ReadFile(path string, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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

	hash := sha256.New()
	tee := io.TeeReader(r, hash)
	inputpath := filepath.Join(dir, "input.pdf")
	if err := WriteFile(inputpath, tee); err != nil {
		return fmt.Errorf("request read error %w", err)
	}
	digest := hash.Sum(nil)
	name := fmt.Sprintf("%x.r%s.t%s.mp4", digest, size, interval)

	ctx_, cancel := context.WithTimeout(ctx, *jobTimeout)
	defer cancel()
	c := exec.CommandContext(ctx_, "/run.sh", dir, size, interval)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := h.RunSequencial(c); err != nil {
		return err
	}

	outputpath := filepath.Join(dir, "output.mp4")
	cachePath := filepath.Join(*cacheDir, name)
	if err = MoveFile(outputpath, cachePath); err != nil {
		return err
	}
	json.NewEncoder(w).Encode(map[string]string{
		"file":     name,
		"size":     size,
		"interval": interval,
	})
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

// MoveFile move file to dst
// If dst is inter-mountpoint, fallback copy and remove
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		// Success rename
		return nil
	}
	if !errors.Is(err, syscall.EXDEV) {
		// Other error
		return fmt.Errorf("rename operation error %w", err)
	}
	// Rename return EXDEV when inter-mountpoint
	// fallback copy and remove
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file error %w", err)
	}
	defer f.Close()
	if err := WriteFile(dst, f); err != nil {
		return fmt.Errorf("open destination file and copy error %w", err)
	}
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("remove source file %w", err)
	}
	return nil
}

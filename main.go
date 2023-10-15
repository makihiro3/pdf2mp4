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
)

var (
	listen = flag.String("listen", "./listen.socket", "listen unix domain socket")
	debug  = flag.Bool("debug", false, "debug flag")
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
	http.HandleFunc("/convert.cgi", HandleFunc)
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

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	// Accpet POST method only
	if r.Method != http.MethodPost {
		m := http.StatusMethodNotAllowed
		io.WriteString(w, http.StatusText(m))
		w.WriteHeader(m)
		return
	}

	if err := Process(w, r.Body); err != nil {
		log.Print(err)
		m := http.StatusInternalServerError
		w.WriteHeader(m)
		io.WriteString(w, http.StatusText(m))
		return
	}
	r.Body.Close()
}

func Process(w io.Writer, r io.Reader) error {
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
	input.ReadFrom(r)
	input.Close()

	outputpath := filepath.Join(dir, "output.mp4")
	c := exec.Command("/run.sh", inputpath, outputpath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}

	output, err := os.Open(outputpath)
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(w, output)
	return err
}

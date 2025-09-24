package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

func GrabageCollect(ctx context.Context, dir string, interval, expire time.Duration) error {
	log.Printf("Start GC: %s, interval: %s, expire: %s", dir, interval, expire)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Print(err)
			return err
		}
		th := time.Now().Add(-expire)
		log.Printf("Do GC: threshold: %v", th)
		for _, file := range files {
			st, err := file.Info()
			if err != nil {
				log.Print(err)
				return err
			}
			if st.ModTime().Before(th) {
				n := file.Name()
				log.Printf("Remove: %s", n)
				p := filepath.Join(dir, n)
				if err := os.Remove(p); err != nil {
					log.Print(err)
					return err
				}
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

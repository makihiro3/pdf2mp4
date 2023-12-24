package main

import (
	"context"
	"log"
	"os"
	"time"
)

func GrabageCollect(ctx context.Context, dir string, interval, expire time.Duration) error {
	log.Printf("Start GC: %s, interval: %s, expire: %s", dir, interval, expire)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		th := time.Now().Add(-expire)
		log.Printf("Do GC: threshold: %v", th)
		for _, file := range files {
			st, err := file.Info()
			if err != nil {
				return err
			}
			if st.ModTime().Before(th) {
				log.Printf("Remove: %s", file.Name())
				if err := os.Remove(file.Name()); err != nil {
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

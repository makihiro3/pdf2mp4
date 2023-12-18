package main

import (
	"context"
	"os"
	"time"
)

func GrabageCollect(ctx context.Context, dir string, interval, expire time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			files, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			th := time.Now().Add(-expire)
			for _, file := range files {
				st, err := file.Info()
				if err != nil {
					return err
				}
				if st.ModTime().Before(th) {
					if err := os.Remove(file.Name()); err != nil {
						return err
					}
				}
			}
		}
	}
}

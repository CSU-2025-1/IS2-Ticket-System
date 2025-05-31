package main

import (
	"log/slog"
	"stress-test/queries"
	"sync"
	"time"
)

const (
	writeRatio = 0.2
	token      = "ory_at_d3T4QQUtLX_HVf-Gp4yxtv-3ouUGK8hC-wXg9FM-PaQ.5EkNAjKPYgWa_CF8zGycTcPXzzSwxvzD_C1hHENevDQ"
	target     = "http://localhost:80"

	rpc = 1000

	readRatio = 1 - writeRatio
)

func main() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			wg := sync.WaitGroup{}
			wg.Add(rpc)
			for _ = range rpc {
				go func() {
					defer wg.Done()

					err := queries.DoRequest(target, token, writeRatio)
					if err != nil {
						slog.Error(err.Error())
					}
				}()
			}
			wg.Wait()
		}
	}
}

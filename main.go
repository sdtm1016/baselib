package main

import (
	"time"
	"baselib/logger"
)

func main() {

	go func() {
		for{
			logger.Info("goroutine loop...")
			time.Sleep(2 * time.Second)
		}
	}()

	for{
		logger.Info("main loop...")
		time.Sleep(3 * time.Second)
	}
}


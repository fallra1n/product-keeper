package shutdown

import (
	"io"
	"log"
	"os"
	"os/signal"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)
	sig := <-sigChan
	log.Printf("Caught signal %s. Shutting down...", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			log.Printf("failed to close %v: %v", closer, err)
		}
	}
}

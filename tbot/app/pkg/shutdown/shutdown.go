package shutdown

import (
	"io"
	"log"
	"os"
	"os/signal"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc
	log.Printf("Caught signal %s. Shutting down...", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			log.Printf("failed to close %v: %v", closer, err)
		}
	}
}

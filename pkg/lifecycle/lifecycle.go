package lifecycle

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func New(ctx context.Context, name string, onStart, onStop func(ctx context.Context) error) {
	lifeCtx, cancel := context.WithCancel(ctx)

	//Capture system signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		sig := <-c
		log.Printf("Received signal: %v, initiating shutdown", sig)
		cancel()
	}()

	g, gCtx := errgroup.WithContext(lifeCtx)

	g.Go(func() error {
		<-gCtx.Done()

		log.Printf("Initiating graceful shutdown for %s", name)

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()

		err := onStop(stopCtx)

		return err
	})
	if err := g.Wait(); err != nil{
		log.Printf("Shutdown completed with error from %s: %s", name, err)
	}else {
		log.Printf("Shutdown completed successfully for %s", name)
	}
}

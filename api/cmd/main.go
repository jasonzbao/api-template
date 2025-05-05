package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jasonzbao/api-template/api/server"
)

var (
	configFile = flag.String(
		"config",
		"./worker/configs/local/config.json",
		"config file")

	version = flag.String(
		"version",
		"",
		"git version",
	)
)

func main() {
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for sig := range c {
			log.Println("=> proc signal:", sig.String())
			cancel()
			time.Sleep(20 * time.Second)
			log.Println("=> proc exit")
			os.Exit(0)
		}
	}()

	fmt.Println("Starting initialize server")

	apiServer := server.NewServer(ctx, *configFile, *version)
	defer apiServer.Stop()

	fmt.Println("Finished initializing server. Starting http server...")

	// blocking call
	err := apiServer.Run(ctx, ":9001")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error running http server: %v", err)
	}
}

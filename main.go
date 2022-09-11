package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/powenyu/split-order/config"
	postgresql "github.com/powenyu/split-order/postgres"
	"github.com/powenyu/split-order/routes"
	"golang.org/x/sync/errgroup"
)

func init() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	postgresql.Initialize()
}

func main() {
	port := config.Port
	routesInit := routes.InitRouter()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        routesInit,
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		log.Println("[info] start http server listening", port)
		return server.ListenAndServe()
	})

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGILL, syscall.SIGFPE)
		select {
		case <-ctx.Done():
			return server.Shutdown(ctx)
		case s := <-c:
			close(c)
			if err := server.Shutdown(ctx); err != nil {
				return err
			}
			return fmt.Errorf("os signal: %v", s)
		}
	})

	if err := g.Wait(); err != nil {
		log.Println("[error]", err)
	}
	log.Println("[info] HTTP Server Exited")
	postgresql.Dispose()
}

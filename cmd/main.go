package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mtumbuu/url-shortener/pkg"
)

const serverPort = ":8000"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := run(); err != nil {
		log.Println(err)
	}
}

func run() error {
	handler := initHandler()
	http.HandleFunc("/url/shorten", handler.GetShortenedURL)
	http.HandleFunc("/url/full", handler.GetFullURL)
	http.HandleFunc("/url/shortened/", handler.GetRedirectURL)
	http.Handle("/notFound", http.NotFoundHandler())
	server := &http.Server{Addr: serverPort, Handler: nil}
	gracefulShutdownFor(server, os.Interrupt, os.Kill, syscall.SIGTERM)
	return server.ListenAndServe()
}

func initHandler() *pkg.HandlerImpl {
	service := pkg.NewServiceImpl()
	return pkg.NewHandlerImpl(service)
}

func gracefulShutdownFor(server *http.Server, signals ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, signals...)
	go func() {
		select {
		case <-c:
			log.Println("shutting down gracefully")
			server.Shutdown(context.Background())
		}
	}()
}

package main

import (
	"fmt"
	"log"
	"flag"
	"net/http"
	//"github.com/gorilla/mux"
	"github.com/nekonotes/router"
	"time"
	"os"
    "os/signal"
	"context"
    "github.com/rs/cors"

)



func main() {

	var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()
	fmt.Println("welcome to the go lang server")
	r :=router.Router()
	// Routes consist of a path and a handler function.
    cors := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedHeaders: []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
        AllowedMethods: []string{"GET", "PATCH", "POST", "PUT", "OPTIONS", "DELETE"},
    })
	//Bind to a port and pass our router in
	srv := &http.Server{
        Handler:      cors.Handler(r),
        Addr:         "0.0.0.0:7000",
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
    }
	// Run our server in a goroutine so that it doesn't block.
    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Println(err)
        }
    }()
	c := make(chan os.Signal, 1)
    // We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    srv.Shutdown(ctx)
    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
    log.Println("shutting down")
    os.Exit(0)
}

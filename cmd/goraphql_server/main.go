package main

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"

	"context"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pseudomuto/goraphql/pkg/graph"
	"github.com/pseudomuto/goraphql/pkg/storage"
)

var (
	listenAddr string
)

func main() {
	flag.StringVar(&listenAddr, "addr", ":8080", "The address to bind to")
	flag.Parse()

	logger := logrus.New()
	schema, _ := graph.NewSchema(graph.NewTodoNode(storage.NewEphemeralTodoRepo()))
	server := createServer(listenAddr, schema)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// turn off keep alives for new connections
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could'nt gracefully shutdown server: %v\n", err)
		}

		close(done)
	}()

	logger.Println("Server is ready and listening on", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Couldn't listen on %s: %v", listenAddr, err)
	}

	<-done
}

func createServer(addr string, schema graphql.Schema) *http.Server {
	svr := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := graphql.Do(graphql.Params{
				Schema:        schema,
				RequestString: r.URL.Query().Get("query"),
			})

			json.NewEncoder(w).Encode(result)
		}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return svr
}

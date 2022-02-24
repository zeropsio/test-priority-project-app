package main

import (
	"app/pkg/graceful"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := graceful.Context()
	defer cancel()

	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(handleIndex)
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", 8080),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()
	log.Println("server started")

	<-ctx.Done()
	log.Println("context done")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %+v\n", err)
	}
	log.Println("server shutdown")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "db", 6379),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Printf("connected to redis: %s", cn)
			return nil
		},
		DB: 1,
	})

	cmd := rdb.Get(r.Context(), "key_as_never_seen_before")
	if err := cmd.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, cmd.Val())
}

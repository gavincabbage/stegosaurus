package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gavincabbage/stegosaurus/image"

	"github.com/gavincabbage/stegosaurus"
)

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	terminated := make(chan os.Signal, 1)
	signal.Notify(terminated, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(terminated)
		close(terminated)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-terminated
		cancel()
	}()

	router := http.NewServeMux()
	router.Handle("/encode", encode())
	router.Handle("/decode", decode())
	router.Handle("/health", health())

	addr := ":8080"
	if a := os.Getenv("PORT"); a != "" {
		i, err := strconv.Atoi(a)
		if err != nil {
			logger.Printf("invalid port %s\n", a)
		}

		addr = fmt.Sprintf(":%d", i)
	}
	server := &http.Server{
		Addr:         addr,
		Handler:      logging(logger)(router),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	serverError := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverError <- err
		}
	}()

	select {
	case err := <-serverError:
		logger.Println("server error", err)
		cancel()
	case <-ctx.Done():
	}

	timeout, _ := context.WithTimeout(ctx, 5*time.Second)
	if err := server.Shutdown(timeout); err != nil {
		logger.Fatal(err)
	}
}

func encode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/encode" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var (
			encoder stegosaurus.Encoder
			payload io.Reader
			key     []byte
		)

		payloads, ok := r.URL.Query()["payload"]
		if ok && len(payloads) == 1 {
			payload = strings.NewReader(payloads[0])
		}
		keys, ok := r.URL.Query()["key"]
		if ok && len(keys) == 1 {
			key = []byte(keys[0])
		}

		switch r.Header.Get("Content-Type") {
		case "image/png", "image/jpeg":
			w.Header().Set("Content-Type", "image/png")
			encoder = image.NewEncoder(key)
		default:
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		if err := encoder.Encode(payload, r.Body, w); err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func decode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/decode" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var (
			decoder stegosaurus.Encoder
			key     []byte
		)

		keys, ok := r.URL.Query()["key"]
		if ok && len(keys) == 1 {
			key = []byte(keys[0])
		}

		switch r.Header.Get("Content-Type") {
		case "image/png":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			decoder = image.NewEncoder(key)
		default:
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		if err := decoder.Decode(r.Body, w); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(),
					"Content-Type:", r.Header.Get("content-type"),
					"Accept:", r.Header.Get("accept"))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

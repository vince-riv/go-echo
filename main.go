package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func echo(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(dump) // client may have disconnected; nothing to do
}

func accessLog(skipPath string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skipPath != "" && r.URL.Path == skipPath {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		fmt.Printf("%s %s %s %s %d %d %s\n",
			start.Format(time.RFC3339), r.RemoteAddr, r.Method, r.URL.RequestURI(),
			sw.status, sw.bytes, time.Since(start))
	})
}

func main() {
	addr := ":" + getenv("GO_ECHO_PORT", "8080")
	healthz := getenv("GO_ECHO_HEALTHZ_PATH", "/healthz")

	srv := &http.Server{
		Addr:              addr,
		Handler:           accessLog(healthz, http.HandlerFunc(echo)),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ctx.Done()
		stop() // a second signal kills the process immediately
		fmt.Println("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()

	fmt.Printf("listening on %s (healthz %s not logged)\n", addr, healthz)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-done
}

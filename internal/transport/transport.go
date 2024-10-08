package transport

import (
	"fmt"
	"metodichka/internal/config"
	"metodichka/internal/logger"
	"metodichka/internal/service"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

type Transport interface {
	Run() error
	Close() error
}

type transport struct {
	router  *http.Server
	service service.Service
}

func NewTransport(cfg config.Transport, s service.Service) Transport {
	router := http.Server{
		Addr: net.JoinHostPort(cfg.Host, cfg.Port),
	}

	tr := transport{
		router:  &router,
		service: s,
	}

	tr.setupRoutes()

	return &tr
}

func (t *transport) setupRoutes() {
	mx := http.NewServeMux()

	mx.HandleFunc("/fib", t.handle)
	mx.HandleFunc("/fib_cache", t.handleWithCache)

	t.router.Handler = mx
}

func (t *transport) Run() error {
	if err := t.router.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to listen and serve. %w", err)
	}

	return nil
}

func (t *transport) Close() error {
	if err := t.service.Close(); err != nil {
		return fmt.Errorf("failed to close server. %w", err)
	}

	return nil
}

func queryWrapper(queryValues url.Values) (int, error) {
	query := queryValues.Get("number")

	n, err := strconv.Atoi(query)
	if err != nil {
		return 0, fmt.Errorf("failed to convert value from query. %w", err)
	}

	return n, nil
}

func (t *transport) handle(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	n, err := queryWrapper(queryValues)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		w.Write([]byte("Failed to parse number from query"))

		return
	}

	logger.GetLogger().Printf("handle %d", n)

	result := t.service.GetFibonacci(n)

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(result))
}

func (t *transport) handleWithCache(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	n, err := queryWrapper(queryValues)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		w.Write([]byte("Failed to parse number from query"))

		return
	}

	logger.GetLogger().Printf("handle with cache %d", n)

	result := t.service.GetFibonacciCached(r.Context(), n)

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(result))
}

package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"net/http"
	"time"
)

type server struct {
	app    app.App
	logger logger.Logger
	srv    *http.Server
	mux    *http.ServeMux
}

func newServer(app app.App, logger logger.Logger) *server {
	s := &server{
		app:    app,
		logger: logger,
		mux:    http.NewServeMux(),
	}
	s.configureRouter()
	return s
}

func (s *server) Start(addr string) error {
	s.srv = &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	s.logger.Info("starting http server")
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}

func (s *server) configureRouter() {
	s.mux.HandleFunc("GET /hello", loggingMiddleware(s.handleHello, s.logger))
}

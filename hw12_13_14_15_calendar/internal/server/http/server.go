package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
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
	s.mux.HandleFunc("POST /event", loggingMiddleware(s.createEvent(s.app), s.logger))
	s.mux.HandleFunc("PUT /event/{id}", loggingMiddleware(s.updateEvent(s.app), s.logger))
	s.mux.HandleFunc("DELETE /event/{id}", loggingMiddleware(s.deleteEvent(s.app), s.logger))
	s.mux.HandleFunc("POST /listday", loggingMiddleware(s.getListDay(s.app), s.logger))
	s.mux.HandleFunc("POST /listweek", loggingMiddleware(s.getListWeek(s.app), s.logger))
	s.mux.HandleFunc("POST /listmonth", loggingMiddleware(s.getListMonth(s.app), s.logger))
}

type M map[string]interface{}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println(err)
	}
}

func serverError(w http.ResponseWriter, err error) {
	log.Println(err)
	errorResponse(w, http.StatusInternalServerError, "internal error")
}

func errorResponse(w http.ResponseWriter, code int, errs interface{}) {
	writeJSON(w, code, M{"errors": errs})
}

func httpEventToStorageEvent(event Event) storage.Event {
	return storage.Event{
		ID:           event.ID,
		Title:        event.Title,
		Start:        event.Start,
		Stop:         event.Stop,
		Description:  event.Description,
		UserID:       event.UserID,
		Notification: event.Notification,
	}
}

func storageEventToHTTPEvent(event storage.Event) Event {
	return Event{
		ID:           event.ID,
		Title:        event.Title,
		Start:        event.Start,
		Stop:         event.Stop,
		Description:  event.Description,
		UserID:       event.UserID,
		Notification: event.Notification,
	}
}

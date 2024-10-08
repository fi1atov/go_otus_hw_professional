package internalhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
)

func (s *server) handleHello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, world\n"))
}

func (s *server) createEvent(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение json-данных в структуру
		var id int
		event := Event{}

		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = event.Validate()
		if err != nil {
			writeJSON(w, http.StatusBadRequest, M{"error": err.Error()})
			return
		}
		// Создание события
		id, err = app.CreateEvent(r.Context(), event.UserID, event.Title, event.Description, event.Start, event.Stop, event.Notification) //nolint:lll
		if err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, CreateResult{id})
	}
}

func (s *server) updateEvent(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event Event
		// Получение json-данных в структуру
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = event.Validate()
		if err != nil {
			writeJSON(w, http.StatusBadRequest, M{"error": err.Error()})
			return
		}
		// Получаем ID события из URL и конвертируем в int
		eventID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatal(err)
		}
		// Обновление события
		change := httpEventToStorageEvent(event)
		err = app.UpdateEvent(r.Context(), eventID, change)
		if err != nil {
			log.Fatal(err)
		}

		writeJSON(w, http.StatusAccepted, OkResult{Ok: true})
	}
}

func (s *server) deleteEvent(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID события из URL и конвертируем в int
		eventID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatal(err)
		}
		// Удаление события
		err = app.DeleteEvent(r.Context(), eventID)
		if err != nil {
			log.Fatal(err)
		}

		writeJSON(w, http.StatusOK, OkResult{Ok: true})
	}
}

func (s *server) getListDay(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getList(w, r, app.ListDayEvent)
	}
}

func (s *server) getListWeek(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getList(w, r, app.ListWeekEvent)
	}
}

func (s *server) getListMonth(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getList(w, r, app.ListMonthEvent)
	}
}

func getList(w http.ResponseWriter, r *http.Request, fn app.ListEvents) {
	req := ListRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := fn(r.Context(), req.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := make(ListResult, 0, len(events))
	for _, event := range events {
		result = append(result, storageEventToHTTPEvent(event))
	}
	writeJSON(w, http.StatusOK, result)
}

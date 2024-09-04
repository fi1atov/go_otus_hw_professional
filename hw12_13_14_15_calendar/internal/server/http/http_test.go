package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage/storecreator"
	"github.com/stretchr/testify/suite"
)

type SuiteTest struct {
	suite.Suite
	ts   *httptest.Server
	app  app.App
	logg logger.Logger
	db   storage.Storage
}

func (s *SuiteTest) SetupTest() {
	ctx := context.Background()

	s.logg = logger.New("INFO", os.Stdout)

	dbConnect := os.Getenv("TEST")
	s.db, _ = storecreator.New(ctx, dbConnect == "", dbConnect)

	s.app = app.New(s.logg, s.db)

	s.ts = httptest.NewServer(newServer(s.app, s.logg).mux)

	_ = s.app.DeleteAllEvent(ctx)
}

func (s *SuiteTest) TearDownTest() {
	ctx := context.Background()
	s.ts.Close()
	_ = s.app.DeleteAllEvent(ctx)
	_ = s.db.Close(ctx)
}

func (s *SuiteTest) Post(endPoint string, data []byte) (resp *http.Response, err error) {
	res, err := http.Post(s.ts.URL+"/"+endPoint, "application/json", bytes.NewReader(data))
	return res, err
}

func (s *SuiteTest) Put(endPoint string, data []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPut, s.ts.URL+"/"+endPoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	return res, err
}

func (s *SuiteTest) NewCommonEvent() Event {
	eventStart := time.Now().Add(time.Hour * 2)
	eventStop := eventStart.Add(time.Hour)
	notification := 4 * time.Hour

	return Event{
		ID:           0,
		Title:        "событие",
		Start:        eventStart,
		Stop:         eventStop,
		Description:  "событие - описание",
		UserID:       1,
		Notification: &notification,
	}
}

func (s *SuiteTest) EqualEvents(event1, event2 Event) {
	s.Require().Equal(event1.Title, event2.Title)
	s.Require().Equal(event1.Description, event2.Description)
	s.Require().Equal(event1.Start.Unix(), event2.Start.Unix())
	s.Require().Equal(event1.Stop.Unix(), event2.Stop.Unix())
	s.Require().Equal(event1.UserID, event2.UserID)
	s.Require().Equal(event1.Notification, event2.Notification)
}

func (s *SuiteTest) AddEvent(event Event) int {
	data, _ := json.Marshal(event)

	res, err := http.Post(s.ts.URL+"/event", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, res.StatusCode)
	return s.readCreateID(res.Body)
}

func (s *SuiteTest) readCreateID(body io.ReadCloser) int {
	data, _ := io.ReadAll(body)
	defer body.Close()

	result := CreateResult{}
	err := json.Unmarshal(data, &result)
	s.Require().NoError(err)
	return result.ID
}

func (s *SuiteTest) readEvents(body io.ReadCloser) ListResult {
	data, _ := io.ReadAll(body)
	defer body.Close()

	result := ListResult{}
	err := json.Unmarshal(data, &result)
	s.Require().NoError(err)
	return result
}

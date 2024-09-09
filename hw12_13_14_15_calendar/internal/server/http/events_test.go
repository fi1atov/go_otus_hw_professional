package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type HTTPTest struct {
	SuiteTest
}

func (s *HTTPTest) TestCreate() {
	tests := []struct {
		name  string
		event Event
	}{
		{
			"with notification",
			s.NewCommonEvent(),
		},
		{
			"without notification",
			func() Event {
				event := s.NewCommonEvent()
				event.Notification = nil
				return event
			}(),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(_ *testing.T) {
			data, _ := json.Marshal(tt.event)

			res, err := s.Post("event", data)
			s.Require().NoError(err)
			defer res.Body.Close()

			s.Require().Equal(http.StatusCreated, res.StatusCode)
			id := s.readCreateID(res.Body)
			s.Require().Greater(id, 0)

			data, _ = json.Marshal(ListRequest{Date: tt.event.Start})

			res, err = s.Post("listday", data)
			s.Require().NoError(err)
			defer res.Body.Close()

			s.Require().Equal(http.StatusOK, res.StatusCode)
			events := s.readEvents(res.Body)
			s.Require().Equal(1, len(events))
			s.EqualEvents(tt.event, events[0])

			_ = s.app.DeleteAllEvent(context.Background())
		})
	}
}

func (s *HTTPTest) TestCreateFail() {
	res, err := s.Post("event", []byte("Hello, world\n"))
	s.Require().NoError(err)
	defer res.Body.Close()

	s.Require().Equal(http.StatusBadRequest, res.StatusCode)
}

func (s *HTTPTest) TestUpdate() {
	event := s.NewCommonEvent()
	id := s.AddEvent(event)

	event.ID = id
	event.Stop = event.Stop.Add(time.Hour)
	data, _ := json.Marshal(event)

	res, err := s.Put("event/"+strconv.Itoa(id), data) //nolint:bodyclose
	s.Require().NoError(err)
	s.Require().Equal(http.StatusAccepted, res.StatusCode)
}

func (s *HTTPTest) TestDelete() {
	event := s.NewCommonEvent()
	id := s.AddEvent(event)

	data, _ := json.Marshal(DeleteRequest{ID: id})

	res, err := s.Delete("event/"+strconv.Itoa(id), data) //nolint:bodyclose
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
}

func (s *HTTPTest) TestListDay() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	data, _ := json.Marshal(ListRequest{Date: event.Start})

	res, err := s.Post("listday", data) //nolint:bodyclose
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	events := s.readEvents(res.Body)
	s.Require().Equal(1, len(events))
	s.EqualEvents(event, events[0])
}

func (s *HTTPTest) TestListWeek() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	data, _ := json.Marshal(ListRequest{Date: event.Start})

	res, err := s.Post("listweek", data) //nolint:bodyclose
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	events := s.readEvents(res.Body)
	s.Require().Equal(1, len(events))
	s.EqualEvents(event, events[0])
}

func (s *HTTPTest) TestListMonth() {
	event := s.NewCommonEvent()
	s.AddEvent(event)

	data, _ := json.Marshal(ListRequest{Date: event.Start})

	res, err := s.Post("listmonth", data) //nolint:bodyclose
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	events := s.readEvents(res.Body)
	s.Require().Equal(1, len(events))
	s.EqualEvents(event, events[0])
}

func TestHttpCreateTest(t *testing.T) {
	suite.Run(t, new(HTTPTest))
}

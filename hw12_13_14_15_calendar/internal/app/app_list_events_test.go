package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/suite"
)

type ListEventTest struct {
	SuiteTest
}

func (s *ListEventTest) TestList() {
	ctx := context.Background()
	event1 := storage.Event{
		ID:           1,
		Title:        "Работа",
		Start:        time.Date(2024, 12, 13, 12, 42, 5, 0, time.UTC),
		Stop:         time.Date(2024, 12, 13, 13, 0, 0, 0, time.UTC),
		Description:  "Пойти поработать",
		UserID:       1,
		Notification: nil,
	}
	event2 := storage.Event{
		ID:           2,
		Title:        "Учеба",
		Start:        time.Date(2024, 12, 13, 17, 42, 5, 0, time.UTC),
		Stop:         time.Date(2024, 12, 13, 18, 0, 0, 0, time.UTC),
		Description:  "Делать домашки",
		UserID:       1,
		Notification: nil,
	}
	event3 := storage.Event{
		ID:           3,
		Title:        "Кушать",
		Start:        time.Date(2024, 12, 14, 9, 13, 17, 0, time.UTC),
		Stop:         time.Date(2024, 12, 14, 9, 15, 9, 0, time.UTC),
		Description:  "Пойти покушать",
		UserID:       1,
		Notification: nil,
	}
	event4 := storage.Event{
		ID:           4,
		Title:        "Смотреть сериал",
		Start:        time.Date(2024, 11, 14, 9, 13, 17, 0, time.UTC),
		Stop:         time.Date(2024, 11, 14, 9, 15, 9, 0, time.UTC),
		Description:  "Смотреть сериал",
		UserID:       1,
		Notification: nil,
	}

	_, err := s.AddEvent(event1)
	s.Require().NoError(err)
	_, err = s.AddEvent(event2)
	s.Require().NoError(err)
	_, err = s.AddEvent(event3)
	s.Require().NoError(err)
	_, err = s.AddEvent(event4)
	s.Require().NoError(err)

	// за 1 день
	list, err := s.calendar.ListDayEvent(ctx, event1.Start)
	s.Require().NoError(err)
	s.Require().Equal(2, len(list))
	s.EqualEvents(event1, list[0])
	s.EqualEvents(event2, list[1])

	// за другой день
	list, err = s.calendar.ListDayEvent(ctx, event3.Start)
	s.Require().NoError(err)
	s.Require().Equal(1, len(list))
	s.EqualEvents(event3, list[0])

	// за неделю
	list, err = s.calendar.ListWeekEvent(ctx, event1.Start)
	s.Require().NoError(err)
	s.Require().Equal(3, len(list))
	s.EqualEvents(event1, list[0])
	s.EqualEvents(event2, list[1])
	s.EqualEvents(event3, list[2])

	// за месяц
	list, err = s.calendar.ListWeekEvent(ctx, event1.Start)
	s.Require().NoError(err)
	s.Require().Equal(3, len(list))
	s.EqualEvents(event1, list[0])
	s.EqualEvents(event2, list[1])
	s.EqualEvents(event3, list[2])

	// за другой месяц
	list, err = s.calendar.ListWeekEvent(ctx, event4.Start)
	s.Require().NoError(err)
	s.Require().Equal(1, len(list))
	s.EqualEvents(event4, list[0])
}

func TestListEventTest(t *testing.T) {
	suite.Run(t, new(ListEventTest))
}

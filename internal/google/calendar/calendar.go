package calendar

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/calendar/v3"
)

type Event struct {
	ID 		string `json:"id"`
	Summary string `json:"summary"`
	ColorID  string `json:"color"`
	Start   *calendar.EventDateTime `json:"start"`
	End     *calendar.EventDateTime `json:"end"`
}

type Calendar struct {
	Service *calendar.Service
}

func NewCalendar(ctx context.Context, token *oauth2.Token) (*Calendar, error) {
	service, err := calendar.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(token)))
	if err != nil {
		return nil, err
	}

	return &Calendar{Service: service}, nil
}

func (c *Calendar) FetchEvents(startTime, endTime time.Time) ([]Event, error) {
	events, err := c.Service.Events.List("primary").
		TimeMin(startTime.Format(time.RFC3339)).
		TimeMax(endTime.Format(time.RFC3339)).
		Do()
	if err != nil {
		return nil, err
	}

	var eventsList []Event
	for _, item := range events.Items {
		event := Event{
			ID:      item.Id,
			Summary: item.Summary,
			ColorID: item.ColorId,
			Start:   item.Start,
			End:     item.End,
		}
		eventsList = append(eventsList, event)
	}


	return eventsList, nil
}

func (c *Calendar) InsertEvent(event *calendar.Event) (*calendar.Event, error) {
	event, err := c.Service.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Calendar) UpdateEvent(eventID string, event *calendar.Event) (*calendar.Event, error) {
	event, err := c.Service.Events.Update("primary", eventID, event).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Calendar) DeleteEvent(eventID string) error {
	err := c.Service.Events.Delete("primary", eventID).Do()
	if err != nil {
		return err
	}

	return nil
}
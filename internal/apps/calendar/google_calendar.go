package calendar

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/koo-arch/adjusta-backend/ent"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
)

type GoogleCalendarManager struct {
	client *ent.Client
}

func NewGoogleCalendarManager(client *ent.Client) *GoogleCalendarManager {
	return &GoogleCalendarManager{
		client: client,
	}
}

func (gcm *GoogleCalendarManager) FetchEventsFromCalendars(calendarService *customCalendar.Calendar, calendars []*ent.Calendar, startTime, endTime time.Time) ([]*models.Event, error) {
	var events []*models.Event

	for _, cal := range calendars {
		calEvents, err := calendarService.FetchEvents(cal.CalendarID, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch events from calendar: %s, error: %w", cal.Summary, err)
		}

		events = append(events, calEvents...)
	}

	return events, nil
}

func (gcm *GoogleCalendarManager) CreateGoogleEvents(calendarService *customCalendar.Calendar, eventReq *models.EventDraftCreation) ([]*calendar.Event, error) {
	// Googleカレンダーに登録されたイベントを追跡するスライス
	insertedGoogleEvents := make([]*calendar.Event, len(eventReq.SelectedDates))

	// 並列処理でイベントを登録
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(eventReq.SelectedDates))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(len(eventReq.SelectedDates))
	for i, date := range eventReq.SelectedDates {
		go func(i int, date models.SelectedDate) {
			defer wg.Done()

			select {
			case <-ctx.Done(): // エラー発生時に他のゴルーチンをキャンセル
				return
			default:
				event := gcm.ConvertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, date.Start, date.End)

				insertEvent, err := calendarService.InsertEvent(event)
				if err != nil {
					errCh <- fmt.Errorf("failed to insert event to Google Calendar: %w", err)
					cancel() // エラーが発生したら他のゴルーチンをキャンセル
					return
				}

				mu.Lock()
				insertedGoogleEvents[i] = insertEvent
				mu.Unlock()
			}
		}(i, date)
	}

	wg.Wait()
	close(errCh)

	// エラーが発生していた場合、登録したイベントを削除
	var errList []error
	for err := range errCh {
		errList = append(errList, err)
	}

	if len(errList) > 0 {
		if delErr := gcm.DeleteGoogleEvents(calendarService, insertedGoogleEvents); delErr != nil {
			return nil, fmt.Errorf("failed to delete events from Google Calendar: %w", delErr)
		}
		return nil, fmt.Errorf("multiple errors occurred: %v", errList)
	}

	return insertedGoogleEvents, nil
}

func (gcm *GoogleCalendarManager) UpdateGoogleCalendarEvents(calendarService *customCalendar.Calendar, eventReq *models.EventDraftDetail) error {
	// Googleカレンダーに登録されたイベントを追跡するスライス
	var backupGoogleEvents []*calendar.Event

	// 並列処理でイベントを登録
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(eventReq.ProposedDates))

	wg.Add(len(eventReq.ProposedDates))
	for _, date := range eventReq.ProposedDates {
		go func(date models.ProposedDate) {
			defer wg.Done()

			event := gcm.ConvertToCalendarEvent(&date.GoogleEventID, eventReq.Title, eventReq.Location, eventReq.Description, *date.Start, *date.End)

			// 更新前にイベントをバックアップできるように、Googleカレンダーからイベントを取得
			backupEvent, err := calendarService.FetchEvent(date.GoogleEventID)
			if err != nil {
				errCh <- fmt.Errorf("failed to fetch event from Google Calendar: %w", err)
				return
			}

			_, err = calendarService.UpdateEvent(date.GoogleEventID, event)
			if err != nil {
				errCh <- fmt.Errorf("failed to insert event to Google Calendar: %w", err)
				return
			}

			mu.Lock()
			backupGoogleEvents = append(backupGoogleEvents, backupEvent)
			mu.Unlock()
		}(date)
	}

	wg.Wait()
	close(errCh)

	// エラーが発生していた場合、更新したイベントを元に戻す
	for err := range errCh {
		if err != nil {
			for _, event := range backupGoogleEvents {
				if _, err := calendarService.UpdateEvent(event.Id, event); err != nil {
					return fmt.Errorf("failed to update event to Google Calendar: %w", err)
				}
			}
			return err
		}
	}

	return nil
}

func (gcm *GoogleCalendarManager) DeleteGoogleEvents(calendarService *customCalendar.Calendar, events []*calendar.Event) error {
	for _, event := range events {
		if event == nil || event.Id == "" {
			continue // eventまたはIDがnilの場合はスキップ
		}

		err := calendarService.DeleteEvent(event.Id)
		if err != nil {
			if gErr, ok := err.(*googleapi.Error); ok && gErr.Code == 410 {
				// 410エラーはリソースが既に削除されているため、無視
				fmt.Printf("Warning: Event ID %s is already deleted.\n", event.Id)
				continue
			}
			return fmt.Errorf("failed to delete event with ID %s: %w", event.Id, err)
		}
	}
	return nil
}

func (gcm *GoogleCalendarManager) ConvertToCalendarEvent(ID *string, title, location, description string, start, end time.Time) *calendar.Event {
	event := &calendar.Event{
		Summary:     title,
		Location:    location,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Asia/Tokyo",
		},
	}
	// IDがnilでない場合のみ設定
	if ID != nil && *ID != "" {
		event.Id = *ID
	}

	return event
}

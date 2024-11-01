package event_operations

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"google.golang.org/api/calendar/v3"
)

type EventDeleteManager struct {
	event *events.EventManager
}

func NewEventDeleteManager(event *events.EventManager) *EventDeleteManager {
	return &EventDeleteManager{
		event: event,
	}
}

func (edm *EventDeleteManager) DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftDetail) error {
	tx, err := edm.event.Client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	token, err := edm.event.AuthManager.VerifyOAuthToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to verify token for account: %s, error: %w", email, err)
	}

	calendarService, err := customCalendar.NewCalendar(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to create calendar service for account: %s, error: %w", email, err)
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	_, err = edm.event.CalendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	// Googleカレンダーからイベントを削除
	backupEvents, err := edm.event.CalendarApp.DeleteGoogleCalendarEvents(calendarService, eventReq)
	if err != nil {
		return fmt.Errorf("failed to delete google calendar events for account: %s, error: %w", email, err)
	}

	// データベースからイベントを削除
	err = edm.event.EventRepo.Delete(ctx, tx, eventReq.ID)
	if err != nil {
		// 削除されたGoogleカレンダーのイベントを復元
		if rollbackErr := edm.rollbackGoogleEvents(calendarService, backupEvents); rollbackErr != nil {
			log.Printf("failed to rollback google calendar events for account: %s, error: %v", email, rollbackErr)
		}
		return fmt.Errorf("failed to delete event for account: %s, error: %w", email, err)
	}

	return nil
}

func (edm *EventDeleteManager) rollbackGoogleEvents(calendarService *customCalendar.Calendar, backupEvents []*calendar.Event) error {
	// Googleカレンダーからイベントを復元
	for _, event := range backupEvents {
		_, err := calendarService.InsertEvent(event)
		if err != nil {
			return fmt.Errorf("failed to rollback google calendar events, error: %w", err)
		}
	}

	return nil
}
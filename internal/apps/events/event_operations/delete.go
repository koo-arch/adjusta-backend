package event_operations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
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

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	_, err = edm.event.CalendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}


	// データベースからイベントを削除
	err = edm.event.EventRepo.Delete(ctx, tx, eventReq.ID)
	if err != nil {
		return fmt.Errorf("failed to delete event for account: %s, error: %w", email, err)
	}

	return nil
}

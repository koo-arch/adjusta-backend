package calendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/account"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
)

type CalendarRepositoryImpl struct {
	client *ent.Client
}

func NewCalendarRepository(client *ent.Client) *CalendarRepositoryImpl {
	return &CalendarRepositoryImpl{
		client: client,
	}
}

func (r *CalendarRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Calendar, error) {
	if tx != nil {
		return tx.Calendar.Get(ctx, id)
	}
	return r.client.Calendar.Get(ctx, id)
}

func (r *CalendarRepositoryImpl) FilterByAccountID(ctx context.Context, tx *ent.Tx, accountID uuid.UUID) ([]*ent.Calendar, error) {
	fileterCalendar := r.client.Calendar.Query()
	if tx != nil {
		fileterCalendar = tx.Calendar.Query()
	}
	return fileterCalendar.
		Where(calendar.HasAccountWith(account.ID(accountID))).
		All(ctx)
}

func (r *CalendarRepositoryImpl) FindByAccountIDAndCalendarID(ctx context.Context, tx *ent.Tx, accountID uuid.UUID, calendarID string) (*ent.Calendar, error) {
	findCalendar := r.client.Calendar.Query()
	if tx != nil {
		findCalendar = tx.Calendar.Query()
	}
	return findCalendar.
		Where(
			calendar.HasAccountWith(account.IDEQ(accountID)),
			calendar.CalendarIDEQ(calendarID),
		).
		Only(ctx)
}

func (r *CalendarRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, calendarID string, summary string, account *ent.Account) (*ent.Calendar, error) {
	createCalendar := r.client.Calendar.Create()
	if tx != nil {
		createCalendar = tx.Calendar.Create()
	}
	return createCalendar.
		SetCalendarID(calendarID).
		SetSummary(summary).
		SetAccount(account).
		Save(ctx)
}

func (r *CalendarRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, summary string) (*ent.Calendar, error) {
	updateCalendar := r.client.Calendar.UpdateOneID(id)
	if tx != nil {
		updateCalendar = tx.Calendar.UpdateOneID(id)
	}
	return updateCalendar.
		SetSummary(summary).
		Save(ctx)
}

func (r *CalendarRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.Calendar.DeleteOneID(id).Exec(ctx)
	}
	return r.client.Calendar.DeleteOneID(id).Exec(ctx)
}
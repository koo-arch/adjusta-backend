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

func (r *CalendarRepositoryImpl) FindByFields(ctx context.Context, tx *ent.Tx, accountID uuid.UUID, calendarID, summary *string, isPrimary *bool) (*ent.Calendar, error) {
	findCalendar := r.client.Calendar.Query()
	if tx != nil {
		findCalendar = tx.Calendar.Query()
	}
	query := findCalendar.
		Where(
			calendar.HasAccountWith(account.IDEQ(accountID)),
		)
	if calendarID != nil {
		query = query.Where(calendar.CalendarIDEQ(*calendarID))
	}
	if summary != nil {
		query = query.Where(calendar.SummaryEQ(*summary))
	}
	if isPrimary != nil {
		query = query.Where(calendar.IsPrimaryEQ(*isPrimary))
	}
	return query.Only(ctx)
}

func (r *CalendarRepositoryImpl) FilterByFields(ctx context.Context, tx *ent.Tx, accountID uuid.UUID, calendarID, summary *string, isPrimary *bool) ([]*ent.Calendar, error) {
	filterCalendar := r.client.Calendar.Query()
	if tx != nil {
		filterCalendar = tx.Calendar.Query()
	}
	query := filterCalendar.
		Where(
			calendar.HasAccountWith(account.IDEQ(accountID)),
		)
	if calendarID != nil {
		query = query.Where(calendar.CalendarIDEQ(*calendarID))
	}
	if summary != nil {
		query = query.Where(calendar.SummaryEQ(*summary))
	}
	if isPrimary != nil {
		query = query.Where(calendar.IsPrimaryEQ(*isPrimary))
	}
	return query.All(ctx)
}

func (r *CalendarRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, calendarID string, summary string, is_primary bool, account *ent.Account) (*ent.Calendar, error) {
	createCalendar := r.client.Calendar.Create()
	if tx != nil {
		createCalendar = tx.Calendar.Create()
	}
	return createCalendar.
		SetCalendarID(calendarID).
		SetSummary(summary).
		SetAccount(account).
		SetIsPrimary(is_primary).
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
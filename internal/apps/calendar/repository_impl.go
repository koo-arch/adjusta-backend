package calendar

import (
	"context"
	"fmt"

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

func (r *CalendarRepositoryImpl) FindByFields(ctx context.Context, tx *ent.Tx, accountID uuid.UUID, opt CalendarQueryOptions) (*ent.Calendar, error) {
	if !opt.WithEvents && opt.WithProposedDates {
		return nil, fmt.Errorf("WithDates is only available when withEvents is true")
	}

	findCalendar := r.client.Calendar.Query()
	if tx != nil {
		findCalendar = tx.Calendar.Query()
	}
	query := r.applyCalendarQueryOptions(findCalendar, accountID, opt)

	return query.Only(ctx)
}

func (r *CalendarRepositoryImpl) FilterByFields(ctx context.Context, tx *ent.Tx, accountID uuid.UUID, opt CalendarQueryOptions) ([]*ent.Calendar, error) {
	filterCalendar := r.client.Calendar.Query()
	
	if !opt.WithEvents && opt.WithProposedDates {
		return nil, fmt.Errorf("WithDates is only available when withEvents is true")
	}
	if tx != nil {
		filterCalendar = tx.Calendar.Query()
	}
	
	query := r.applyCalendarQueryOptions(filterCalendar, accountID, opt)

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

func (r *CalendarRepositoryImpl) applyCalendarQueryOptions(query *ent.CalendarQuery, accountID uuid.UUID, opt CalendarQueryOptions) *ent.CalendarQuery {
	query = query.Where(calendar.HasAccountWith(account.IDEQ(accountID)))

	if opt.CalendarID != nil {
		query = query.Where(calendar.CalendarIDEQ(*opt.CalendarID))
	}
	if opt.Summary != nil {
		query = query.Where(calendar.SummaryEQ(*opt.Summary))
	}
	if opt.IsPrimary != nil {
		query = query.Where(calendar.IsPrimaryEQ(*opt.IsPrimary))
	}

	if opt.WithEvents {
		query = query.WithEvents(func(eventQuery *ent.EventQuery) {
			if opt.EventOffset > 0 {
				eventQuery = eventQuery.Offset(opt.EventOffset)
			}
			if opt.EventLimit > 0 {
				eventQuery = eventQuery.Limit(opt.EventLimit)
			}

			if opt.WithProposedDates {
				eventQuery = eventQuery.WithProposedDates(func(dateQuery *ent.ProposedDateQuery) {
					if opt.ProposedDateOffset > 0 {
						dateQuery = dateQuery.Offset(opt.ProposedDateOffset)
					}
					if opt.ProposedDateLimit > 0 {
						dateQuery = dateQuery.Limit(opt.ProposedDateLimit)
					}
				})
			}
		})
	}

	return query
}
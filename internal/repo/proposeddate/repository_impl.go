package proposeddate

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type ProposedDateRepositoryImpl struct {
	client *ent.Client
}

func NewProposedDateRepository(client *ent.Client) *ProposedDateRepositoryImpl {
	return &ProposedDateRepositoryImpl{
		client: client,
	}
}

func (r *ProposedDateRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.ProposedDate, error) {
	if tx != nil {
		return tx.ProposedDate.Get(ctx, id)
	}
	return r.client.ProposedDate.Get(ctx, id)
}

func (r *ProposedDateRepositoryImpl) FilterByEventID(ctx context.Context, tx *ent.Tx, eventID uuid.UUID) ([]*ent.ProposedDate, error) {
	filterProposedDate := r.client.ProposedDate.Query()
	if tx != nil {
		filterProposedDate = tx.ProposedDate.Query()
	}
	return filterProposedDate.
		Where(proposeddate.HasEventWith(event.IDEQ(eventID))).
		All(ctx)
}

func (r* ProposedDateRepositoryImpl) ExclusionEventID(ctx context.Context, tx *ent.Tx, eventID uuid.UUID) ([]*ent.ProposedDate, error) {
	filterProposedDate := r.client.ProposedDate.Query()
	if tx != nil {
		filterProposedDate = tx.ProposedDate.Query()
	}
	return filterProposedDate.
		Where(proposeddate.Not(proposeddate.HasEventWith(event.IDEQ(eventID)))).
		All(ctx)
}

func (r *ProposedDateRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, googleEventID *string, opt ProposedDateQueryOptions, entEvent *ent.Event) (*ent.ProposedDate, error) {
	proposedDateCreate := r.client.ProposedDate.Create()
	if tx != nil {
		proposedDateCreate = tx.ProposedDate.Create()
	}

	if googleEventID != nil {
		proposedDateCreate = proposedDateCreate.SetGoogleEventID(*googleEventID)
	}

	proposedDateCreate = proposedDateCreate.
		SetStartTime(*opt.StartTime).
		SetEndTime(*opt.EndTime).
		SetPriority(*opt.Priority).
		SetEvent(entEvent)

	return proposedDateCreate.Save(ctx)
}

func (r *ProposedDateRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt ProposedDateQueryOptions) (*ent.ProposedDate, error) {
	proposedDateUpdate := r.client.ProposedDate.UpdateOneID(id)
	if tx != nil {
		proposedDateUpdate = tx.ProposedDate.UpdateOneID(id)
	}

	if opt.GoogleEventID != nil {
		proposedDateUpdate = proposedDateUpdate.SetGoogleEventID(*opt.GoogleEventID)
	}

	if opt.StartTime!= nil {
		proposedDateUpdate = proposedDateUpdate.SetStartTime(*opt.StartTime)
	}

	if opt.EndTime != nil {
		proposedDateUpdate = proposedDateUpdate.SetEndTime(*opt.EndTime)
	}

	if opt.Priority != nil {
		proposedDateUpdate = proposedDateUpdate.SetPriority(*opt.Priority)
	}

	return proposedDateUpdate.Save(ctx)
}

func (r *ProposedDateRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.ProposedDate.DeleteOneID(id).Exec(ctx)
	}
	return r.client.ProposedDate.DeleteOneID(id).Exec(ctx)
}


func (r *ProposedDateRepositoryImpl) CreateBulk(ctx context.Context, tx *ent.Tx, selectedDates []models.SelectedDate, googleEvents []*calendar.Event, entEvent *ent.Event) ([]*ent.ProposedDate, error) {
	var proposedDateCreates []*ent.ProposedDateCreate

	for i, selectedDate := range selectedDates {
		proposedDateCreate := r.client.ProposedDate.Create()
		if tx != nil {
			proposedDateCreate = tx.ProposedDate.Create()
		}

		proposedDateCreate = proposedDateCreate.
			SetStartTime(selectedDate.Start).
			SetEndTime(selectedDate.End).
			SetPriority(selectedDate.Priority).
			SetEvent(entEvent)

		if googleEvents != nil {
			proposedDateCreate = proposedDateCreate.SetGoogleEventID(googleEvents[i].Id)
		}

		proposedDateCreates = append(proposedDateCreates, proposedDateCreate)
	}

	if (tx != nil) {
		return tx.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
	}

	return r.client.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
}

func (r *ProposedDateRepositoryImpl) UpdateByGoogleEventID(ctx context.Context, tx *ent.Tx, oldGoogleEvent *string, opt ProposedDateQueryOptions) error {
	update := r.client.ProposedDate.Update()
	if tx != nil {
		update = tx.ProposedDate.Update()
	}

	update = update.Where(proposeddate.GoogleEventIDEQ(*oldGoogleEvent))

	if opt.GoogleEventID != nil {
		update = update.SetGoogleEventID(*opt.GoogleEventID)
	}

	if opt.StartTime!= nil {
		update = update.SetStartTime(*opt.StartTime)
	}

	if opt.EndTime != nil {
		update = update.SetEndTime(*opt.EndTime)
	}

	if opt.Priority != nil {
		update = update.SetPriority(*opt.Priority)
	}

	_, err := update.Save(ctx)
	return err
}
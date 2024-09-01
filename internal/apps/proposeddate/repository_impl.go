package proposeddate

import (
	"context"
	"time"

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

func (r *ProposedDateRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, googleEventID *string, startTime, endTime time.Time, priority int, entEvent *ent.Event) (*ent.ProposedDate, error) {
	proposedDateCreate := r.client.ProposedDate.Create()
	if tx != nil {
		proposedDateCreate = tx.ProposedDate.Create()
	}

	if googleEventID != nil {
		proposedDateCreate = proposedDateCreate.SetGoogleEventID(*googleEventID)
	}

	proposedDateCreate = proposedDateCreate.
		SetStartTime(startTime).
		SetEndTime(endTime).
		SetPriority(priority).
		SetEvent(entEvent)

	return proposedDateCreate.Save(ctx)
}

func (r *ProposedDateRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, googleEventID *string, startTime, endTime *time.Time, priority *int, isFinalized *bool) (*ent.ProposedDate, error) {
	proposedDateUpdate := r.client.ProposedDate.UpdateOneID(id)
	if tx != nil {
		proposedDateUpdate = tx.ProposedDate.UpdateOneID(id)
	}

	if googleEventID != nil {
		proposedDateUpdate = proposedDateUpdate.SetGoogleEventID(*googleEventID)
	}

	if startTime != nil {
		proposedDateUpdate = proposedDateUpdate.SetStartTime(*startTime)
	}

	if endTime != nil {
		proposedDateUpdate = proposedDateUpdate.SetEndTime(*endTime)
	}

	if priority != nil {
		proposedDateUpdate = proposedDateUpdate.SetPriority(*priority)
	}

	if isFinalized != nil {
		proposedDateUpdate = proposedDateUpdate.SetIsFinalized(*isFinalized)
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
package proposeddate

import (
	"context"

	"github.com/google/uuid"
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

func (r *ProposedDateRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, selectedDates models.SelectedDate, entEvent *ent.Event, priority int) (*ent.ProposedDate, error) {
	proposedDateCreate := r.client.ProposedDate.Create()
	if tx != nil {
		proposedDateCreate = tx.ProposedDate.Create()
	}

	proposedDateCreate = proposedDateCreate.
		SetStartTime(selectedDates.Start).
		SetEndTime(selectedDates.End).
		SetPriority(priority).
		SetEvent(entEvent)

	return proposedDateCreate.Save(ctx)
}

func (r *ProposedDateRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, selectedDates *models.SelectedDate, priority *int, isFinalized *bool) (*ent.ProposedDate, error) {
	proposedDateUpdate := r.client.ProposedDate.UpdateOneID(id)
	if tx != nil {
		proposedDateUpdate = tx.ProposedDate.UpdateOneID(id)
	}

	if selectedDates != nil {
		proposedDateUpdate = proposedDateUpdate.
			SetStartTime(selectedDates.Start).
			SetEndTime(selectedDates.End)
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

func (r *ProposedDateRepositoryImpl) CreateBulk(ctx context.Context, tx *ent.Tx, selectedDates []models.SelectedDate, entEvent *ent.Event) ([]*ent.ProposedDate, error) {
	var proposedDateCreates []*ent.ProposedDateCreate

	for index, selectedDate := range selectedDates {
		proposedDateCreate := r.client.ProposedDate.Create()
		if tx != nil {
			proposedDateCreate = tx.ProposedDate.Create()
		}
		
		proposedDateCreate.
			SetStartTime(selectedDate.Start).
			SetEndTime(selectedDate.End).
			SetPriority(index+1).
			SetEvent(entEvent)
	
		proposedDateCreates = append(proposedDateCreates, proposedDateCreate)

	}
	if (tx != nil) {
		return tx.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
	}

	return r.client.ProposedDate.CreateBulk(proposedDateCreates...).Save(ctx)
}
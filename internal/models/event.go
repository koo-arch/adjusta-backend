package models

import (
	"github.com/google/uuid"

	"time"
)

type AccountsEvents struct {
	AccountID uuid.UUID `json:"account_id"`
	Email     string    `json:"email"`
	Events    []*Event  `json:"events"`
}

type Event struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
	ColorID string `json:"color"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

type EventDraftCreation struct {
	Title         string         `json:"title" binding:"required"`
	Location      string         `json:"location"`
	Description   string         `json:"description"`
	SelectedDates []SelectedDate `json:"selected_dates" binding:"required"`
}

type SelectedDate struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Priority int       `json:"priority"`
}


type EventDraftDetail struct {
	ID            uuid.UUID      `json:"id" binding:"required"`
	Title         string         `json:"title"`
	Location      string         `json:"location"`
	Description   string         `json:"description"`
	ProposedDates []ProposedDate `json:"proposed_dates"`
}

type ProposedDate struct {
	ID            uuid.UUID  `json:"id"`
	GoogleEventID string     `json:"event_id"`
	Start         *time.Time `json:"start_date"`
	End           *time.Time `json:"end_date"`
	Priority      int        `json:"priority"`
	IsFinalized   bool       `json:"is_finalized"`
}


type ConfirmEvent struct {
	ConfirmDate ConfirmDate `json:"confirm_date" binding:"required"`
}

type ConfirmDate struct {
	ID            *uuid.UUID `json:"id"`
	GoogleEventID string     `json:"event_id"`
	Start         *time.Time `json:"start_date"`
	End           *time.Time `json:"end_date"`
	Priority      int        `json:"priority"`
}

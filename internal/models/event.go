package models

import (
	"time"
	
	"github.com/google/uuid"
)

type EventStatus string

const (
	StatusPending   EventStatus = "pending"
	StatusConfirmed EventStatus = "confirmed"
	StatusCancelled EventStatus = "cancelled"
)

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
	Status        EventStatus    `json:"status"`
	ConfirmedDateID *uuid.UUID    `json:"confirmed_date_id"`
	GoogleEventID string     `json:"google_event_id"`
	ProposedDates []ProposedDate `json:"proposed_dates"`
}

type ProposedDate struct {
	ID            *uuid.UUID  `json:"id"`
	Start         *time.Time `json:"start"`
	End           *time.Time `json:"end"`
	Priority      int        `json:"priority"`
}


type ConfirmEvent struct {
	ConfirmDate ConfirmDate `json:"confirm_date" binding:"required"`
}

type ConfirmDate struct {
	ID            *uuid.UUID `json:"id"`
	GoogleEventID string     `json:"google_event_id"`
	Start         *time.Time `json:"start"`
	End           *time.Time `json:"end"`
	Priority      int        `json:"priority"`
}

type EventDraftQueryOptions struct {
	Title 	 string
	Location string
	Description string
	Status   EventStatus
	StartDate time.Time
	EndDate   time.Time
	GoogleEventID string
}
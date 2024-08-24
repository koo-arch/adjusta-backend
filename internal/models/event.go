package models

import (
	"github.com/google/uuid"

	"time"
)

type AccountsEvents struct {
	AccountID uuid.UUID   `json:"account_id"`
	Email     string      `json:"email"`
	Events    []*Event    `json:"events"`
}

type Event struct {
	ID 		string `json:"id"`
	Summary string `json:"summary"`
	ColorID  string `json:"color"`
	Start   string `json:"start"`
	End     string`json:"end"`
}


type EventDraft struct {
	Title    string    `json:"title" binding:"required"`
	Location string    `json:"location"`
	Description string `json:"description"`
	SelectedDates []SelectedDate `json:"selected_dates" binding:"required"`
}

type SelectedDate struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
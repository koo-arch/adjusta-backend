package queryparser

import (
	"fmt"

	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/models"

)

func(qp *QueryParser) ParseSearchEventQuery() (*event.EventQueryOptions, error) {

	title, err := qp.ParseString("title")
	if err != nil {
		return nil, fmt.Errorf("failed to parse title: %w", err)
	}

	location, err := qp.ParseString("location")
	if err != nil {
		return nil, fmt.Errorf("failed to parse location: %w", err)
	}

	description, err := qp.ParseString("description")
	if err != nil {
		return nil, fmt.Errorf("failed to parse description: %w", err)
	}

	status, err := qp.ParseString("status")
	if err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}
	eventStatus, err := qp.vaildateStatus(status)
	if err != nil {
		return nil, fmt.Errorf("failed to validate status: %w", err)
	}

	startTime, err := qp.ParseTime("start_time")
	if err != nil {
		return nil, fmt.Errorf("failed to parse start_time: %w", err)
	}

	endTime, err := qp.ParseTime("end_time")
	if err != nil {
		return nil, fmt.Errorf("failed to parse end_time: %w", err)
	}

	options := event.EventQueryOptions{
		Summary: title,
		Location: location,
		Description: description,
		Status: eventStatus,
		ProposedDateStartTime: startTime,
		ProposedDateEndTime: endTime,
	}

	return &options, nil
}

func(qp *QueryParser) vaildateStatus(status *string) (*models.EventStatus, error) {
	if status == nil {
		return nil, nil
	}

	var result models.EventStatus

	switch *status {
	case "pending":
		result = models.StatusPending
	case "confirmed":
		result = models.StatusConfirmed
	case "cancelled":
		result = models.StatusCancelled
	default:
		return nil, fmt.Errorf("invalid status: %s", *status)
	}

	return &result, nil
}
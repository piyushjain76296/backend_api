package services

import (
	"errors"

	"ticket-system/internal/core/models"
	"ticket-system/internal/infrastructure/repositories"
)

var (
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrUnauthorized      = errors.New("unauthorized to access this ticket")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrTitleRequired     = errors.New("title is required")
	ErrDescRequired      = errors.New("description is required")
)

type TicketService struct {
	ticketRepo *repositories.TicketRepository
}

func NewTicketService(ticketRepo *repositories.TicketRepository) *TicketService {
	return &TicketService{ticketRepo: ticketRepo}
}

func (s *TicketService) CreateTicket(userID int, title, description string) (*models.Ticket, error) {
	if title == "" {
		return nil, ErrTitleRequired
	}
	if description == "" {
		return nil, ErrDescRequired
	}

	ticket := &models.Ticket{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      models.StatusOpen, // Enforce initial state
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}

	// Fetch full ticket to get DB generated timestamps
	return s.ticketRepo.GetByID(ticket.ID)
}

func (s *TicketService) GetUserTickets(userID int) ([]models.Ticket, error) {
	return s.ticketRepo.GetAllByUserID(userID)
}

func (s *TicketService) GetTicketByID(userID, ticketID int) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.GetByID(ticketID)
	if err != nil {
		if errors.Is(err, repositories.ErrTicketNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	// Enforce ownership
	if ticket.UserID != userID {
		// Return not found to not leak existence of ticket
		return nil, ErrTicketNotFound
	}

	return ticket, nil
}

func (s *TicketService) UpdateTicketStatus(userID, ticketID int, newStatus string) error {
	// 1. Validate requested status string
	if newStatus != models.StatusOpen && newStatus != models.StatusInProgress && newStatus != models.StatusClosed {
		return ErrInvalidStatus
	}

	// 2. Fetch ticket and check ownership
	ticket, err := s.ticketRepo.GetByID(ticketID)
	if err != nil {
		if errors.Is(err, repositories.ErrTicketNotFound) {
			return ErrTicketNotFound
		}
		return err
	}

	if ticket.UserID != userID {
		// Return not found to not leak existence
		return ErrTicketNotFound
	}

	// 3. Validate transition rules
	if ticket.Status == models.StatusClosed {
		// Closed tickets can never be reopened or changed
		return ErrInvalidTransition
	}

	if ticket.Status == models.StatusOpen && newStatus != models.StatusInProgress {
		// From open, can only go to in_progress
		return ErrInvalidTransition
	}

	if ticket.Status == models.StatusInProgress && newStatus != models.StatusClosed {
		// From in_progress, can only go to closed
		return ErrInvalidTransition
	}

	// 4. Perform update
	return s.ticketRepo.UpdateStatus(ticketID, newStatus)
}

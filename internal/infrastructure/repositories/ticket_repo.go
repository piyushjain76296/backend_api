package repositories

import (
	"database/sql"
	"errors"

	"ticket-system/internal/core/models"
)

var ErrTicketNotFound = errors.New("ticket not found")

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	query := `INSERT INTO tickets (user_id, title, description, status) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, ticket.UserID, ticket.Title, ticket.Description, ticket.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	ticket.ID = int(id)
	return nil
}

func (r *TicketRepository) GetAllByUserID(userID int) ([]models.Ticket, error) {
	query := `SELECT id, user_id, title, description, status, created_at, updated_at FROM tickets WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	// Return empty array instead of nil for better JSON representation
	if tickets == nil {
		tickets = []models.Ticket{}
	}

	return tickets, nil
}

func (r *TicketRepository) GetByID(id int) (*models.Ticket, error) {
	query := `SELECT id, user_id, title, description, status, created_at, updated_at FROM tickets WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var t models.Ticket
	err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return &t, nil
}

func (r *TicketRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE tickets SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTicketNotFound
	}

	return nil
}

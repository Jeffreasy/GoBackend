// service.go (contact directory)

// De Contact service slaat ontvangen contactformulieren op in de database.

package contact

import (
	"database/sql"
	"dklbackendGolang/internal/models"
)

type Service interface {
	SaveContact(c *models.Contact) error
}

type contactService struct {
	db *sql.DB
}

func NewService(db *sql.DB) Service {
	return &contactService{db: db}
}

func (s *contactService) SaveContact(c *models.Contact) error {
	query := `INSERT INTO contacts (name, email, message) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, c.Name, c.Email, c.Message)
	return err
}

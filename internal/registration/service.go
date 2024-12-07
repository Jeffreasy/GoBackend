package registration

import (
	"database/sql"

	"github.com/Jeffreasy/GoBackend/internal/models"
)

type Service interface {
	CreateRegistration(reg *models.Registration) error
	GetRegistrations() ([]models.Registration, error)
}

type registrationService struct {
	db *sql.DB
}

func NewService(db *sql.DB) Service {
	return &registrationService{db: db}
}

func (s *registrationService) CreateRegistration(reg *models.Registration) error {
	query := `
        INSERT INTO registrations 
        (name, email, role, distance, needs_support, support_details, terms_accepted)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

	return s.db.QueryRow(
		query,
		reg.Name,
		reg.Email,
		reg.Role,
		reg.Distance,
		reg.NeedsSupport,
		reg.SupportDetails,
		reg.TermsAccepted,
	).Scan(&reg.ID, &reg.CreatedAt, &reg.UpdatedAt)
}

func (s *registrationService) GetRegistrations() ([]models.Registration, error) {
	query := `SELECT * FROM registrations ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []models.Registration
	for rows.Next() {
		var reg models.Registration
		err := rows.Scan(
			&reg.ID,
			&reg.Name,
			&reg.Email,
			&reg.Role,
			&reg.Distance,
			&reg.NeedsSupport,
			&reg.SupportDetails,
			&reg.TermsAccepted,
			&reg.CreatedAt,
			&reg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		registrations = append(registrations, reg)
	}
	return registrations, nil
}

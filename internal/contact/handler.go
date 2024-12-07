// handler.go (contact directory)

// De Contact handler verwerkt HTTP-verzoeken om contactberichten op te slaan en eventueel notificaties te versturen.

package contact

import (
	"encoding/json"
	"net/http"

	"dklbackendGolang/internal/email"
	"dklbackendGolang/internal/models"
	"dklbackendGolang/pkg/validator"
)

type Handler struct {
	service      Service
	validator    *validator.Validator
	emailService email.Service
}

func NewHandler(service Service, v *validator.Validator, emailService email.Service) *Handler {
	return &Handler{
		service:      service,
		validator:    v,
		emailService: emailService,
	}
}

func (h *Handler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Ongeldige invoer", http.StatusBadRequest)
		return
	}

	// Valideer input
	if err := h.validator.Validate(contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sla het contactbericht op in de DB
	if err := h.service.SaveContact(&contact); err != nil {
		http.Error(w, "Kon contact niet opslaan", http.StatusInternalServerError)
		return
	}

	// Stuur een notificatie email naar de beheerder of een vastgesteld e-mailadres.
	h.emailService.SendMail(h.emailService.FromEmail(), "Nieuw contactbericht", "Er is een nieuw contactbericht ontvangen.")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Contactbericht ontvangen"})
}

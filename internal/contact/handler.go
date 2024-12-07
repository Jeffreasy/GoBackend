package contact

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jeffreasy/GoBackend/internal/email"
	"github.com/Jeffreasy/GoBackend/internal/models"
	"github.com/Jeffreasy/GoBackend/pkg/validator"
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
		log.Printf("Error decoding contact: %v", err)
		http.Error(w, "Ongeldige invoer", http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(contact); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SaveContact(&contact); err != nil {
		log.Printf("Error saving contact: %v", err)
		http.Error(w, "Kon contact niet opslaan", http.StatusInternalServerError)
		return
	}

	if err := h.emailService.SendMail(h.emailService.FromEmail(), "Nieuw contactbericht", "Er is een nieuw contactbericht ontvangen."); err != nil {
		log.Printf("Error sending email: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Contactbericht ontvangen"})
}

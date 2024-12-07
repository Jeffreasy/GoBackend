package contact

import (
	"encoding/json"
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
		http.Error(w, "Ongeldige invoer", http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SaveContact(&contact); err != nil {
		http.Error(w, "Kon contact niet opslaan", http.StatusInternalServerError)
		return
	}

	h.emailService.SendMail(h.emailService.FromEmail(), "Nieuw contactbericht", "Er is een nieuw contactbericht ontvangen.")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Contactbericht ontvangen"})
}

package registration

import (
	"encoding/json"
	"net/http"

	"github.com/Jeffreasy/GoBackend/internal/models"
	"github.com/Jeffreasy/GoBackend/pkg/validator"
)

type Handler struct {
	service   Service
	validator *validator.Validator
}

func NewHandler(service Service, v *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: v,
	}
}

func (h *Handler) CreateRegistration(w http.ResponseWriter, r *http.Request) {
	var reg models.Registration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, "Ongeldige invoer", http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(reg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateRegistration(&reg); err != nil {
		http.Error(w, "Kon registratie niet opslaan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reg)
}

func (h *Handler) GetRegistrations(w http.ResponseWriter, r *http.Request) {
	registrations, err := h.service.GetRegistrations()
	if err != nil {
		http.Error(w, "Kon registraties niet ophalen", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(registrations)
}

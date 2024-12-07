// handler.go (auth directory)

// De Auth handler regelt de HTTP-logica voor registratie en inloggen.
// Hier wordt input gevalideerd en worden diensten aangeroepen die de eigenlijke logica uitvoeren.

package auth

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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ongeldige invoer", http.StatusBadRequest)
		return
	}

	// Valideer de user input
	if err := h.validator.Validate(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Registreer de gebruiker via de service
	if err := h.service.RegisterUser(&user); err != nil {
		http.Error(w, "Kon gebruiker niet registreren", http.StatusInternalServerError)
		return
	}

	// Stuur een bevestigingsmail
	h.emailService.SendMail(user.Email, "Registratie succesvol", "Bedankt voor je registratie!")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Gebruiker geregistreerd"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Ongeldige invoer", http.StatusBadRequest)
		return
	}

	// Vraag de service om een JWT token op basis van email en wachtwoord
	token, err := h.service.Authenticate(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Ongeldige inloggegevens", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

package auth

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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ongeldige invoer", http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterUser(&user); err != nil {
		http.Error(w, "Kon gebruiker niet registreren", http.StatusInternalServerError)
		return
	}

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

	token, err := h.service.Authenticate(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Ongeldige inloggegevens", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

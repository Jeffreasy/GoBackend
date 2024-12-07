// user.go (models directory)

// De User model vertegenwoordigt een gebruiker in de applicatie.
// We gebruiken tags om velden te valideren (bijv. "required", "email").

package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

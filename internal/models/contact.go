// contact.go (models directory)

// De Contact model vertegenwoordigt een contactbericht.
// Ook hier gebruiken we validatietags om invoer te controleren.

package models

type Contact struct {
	ID      int    `json:"id"`
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Message string `json:"message" validate:"required"`
}

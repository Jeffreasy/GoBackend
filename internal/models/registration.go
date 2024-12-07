package models

import "time"

type Registration struct {
	ID             int       `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Email          string    `json:"email" validate:"required,email"`
	Role           string    `json:"role" validate:"required,oneof=deelnemer begeleider vrijwilliger"`
	Distance       string    `json:"distance" validate:"required,oneof=2.5km 6km 10km 15km"`
	NeedsSupport   string    `json:"needs_support" validate:"required,oneof=ja nee anders"`
	SupportDetails string    `json:"support_details,omitempty"`
	TermsAccepted  bool      `json:"terms_accepted" validate:"required,eq=true"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

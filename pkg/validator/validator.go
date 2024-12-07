// validator.go (pkg/validator)

// De validator maakt gebruik van de "go-playground/validator" package om struct velden te controleren op basis van tags.
// Op deze manier kan bijv. worden gecontroleerd of een email veld werkelijk een email is.

package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

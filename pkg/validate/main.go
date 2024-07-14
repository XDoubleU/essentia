// Package validate provides an easy way to validate the data contained in structs,
// typically being DTOs containing user input.
package validate

// ValidatedType is implemented by any struct with a Validate method.
type ValidatedType interface {
	Validate() *Validator
}

// Validator is used to validate contents
// of structs using [Check].
type Validator struct {
	Errors map[string]string
}

// New creates a new [Validator].
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid checks if a [Validator] has any errors.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) addError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check checks if value passes the validatorFunc.
// The provided key is used for creating the errors map of the [Validator].
func Check[T any](v *Validator, value T, validatorFunc ValidatorFunc[T], key string) {
	if result, message := validatorFunc(value); !result {
		v.addError(key, message)
	}
}

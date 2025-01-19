// Package validate provides an easy way to validate the data contained in structs,
// typically being DTOs containing user input.
package validate

// ValidatedType is implemented by any struct with a Validate method.
type ValidatedType interface {
	Validate() (bool, map[string]string)
}

// Validator is used to validate contents
// of structs using [Check].
type Validator struct {
	errors map[string]string
}

// New creates a new [Validator].
func New() *Validator {
	return &Validator{errors: make(map[string]string)}
}

// Valid checks if a [Validator] has any errors.
func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() map[string]string {
	return v.errors
}

func (v *Validator) addError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

// Check checks if value passes the validatorFunc.
// The provided key is used for creating the errors map of the [Validator].
func Check[T any](v *Validator, key string, value T, validatorFunc ValidatorFunc[T]) {
	if result, message := validatorFunc(value); !result {
		v.addError(key, message)
	}
}

func CheckOptional[T any](v *Validator, key string, value *T, validatorFunc ValidatorFunc[T]) {
	if value == nil {
		return
	}

	Check(v, key, *value, validatorFunc)
}

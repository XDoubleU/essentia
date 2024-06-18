package validate

type IValidatedType interface {
	Validate() *Validator
}

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func Check[T any](v *Validator, value T, validatorFunc ValidatorFunc[T], key string) {
	if result, message := validatorFunc(value); !result {
		v.AddError(key, message)
	}
}

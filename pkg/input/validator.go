package input

import "github.com/XDoubleU/essentia/pkg/router"

type Validator struct {
	rules  []Rule
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (validator *Validator) AddRule(
	key string,
	paramType ParamType,
	validation Validation,
	message string,
) {
	rule := Rule{
		Key:        key,
		ParamType:  paramType,
		Validation: validation,
		Message:    message,
	}
	validator.rules = append(validator.rules, rule)
}

func (validator *Validator) Validate(context *router.Context) bool {
	valid := true
	for _, rule := range validator.rules {
		if !rule.Check(context) {
			valid = false
			validator.addError(rule.Key, rule.Message)
		}
	}

	return valid
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) addError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

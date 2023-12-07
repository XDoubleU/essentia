package validator

import "github.com/XDoubleU/essentia/pkg/router"

type ParamType int

const (
	Query ParamType = 0
	Path  ParamType = 1
)

type Rule struct {
	Key        string
	ParamType  ParamType
	Validation Validation
	Message    string
}

func (rule Rule) Check(context *router.Context) bool {
	var value any

	switch rule.ParamType {
	case Query:
		value = context.GetQueryValue(rule.Key)
	case Path:
		value = context.GetPathValue(rule.Key)
	default:
		return false
	}

	return rule.Validation(value)
}

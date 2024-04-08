package parser

import "github.com/XDoubleU/essentia/pkg/router"

type ValueType int
type Validation[T any] func(value T) bool

const (
	Query ValueType = 0
	Path  ValueType = 1
	Body  ValueType = 2
)

type rule interface {
	Parse(*router.Context) *Error
}

type Rule[T any] struct {
	Key        string
	ValueType  ValueType
	Validation Validation[T]
	Message    string
}

type Error struct {
	Key     string
	Message string
}

func NewRule[T any](
	key string,
	valueType ValueType,
	validation Validation[T],
	message string,
) Rule[T] {
	return Rule[T]{
		Key:        key,
		ValueType:  valueType,
		Validation: validation,
		Message:    message,
	}
}

func (rule Rule[T]) Parse(c *router.Context) *Error {
	var value any

	// acquire value
	switch rule.ValueType {
	case Query:
		value = c.GetQueryValue(rule.Key)
	case Path:
		value = c.GetPathValue(rule.Key)
	case Body:
		value = c.GetBodyValue(rule.Key)
	default:
		return &Error{
			Key:     rule.Key,
			Message: "TODO",
		}
	}

	//TODO: check if value is required and if is present

	// conversion/parsing
	parsedValue := value.(T)

	// validation
	valid := rule.Validation(parsedValue)

	if !valid {
		return &Error{
			Key:     rule.Key,
			Message: rule.Message,
		}
	}

	// add to context
	switch rule.ValueType {
	case Query:
		c.QueryValues[rule.Key] = parsedValue
	case Path:
		c.PathValues[rule.Key] = parsedValue
	case Body:
		c.BodyValues[rule.Key] = parsedValue
	default:
		return &Error{
			Key:     rule.Key,
			Message: "TODO",
		}
	}

	return nil
}

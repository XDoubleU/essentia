package parser

/*
type Parser struct {
	rules  []rule
	Errors map[string]string
}

func NewParser(rules ...rule) *Parser {
	return &Parser{
		rules:  rules,
		Errors: make(map[string]string),
	}
}

func (p *Parser) Parse(context *router.Context) bool {
	valid := true
	for _, rule := range p.rules {
		if err := rule.Parse(context); err != nil {
			valid = false
			p.addError(*err)
		}
	}

	return valid
}

func (p Parser) Valid() bool {
	return len(p.Errors) == 0
}

func (p *Parser) addError(err Error) {
	if _, exists := p.Errors[err.Key]; !exists {
		p.Errors[err.Key] = err.Message
	}
}
*/

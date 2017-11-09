package rule

type RuleError struct {
	Message string
}

func (ruleError RuleError) Error() string {
	return ruleError.Message
}

package service

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type IdempotencyCheckerSvc interface {
	Execute(prev, curr int) bool
}

type checker struct{}

func NewChecker() IdempotencyCheckerSvc {
	return &checker{}
}

// Execute returns true if tgbotapi.Message.MessageID mismatched
func (checker) Execute(prev, curr int) bool {
	return prev != curr
}

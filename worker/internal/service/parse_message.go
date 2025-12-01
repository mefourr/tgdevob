package service

import (
	"encoding/json"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/domain"
)

var ErrNoData = errors.New("no data in consumed message")

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserRequestParserSvc interface {
	Execute(data []byte) (*domain.Message, error)
}

type parser struct {
}

func NewParserSvc() UserRequestParserSvc {
	return &parser{}
}

func (parser) Execute(data []byte) (*domain.Message, error) {
	if len(data) == 0 {
		return nil, ErrNoData
	}
	var res domain.Message
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

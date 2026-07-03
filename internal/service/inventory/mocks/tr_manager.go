package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type TrManager struct {
	mock.Mock
}

func (m *TrManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

package cache

import "github.com/stretchr/testify/mock"

type hashKeyToIntConverterMock[K comparable] struct {
	mock.Mock
}

// hashKeyToInt implements mocked function hashKeyToInt
func (m *hashKeyToIntConverterMock[K]) hashKeyToInt(key K) int {
	args := m.Called(key)
	return args.Int(0)
}

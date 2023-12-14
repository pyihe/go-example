package testify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCal(t *testing.T) {
	sum, err := Cal("+", 1, 2)
	assert.Equal(t, 3, sum, "sum should be 3")
	assert.Nil(t, err, "err should be nil")

	div, err := Cal("/", 1, 0)
	assert.Equal(t, 0, div, "divide should be zero")
	assert.EqualError(t, err, "divide by zero")
}

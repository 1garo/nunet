package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		filename []string
		e        error
	}{
		{[]string{".env", ".env.test"}, errors.New("cannot pass more than 1 filename")},
		{[]string{".env.not.found"}, errors.New("no .env file found")},
		{[]string{".env.empty"}, errors.New("TARGET_CONTAINER_ADDR environment variable not set")},
		{[]string{"../.env"}, nil},
	}

	for _, tt := range testCases {
		err := NewConfig(tt.filename...)

		assert.Equal(t, tt.e, err)
	}
}

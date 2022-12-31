package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	want := Config{}
	got, _ := New("/test")
	assert.Equal(t, want, got)
}

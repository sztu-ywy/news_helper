package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildUserKey(t *testing.T) {
	assert.Equal(t, "user:1", BuildUserKey(1))
	assert.Equal(t, "user:2:info", BuildUserKey(2, "info"))
	assert.Equal(t, "user:3:info:update_time", BuildUserKey(3, "info", "update_time"))
}

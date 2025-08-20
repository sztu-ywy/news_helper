package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSubjects(t *testing.T) {
	readonly := GetSubjects()
	assert.Equal(t, readonly, subjects)
	// if I change readonly value, it should not affect subjects,
	// says if I change readonly[0] = "test", subjects[0] should not be "test"
	// and if I call GetSubjects() again, it should return the original value of subjects
	readonly[0] = "test"
	readonly = GetSubjects()
	assert.Equal(t, readonly, subjects)
}

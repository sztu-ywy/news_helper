package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToMap(t *testing.T) {
	p := PermissionList{
		{
			Subject: User,
			Action:  Write,
		},
	}

	m := p.ToMap()
	assert.Equal(t, true, m[User][Write])
}

package acl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCan(t *testing.T) {
	p := PermissionList{
		{
			Subject: User,
			Action:  Write,
		},
	}

	m := p.ToMap()

	assert.Equal(t, true, Can(m, User, Write))

	p2 := PermissionList{
		{
			Subject: All,
			Action:  Write,
		},
		{
			Subject: All,
			Action:  Read,
		},
	}

	m2 := p2.ToMap()

	assert.Equal(t, true, Can(m2, All, Write))
	assert.Equal(t, true, Can(m2, All, Read))
	assert.Equal(t, true, Can(m2, User, Read))
	assert.Equal(t, true, Can(m2, User, Write))
}

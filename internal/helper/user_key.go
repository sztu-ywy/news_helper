package helper

import (
	"github.com/spf13/cast"
	"strings"
)

func BuildUserKey(userId uint64, groups ...string) string {
	var sb strings.Builder
	sb.WriteString("user:")
	sb.WriteString(cast.ToString(userId))
	for _, group := range groups {
		sb.WriteString(":")
		sb.WriteString(group)
	}
	return sb.String()
}

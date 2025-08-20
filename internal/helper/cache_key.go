package helper

import "strings"

func BuildKey(prefix string, id string) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteByte(':')
	sb.WriteString(id)
	return sb.String()
}

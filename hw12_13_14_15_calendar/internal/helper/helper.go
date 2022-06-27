package helper

import "strings"

func ConCat(s ...string) string {
	var builder strings.Builder
	var lgt int

	for _, v := range s {
		lgt += len(v)
	}

	builder.Grow(lgt)

	for _, v := range s {
		builder.WriteString(v)
	}

	return builder.String()
}

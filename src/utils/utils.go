package utils

import (
	"strings"
)

func EscapeForLarkMd(source string) string {
	target := strings.Replace(source, "\n", "\\n", -1)
	target = strings.Replace(target, "\r", "\\r", -1)
	target = strings.Replace(target, "\t", "\\t", -1)
	return target
}

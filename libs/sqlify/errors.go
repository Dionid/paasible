package sqlify

import "strings"

func NotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no rows in result")
}

func UniqueConstraintFailed(err error, postfix string) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed: "+postfix)
}

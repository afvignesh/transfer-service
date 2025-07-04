package middleware

import (
	"github.com/lib/pq"
)


func IsUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return true
	}
	return false
} 
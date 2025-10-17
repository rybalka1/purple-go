package utils

import (
	"net/http"
	"strings"
)

func ValidateFields(fields map[string]string, required []string, w http.ResponseWriter) bool {
	for _, field := range required {
		if val, ok := fields[field]; !ok || strings.TrimSpace(val) == "" {
			Error(w, http.StatusBadRequest, "Missing field: "+field)
			return false
		}
	}
	return true
}

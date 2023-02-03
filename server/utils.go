package server

import (
	"net/http"
	"strings"
)

func GetParamFromRequestURL(urlPath string) string {
	i := strings.LastIndexByte(urlPath, '/')
	return (urlPath[i+1:])
}

func IsJsonContentType(r *http.Request) bool {
	headerContentType := r.Header.Get("Content-Type")

	return headerContentType == "application/json"
}

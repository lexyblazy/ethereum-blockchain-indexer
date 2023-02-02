package server

import (
	"strings"
)

func GetAccountFromRequestURL(urlPath string) string {
	i := strings.LastIndexByte(urlPath, '/')
	return (urlPath[i+1:])
}

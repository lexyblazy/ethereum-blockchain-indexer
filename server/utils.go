package server

import (
	"strings"
)

func GetParamFromRequestURL(urlPath string) string {
	i := strings.LastIndexByte(urlPath, '/')
	return (urlPath[i+1:])
}

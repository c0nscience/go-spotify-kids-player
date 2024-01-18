package spotify

import (
	"net/url"
	"strings"
)

func GetIdFrom(urlString string) string {
	u, _ := url.Parse(urlString)
	p := strings.Split(u.Path[1:], "/")
	return p[len(p)-1]
}

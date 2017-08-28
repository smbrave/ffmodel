package main

import (
	"strings"
)

var LintGonicMapper = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func snakeString(s string) string {
	tmp := make([]byte, 0)
	data := make([]string, 0)
	var pre, cur bool = false, false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if d >= 'A' && d <= 'Z' {
			cur = true
		} else {
			cur = false
		}
		if i == 0 {
			tmp = append(tmp, d)
			pre = cur
			continue
		}
		if cur && !pre {
			data = append(data, strings.ToLower(string(tmp)))
			tmp = make([]byte, 0)
			tmp = append(tmp, d)
		} else {
			tmp = append(tmp, d)
		}
		pre = cur
	}
	data = append(data, strings.ToLower(string(tmp)))
	return strings.Join(data, "_")
}

// camel string, xx_yy to XxYy
func camelString(s string) string {
	tmp := strings.Split(s, "_")
	var data string
	for _, o := range tmp {
		if len(o) == 0 {
			continue
		}
		if _, ok := LintGonicMapper[strings.ToUpper(o)]; ok {
			data += strings.ToUpper(o)
			continue
		}
		data += strings.ToUpper(string(o[0]))
		if len(o) > 1 {
			data += strings.ToLower(o[1:])
		}
	}
	return data
}

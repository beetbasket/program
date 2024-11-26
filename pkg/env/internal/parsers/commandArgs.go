package parsers

import (
	"regexp"
	"strings"
)

type Args struct{}

var argsPattern = regexp.MustCompile(`("[^"]*"|[^"\s]+)(\s+|$)`)

func (Args) Name() string { return "args" }

func (Args) Parse(s string) ([]string, error) {
	ss := argsPattern.FindAllString(s, -1)
	out1 := make([]string, len(ss))
	for i, v := range ss {
		out1[i] = strings.TrimSpace(v)
	}
	return out1, nil
}

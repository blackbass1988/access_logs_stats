package core

import (
	"regexp"
	"strings"
	"log"
)

var regularExpressionRex = regexp.MustCompile(`[\[\]{}+*\\()]`)

type Filter struct {
	Matcher *NativeMatcher `json:"filter"`
	Prefix  string         `json:"prefix"`
	Items   []struct {
		Field   string   `json:"field"`
		Metrics []string `json:"metrics"`
	} `json:"items"`
}

func (f *Filter) MatchString(str string) bool {
	return f.Matcher.MatchString(str)
}

type Matcher interface {
	MatchString(str string) bool
	String() string
}

type NativeMatcher struct {
	Matcher
	raw       string
	isRegex   bool
	filterRex *regexp.Regexp
}

func (m *NativeMatcher) MatchString(str string) bool {

	//micro optimization
	if m.String() == ".+" || m.String() == ".*" {
		return true
	}

	if m.isRegex {
		return m.filterRex.MatchString(str)
	} else {
		return strings.Contains(str, m.raw)
	}
}

func (m *NativeMatcher) String() string {
	return m.raw
}

func newNativeMatcher(str string) (NativeMatcher, error) {
	var err error
	m := NativeMatcher{}
	m.raw = str

	if regularExpressionRex.MatchString(str) {
		m.isRegex = true
		log.Printf("filter [%s] was recognized as regular expersion\n", str)
		m.filterRex, err = regexp.Compile(str)
	} else {
		log.Printf("filter [%s] was recognized as regular string\n", str)
	}
	return m, err
}

func (m *NativeMatcher) UnmarshalJSON(data []byte) (err error) {
	*m, err = newNativeMatcher(string(data[1 : len(data)-1]))
	return err
}

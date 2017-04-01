package core

import (
	"log"
	"regexp"
	"strings"
)

var regularExpressionRex = regexp.MustCompile(`[\[\]{}+*\\()]`)

//Filter matching input string
type Filter struct {
	Matcher *nativeMatcher `json:"filter" yaml:"filter"`
	Prefix  string         `json:"prefix" yaml:"prefix"`
	Items   []struct {
		Field   string   `json:"field" yaml:"field"`
		Metrics []string `json:"metrics" yaml:"metrics"`
	} `json:"items" yaml:"items"`
}

//MatchString matches a input string
func (f *Filter) MatchString(str string) bool {
	return f.Matcher.MatchString(str)
}

type matcher interface {
	MatchString(str string) bool
	String() string
}

type nativeMatcher struct {
	matcher
	raw       string
	isRegex   bool
	filterRex *regexp.Regexp
}

func (m *nativeMatcher) MatchString(str string) bool {

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

func (m *nativeMatcher) String() string {
	return m.raw
}

func newNativeMatcher(str string) (nativeMatcher, error) {
	var err error
	m := nativeMatcher{}
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

func (m *nativeMatcher) UnmarshalJSON(data []byte) (err error) {
	*m, err = newNativeMatcher(string(data[1 : len(data)-1]))
	return err
}

func (m *nativeMatcher) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	v := ""
	if err := unmarshal(&v); err != nil {
		return err
	}
	*m, err = newNativeMatcher(v)
	return err
}

//func (m *nativeMatcher) UnmarshalYAML(data []byte) (err error) {
//	return err
//}

package re

import "regexp"

//NativeRegExp implements RegExp interface with bult-in regexp implementation
type NativeRegExp struct {
	RegExp
	re *regexp.Regexp
}

func newNativeRexCompile(expr string) (RegExp, error) {
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &NativeRegExp{re: r}, nil
}

// FindStringSubmatch returns a slice of matched strings
func (n *NativeRegExp) FindStringSubmatch(s string) []string {
	return n.re.FindStringSubmatch(s)
}

// SubexpNames returns the names of the parenthesized subexpressions
func (n *NativeRegExp) SubexpNames() []string {
	return n.re.SubexpNames()
}

// MatchString reports whether the Regexp matches the string s.
func (n *NativeRegExp) MatchString(s string) bool {
	return n.re.MatchString(s)
}

func (n *NativeRegExp) String() string {
	return n.re.String()
}

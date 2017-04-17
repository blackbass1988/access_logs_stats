package re

import (
	"errors"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"regexp"
)

//LibPcreRegexp implements RegExp interface with libcre
type LibPcreRegexp struct {
	RegExp
	re *pcre.Regexp

	currentMatcher *pcre.Matcher
	expr           string
	namedGroups    []string
}

func newLibPcreRegexp(expr string) (RegExp, error) {
	r, err := pcre.Compile(expr, 0)
	if err != nil {
		return nil, errors.New(err.String())
	}

	return &LibPcreRegexp{re: &r, expr: expr, namedGroups: getNamedGroupsFromExpression(expr)}, nil
}

// FindStringSubmatch returns a slice of matched strings
func (n *LibPcreRegexp) FindStringSubmatch(s string) []string {
	n.currentMatcher = n.re.MatcherString(s, 0)
	m := n.currentMatcher

	if !m.Matches() {
		return []string{}
	}

	groupsCnt := m.Groups() + 1

	matches := make([]string, groupsCnt)

	for i := 0; i < groupsCnt; i++ {
		matches[i] = m.GroupString(i)
	}

	return matches
}

// SubexpNames returns the names of the parenthesized subexpressions
func (n *LibPcreRegexp) SubexpNames() []string {
	return n.namedGroups
}

// MatchString reports whether the Regexp matches the string s.
func (n *LibPcreRegexp) MatchString(s string) bool {
	return n.re.MatcherString(s, 0).Matches()
}

func (n *LibPcreRegexp) String() string {
	return n.expr
}

func getNamedGroupsFromExpression(expr string) []string {
	//collect named groups from expression.
	r, err := regexp.Compile(`\?P<([^>]+)>`)
	if err != nil {
		panic(err)
	}

	matches := r.FindAllStringSubmatch(expr, -1)

	namedGroups := make([]string, len(matches)+1)

	for i, m := range matches {
		namedGroups[i+1] = m[1]
	}

	return namedGroups
}

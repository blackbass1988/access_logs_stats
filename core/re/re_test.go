package re

import (
	"reflect"
	"runtime"
	"testing"
)

//TestNewRegexp tests Compile that regexp.Regexp and libpcre implementations returns identical values
func TestNewRegexp(t *testing.T) {

	var check = func(fun func(s string) (RegExp, error), expr string, mustCompile bool) {
		m, err := fun(expr)
		if mustCompile {
			if err != nil {
				t.Errorf("matcher: %+v\nerr : %+v", getFunctionName(fun), err)
			} else if expr != m.String() {
				t.Errorf("matcher: %+v\nexpr != m.String() ; %+v != %+v",
					getFunctionName(fun), expr, m.String())
			}
		} else {
			if err == nil {
				t.Errorf("matcher: %+v\nmust failed but error is nil; expr: %s",
					getFunctionName(fun), expr)
			}
		}
	}


	type testCase struct {
		expr                string
		mustCompile                bool
	}

	var testCases = []testCase{
		{"foo", true},
		{".+", true},
		{"[a-z]+", true},
		{"foo|bar", true},
		{"zoo$", true},
		{"^zoo", true},
		{"[A-Z]+", true},
		{`(?P<code>\d+)`, true},

		//negatives
		{"(", false},
		{"((", false},
		{"\\", false},
		{"abc\\", false},

	}


	for _, tCase := range testCases {
		check(newLibPcreRegexp, tCase.expr, tCase.mustCompile)
		check(newNativeRexCompile, tCase.expr, tCase.mustCompile)
	}
}

//TestMatchString tests MatchString that regexp.Regexp and libpcre implementations returns identical values
func TestMatchString(t *testing.T) {

	var check = func(fun func(s string) (RegExp, error), expr string, s string, expectedMatchString bool) {
		m, _ := fun(expr)

		actualMatchString := m.MatchString(s)

		if actualMatchString != expectedMatchString {
			t.Errorf(
				"matcher: %+v\nactualMatchString != expectedMatchString; %+v != %+v",
				getFunctionName(fun), actualMatchString, expectedMatchString)
		}
	}

	type testCase struct {
		expr                string
		s                   string
		expectedMatchString bool
	}

	var testCases = []testCase{
		{"foo", "foo", true},
		{".+", "foo", true},
		{"[a-z]+", "omn omn zoo", true},
		{"foo|bar", "foo", true},
		{"zoo$", "omn omn zoo", true},
		{"^zoo", "omn omn zoo", false},
		{"[A-Z]+", "omn omn zoo", false},
	}

	for _, tCase := range testCases {
		check(newLibPcreRegexp, tCase.expr, tCase.s, tCase.expectedMatchString)
		check(newNativeRexCompile, tCase.expr, tCase.s, tCase.expectedMatchString)
	}

}
//TestFindStringSubmatch tests FindStringSubmatch that regexp.Regexp and libpcre implementations returns identical values
func TestFindStringSubmatch(t *testing.T) {

	var check = func(fun func(s string) (RegExp, error), expr string, s string, expectedSubmatches []string) {
		m, _ := fun(expr)

		actualSubmatches := m.FindStringSubmatch(s)

		actualGroupsCnt := len(actualSubmatches)
		expectedGroupsCnt := len(expectedSubmatches)

		if expectedGroupsCnt != actualGroupsCnt {
			t.Errorf(
				"matcher: %+v\nexpectedGroupsCnt != actualGroupsCnt; %+v != %+v",
				getFunctionName(fun), expectedGroupsCnt, actualGroupsCnt)
			t.FailNow()
		}

		for i, j := range actualSubmatches {
			if j != expectedSubmatches[i] {
				t.Errorf(
					"matcher: %+v\nactualSubmatches[%d] != expectedSubmatches[%d]; %+v != %+v",
					getFunctionName(fun), i, i, actualSubmatches[i], expectedSubmatches[i])

			}
		}
	}
	type testCase struct {
		expr               string
		s                  string
		expectedSubmatches []string
	}

	var testCases = []testCase{
		{"foo", "foo", []string{"foo"}},
		{"(foo)", "foo", []string{"foo", "foo"}},
		{"(foo)(bar)", "foobar", []string{"foobar", "foo", "bar"}},
		{`([^_]+)(_)(\S+)$`, "foo_bar", []string{"foo_bar", "foo", "_", "bar"}},
	}

	for _, tCase := range testCases {
		check(newLibPcreRegexp, tCase.expr, tCase.s, tCase.expectedSubmatches)
		check(newNativeRexCompile, tCase.expr, tCase.s, tCase.expectedSubmatches)
	}
}
//TestSubexpNames tests SubexpNames that regexp.Regexp and libpcre implementations returns identical values
func TestSubexpNames(t *testing.T) {
	var check = func(fun func(s string) (RegExp, error), expr string, s string, expectedMap map[string]string) {
		m, _ := fun(expr)

		actualSubmatches := m.FindStringSubmatch(s)
		if len(actualSubmatches) == 0 {
			t.Errorf("%s\nactualSubmatches is 0", getFunctionName(fun))
			t.FailNow()
		}

		actualGroups := m.SubexpNames()

		actualMap := make(map[string]string)

		for i, name := range actualGroups {

			if len(name) > 0 {
				actualMap[name] = actualSubmatches[i]
			}
		}

		for expectedName, expectedValue := range expectedMap {


			if _, ok := actualMap[expectedName]; !ok {
				t.Errorf(
					"matcher: %+v\nactualMap[%+v] not set",
					getFunctionName(fun), expectedName)

			} else if actualMap[expectedName] != expectedValue {
				t.Errorf(
					"matcher: %+v\nactualMap[%+v] != expectedValue[%+v]; %+v != %+v",
					getFunctionName(fun), expectedName, expectedName,
					actualMap[expectedName], expectedMap[expectedName])
			}

		}
	}

	type testCase struct {
		expr               string
		s                  string
		expectedGroups map[string]string
	}

	var testCases = []testCase{
		{`(?P<A>\S+)`, "foo", map[string]string{"A":"foo"}},
		{`(?P<A>\S+)_(?P<B>\d{1,}\.\d{3})`,
			"foo_12.345", map[string]string{"A":"foo", "B": "12.345"}},

	}

	for _, tCase := range testCases {
		check(newLibPcreRegexp, tCase.expr, tCase.s, tCase.expectedGroups)
		check(newNativeRexCompile, tCase.expr, tCase.s, tCase.expectedGroups)
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

var inputString = `s.auto.drom.ru s.auto.drom.ru 217.118.78.99 - [2017-03-19T20:57:44+10:00] GET "/i24204/r/photos/249454/gen177_1119983.jpg" HTTP/1.1 200 9453 "http://www.drom.ru/" "Mozilla/5.0 (Linux; Android 4.4.4; SM-T116 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.95 Safari/537.36" 0.200 717 "-" "-" HIT "-/" bad75053zRkT2GgTtk25FGYmFoYnw0a7 -`
//very slow in libpcre and fast in regexp
var expression = `HTTP/\d.?\d?\s(?P<code>\d+)[^"]+"[^"]*" "[^"]*" (?P<time>\d{1,}\.\d{3})`
var iterations = 100000

var r1, r2 RegExp



func BenchmarkLibPcreRegexp_MatchString(b *testing.B) {

	if r1 == nil {
		r1, _ = newLibPcreRegexp(expression)
	}

	for i := iterations ; i > 0; i -- {
		r1.MatchString(inputString)
	}
}

func BenchmarkNativeRegExp_MatchString(b *testing.B) {

	if r2 == nil {
		r2, _ = newNativeRexCompile(expression)
	}

	for i := iterations ; i > 0; i -- {
		r2.MatchString(inputString)
	}
}

func BenchmarkLibPcreRegexp_FindStringSubmatch(b *testing.B) {

	if r1 == nil {
		r1, _ = newLibPcreRegexp(expression)
	}

	for i := iterations ; i > 0; i -- {
		r1.FindStringSubmatch(inputString)
	}
}

func BenchmarkNativeRegExp_FindStringSubmatch(b *testing.B) {

	if r2 == nil {
		r2, _ = newNativeRexCompile(expression)
	}

	for i := iterations ; i > 0; i -- {
		r2.MatchString(inputString)
	}
}
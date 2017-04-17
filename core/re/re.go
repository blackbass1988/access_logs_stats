package re

//RegExp matches regular expression and returns matched strings
type RegExp interface {
	MatchString(s string) bool
	FindStringSubmatch(s string) []string
	SubexpNames() []string
	String() string
}

//Compile returns implementation of RegExp
func Compile(expr string) (RegExp, error) {
	return newNativeRexCompile(expr)
}
func MustCompile(expr string) RegExp {
	r, err := newNativeRexCompile(expr)

	if err != nil {
		panic(err)
	}
	return r
}

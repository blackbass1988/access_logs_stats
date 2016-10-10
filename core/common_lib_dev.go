//+build devel

package core

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//+build devel

package pkg

func check(err error) {
	if err != nil {
		panic(err)
	}
}

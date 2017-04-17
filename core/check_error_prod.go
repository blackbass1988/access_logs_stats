//+build !devel

package core

import (
	"log"
)

func checkOrFail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

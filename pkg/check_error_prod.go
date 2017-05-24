//+build !devel

package pkg

import (
	"log"
)

func checkOrFail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

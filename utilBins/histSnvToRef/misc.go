package main

import (
	"fmt"
	"time"
)

func timeStamp() string {
	t := time.Now()
	return fmt.Sprintf("[%02v:%02v:%02v]", t.Hour(), t.Minute(), t.Second())
}

func setBool(a []bool, b bool) {
	for i := range a {
		a[i] = b
	}
}

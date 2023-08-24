package threading

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strconv"
)

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RoutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func RunSafe(fn func()) {
	defer func() {
		if p := recover(); p != nil {
			msg := fmt.Sprintf("%v\n%s", p, string(debug.Stack()))
			log.Println(msg)
		}
	}()
	fn()
}

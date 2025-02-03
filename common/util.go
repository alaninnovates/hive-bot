package common

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math/rand"
	"time"
)

func ImageToPipe(image image.Image) *io.PipeReader {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		if err := png.Encode(w, image); err != nil {
			fmt.Println(err)
		}
	}()
	return r
}

func ShuffleArray[T any](array []T) []T {
	dest := make([]T, len(array))
	perm := rand.Perm(len(array))
	for i, v := range perm {
		dest[v] = array[i]
	}
	return dest
}

func ArrayIncludes[T comparable](array []T, value T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

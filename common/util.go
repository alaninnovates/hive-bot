package common

import (
	"image"
	"image/png"
	"io"
)

func ImageToPipe(image image.Image) *io.PipeReader {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		if err := png.Encode(w, image); err != nil {
			panic(err)
		}
	}()
	return r
}

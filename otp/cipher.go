//go:build !solution

package otp

import (
	"io"
)

type streamCipherReader struct {
	r    io.Reader
	prng io.Reader
}

func (scr *streamCipherReader) Read(p []byte) (int, error) {
	n, err := scr.r.Read(p)
	if n > 0 {
		prngData := make([]byte, n)
		if _, prngErr := scr.prng.Read(prngData); prngErr != nil {
			return n, prngErr
		}
		for i := 0; i < n; i++ {
			p[i] ^= prngData[i]
		}
	}

	return n, err
}

type streamCipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func (scw *streamCipherWriter) Write(p []byte) (int, error) {
	prngData := make([]byte, len(p))
	if _, err := scw.prng.Read(prngData); err != nil {
		return 0, err
	}
	encoded := make([]byte, len(p))
	for i := 0; i < len(p); i++ {
		encoded[i] = p[i] ^ prngData[i]
	}

	return scw.w.Write(encoded)
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &streamCipherReader{r: r, prng: prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &streamCipherWriter{w: w, prng: prng}
}

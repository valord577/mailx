package mailx

import "io"

// @author valor.

// According to
//  1. RFC 2045, 6.7. (page 21) for quoted-printable
//  2. RFC 2045, 6.8. (page 25) for base64
//
// The encoded output stream must be represented in lines
// with no more than 76 characters per line.
const maxLineLength = 76

type multipartBase64Writer struct {
	w io.Writer
}

// Write implements io.Writer
func (w *multipartBase64Writer) Write(p []byte) (int, error) {
	length := len(p)
	if length <= maxLineLength {
		return w.write(p)
	}

	sum := 0

	doit := true
	index := 0
	for doit {
		if index >= length {
			break
		}

		end := index + maxLineLength
		if end >= length {
			end = length
			doit = false
		}
		n, err := w.write(p[index:end])
		if err != nil {
			return 0, err
		}

		sum += n
		index += maxLineLength
	}
	return sum, nil
}

func (w *multipartBase64Writer) write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	if err != nil {
		return 0, err
	}
	x, err := w.w.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	return n + x, nil
}

// Package byframe provides a simple framing protocol for streaming data.
package byframe

import (
	"errors"
)

const (
	maskContinue = 0b1000_0000
	maskfraction = 0b0111_1111
)

var (
	// ErrInsufficient ...
	ErrInsufficient = errors.New("[byframe] data is not sufficient to construct the body")

	// ErrHeaderInsufficient ...
	ErrHeaderInsufficient = errors.New("[byframe] data is not sufficient to construct the header")

	// ErrHeaderTooLarge ...
	ErrHeaderTooLarge = errors.New("[byframe] header is too long")
)

// EncodeHeader encode data length into header
func EncodeHeader(l int) (header []byte) {
	if l == 0 {
		return []byte{0}
	}

	for l > 0 {
		digit := l & maskfraction
		l >>= 7

		if l > 0 {
			header = append(header, byte(digit|128))
		} else {
			header = append(header, byte(digit))
		}
	}

	return
}

// HeaderLength decodes the length of the header.
// Returns zero if the raw is not sufficient.
// It will return -1 if the header length is too large.
func HeaderLength(raw []byte) int {
	for i, b := range raw {
		if i > 8 {
			return -1
		}
		if b&maskContinue == 0 {
			return i + 1
		}
	}
	return 0
}

// DecodeHeader decode bytes into header length and data length.
func DecodeHeader(raw []byte) (header int, data int) {
	header = HeaderLength(raw)
	for i := 0; i < header; i++ {
		digit := int(raw[i])
		data |= (digit & maskfraction) << (i * 7)
	}
	return
}

// Encode encode data into frame format
func Encode(data []byte) []byte {
	header := EncodeHeader(len(data))
	return append(header, data...)
}

// Decode decode frame into data, decoded bytes and error
func Decode(data []byte) ([]byte, error) {
	headerLen, dataLen := DecodeHeader(data)
	if headerLen > 0 {
		n := headerLen + int(dataLen)
		if len(data) < n {
			return nil, ErrInsufficient
		}
		return data[headerLen:n], nil
	}

	if headerLen == 0 {
		return nil, ErrHeaderInsufficient
	}

	return nil, ErrHeaderTooLarge
}

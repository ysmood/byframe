package byframe

import (
	"bytes"
	"encoding/gob"
	"errors"
)

// ErrInsufficient ...
var ErrInsufficient = errors.New("data is not sufficient to construct the body")

// ErrHeaderInsufficient ...
var ErrHeaderInsufficient = errors.New("data is not sufficient to construct the header")

// EncodeHeader encode data length into header
func EncodeHeader(l int) (header []byte) {
	if l == 0 {
		return []byte{0}
	}

	for l > 0 {
		digit := l & 127
		l >>= 7

		if l > 0 {
			header = append(header, byte(digit|128))
		} else {
			header = append(header, byte(digit))
		}
	}

	return
}

// DecodeHeader decode bytes into data length, header length and whether it's sufficient to
// parse the header from raw.
func DecodeHeader(raw []byte) (int, int, bool) {
	rawLen := len(raw)
	headerLen := 0
	dataLen := 0

	for {
		if headerLen >= rawLen {
			return headerLen, dataLen, false
		}
		digit := int(raw[headerLen])
		dataLen |= (digit & 127) << (uint(headerLen) * 7)
		headerLen++
		if (digit & 128) == 0 {
			break
		}
	}
	return dataLen, headerLen, true
}

// EncodeBytes encode data into frame format
func EncodeBytes(data []byte) []byte {
	header := EncodeHeader(len(data))
	return append(header, data...)
}

// DecodeBytes decode frame into data, decoded bytes and error
func DecodeBytes(data []byte) ([]byte, int, error) {
	dataLen, headerLen, sufficient := DecodeHeader(data)
	if !sufficient {
		return nil, 0, ErrHeaderInsufficient
	}
	n := headerLen + dataLen
	if len(data) < n {
		return nil, 0, ErrInsufficient
	}
	return data[headerLen:n], n, nil
}

// Encode arbitrary value
func Encode(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return EncodeBytes(buf.Bytes()), nil
}

// Decode bytes into arbitrary value
func Decode(data []byte, value interface{}) error {
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)

	frame, _, err := DecodeBytes(data)
	if err != nil {
		return err
	}

	_, _ = buf.Write(frame)

	return dec.Decode(value)
}

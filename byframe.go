package byframe

import (
	"errors"
)

// ErrInsufficient ...
var ErrInsufficient = errors.New("data is not sufficient to construct the body")

// ErrHeaderInsufficient ...
var ErrHeaderInsufficient = errors.New("data is not sufficient to construct the header")

// EncodeHeader encode data length into header
func EncodeHeader(l int) []byte {
	header := []byte{}

	for l > 0 {
		// 128 is 0b10000000
		digit := l % 128
		l = (l - digit) / 128

		if l > 0 {
			header = append(header, byte(digit|128))
		} else {
			header = append(header, byte(digit))
		}
	}

	return header
}

// DecodeHeader decode bytes into data length, header length and whether it's sufficient to
// parse the header from raw.
func DecodeHeader(raw []byte) (int, int, bool) {
	rawLen := len(raw)
	headerLen := 0
	dataLen := 0
	isContinue := true

	for isContinue {
		if headerLen == rawLen {
			return headerLen, dataLen, false
		}
		digit := int(raw[headerLen])
		isContinue = (digit & 128) == 128                   // 128 is 0b10000000
		dataLen += (digit & 127) * (1 << uint(headerLen*7)) // 127 is 0b01111111
		headerLen++
	}
	return dataLen, headerLen, true
}

// Encode encode data into frame format
func Encode(data []byte) []byte {
	header := EncodeHeader(len(data))
	return append(header, data...)
}

// Decode decode frame into data
func Decode(data []byte) ([]byte, error) {
	dataLen, headerLen, sufficient := DecodeHeader(data)
	if !sufficient {
		return nil, ErrHeaderInsufficient
	}
	if len(data) < headerLen+dataLen {
		return nil, ErrInsufficient
	}
	return data[headerLen : headerLen+dataLen], nil
}

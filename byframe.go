package byframe

import (
	"errors"
	"io"
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

// Scanner scan frames based on the length header
type Scanner struct {
	r          io.Reader
	frame      []byte
	dataLen    int
	headerLen  int
	sufficient bool
	buf        []byte
	err        error
}

// NewScanner just like line scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:          r,
		frame:      []byte{},
		dataLen:    0,
		headerLen:  0,
		sufficient: false,
		buf:        []byte{},
	}
}

func (s *Scanner) read() bool {
	buf := make([]byte, 1024)
	n, err := s.r.Read(buf)
	s.buf = append(s.buf, buf[:n]...)
	if err != nil {
		s.err = err
		return false
	}
	return true
}

// Scan scan next frame
func (s *Scanner) Scan() bool {
	for {
		dl, hl, sufficient := DecodeHeader(s.buf[s.headerLen:])
		s.dataLen += dl
		s.headerLen += hl
		s.sufficient = sufficient

		if sufficient {
			break
		}

		if !s.read() {
			return false
		}
	}

	for {
		if len(s.buf) >= s.headerLen+s.dataLen {
			s.frame = s.buf[s.headerLen : s.headerLen+s.dataLen]

			// reset
			s.buf = s.buf[s.headerLen+s.dataLen:]
			s.dataLen = 0
			s.headerLen = 0
			s.sufficient = false

			return true
		}

		if !s.read() {
			return false
		}
	}
}

// Frame current frame
func (s *Scanner) Frame() []byte {
	return s.frame
}

// Err the error
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

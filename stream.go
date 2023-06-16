package byframe

import (
	"fmt"
	"io"
)

// ErrLimitExceeded ...
var ErrLimitExceeded = fmt.Errorf("[byframe] exceeded the limit")

// Scanner scan frames based on the length header
type Scanner struct {
	limit int
	r     io.Reader
	frame []byte
	err   error
}

// NewScanner just like line scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		limit: 10 * 1024 * 1024,
		r:     r,
	}
}

// Limit of frame size, default is 10MB
func (s *Scanner) Limit(size int) *Scanner {
	s.limit = size
	return s
}

// Next scan next frame, returns error if errors happen while scanning
func (s *Scanner) Next() ([]byte, error) {
	if !s.Scan() {
		return nil, s.Err()
	}
	return s.Frame(), nil
}

// Scan scan next frame, returns true to continue the scan
func (s *Scanner) Scan() bool {
	var buf []byte
	b := make([]byte, 1)
	headerLen := 0
	dataLen := 0
	for headerLen == 0 {
		_, s.err = s.r.Read(b)
		if s.err != nil {
			return false
		}

		buf = append(buf, b[0])

		headerLen, dataLen = DecodeHeader(buf)

		if headerLen < 0 {
			s.err = ErrHeaderTooLarge
			return false
		}
	}

	if dataLen > s.limit {
		s.err = ErrLimitExceeded
		return false
	}

	s.frame = make([]byte, dataLen)

	_, s.err = io.ReadFull(s.r, s.frame)
	return s.err == nil
}

// Frame returns the current frame
func (s *Scanner) Frame() []byte {
	return s.frame
}

// Err returns error if errors happen while scanning
func (s *Scanner) Err() error {
	return s.err
}

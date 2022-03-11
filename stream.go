package byframe

import (
	"fmt"
	"io"
)

// ErrLimitExceeded ...
var ErrLimitExceeded = fmt.Errorf("[byframe] exceeded the limit")

// Scanner scan frames based on the length header
type Scanner struct {
	limit   int
	r       io.Reader
	frame   []byte
	buf     []byte
	readBuf []byte
	err     error
}

// NewScanner just like line scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		limit:   1024 * 1024 * 1024,
		r:       r,
		buf:     []byte{},
		readBuf: make([]byte, 64*1024),
	}
}

// Limit of frame size, default is 1GB
func (s *Scanner) Limit(size int) *Scanner {
	s.limit = size
	return s
}

// BufferSize of the buffer for read, default is 64KB
func (s *Scanner) BufferSize(size int) *Scanner {
	s.readBuf = make([]byte, size)
	return s
}

// Scan scan next frame, returns true to continue the scan
func (s *Scanner) Scan() bool {
	headerDone := false
	headerLen := 0
	dataLen := 0

	for {
		if len(s.buf) > s.limit {
			s.err = ErrLimitExceeded
			return false
		}

		if !headerDone {
			headerLen, dataLen = DecodeHeader(s.buf)
			if headerLen > 0 {
				headerDone = true
			} else if headerLen < 0 {
				s.err = ErrHeaderTooLarge
				return false
			}
		}

		if headerDone && len(s.buf) >= headerLen+dataLen {
			s.frame = s.buf[headerLen : headerLen+dataLen]
			s.buf = s.buf[headerLen+dataLen:]
			return true
		}

		n, err := s.r.Read(s.readBuf)
		if err != nil {
			s.err = err
			return false
		}
		s.buf = append(s.buf, s.readBuf[:n]...)
	}
}

// Frame returns the current frame
func (s *Scanner) Frame() []byte {
	return s.frame
}

// Err returns error if errors happen while scanning
func (s *Scanner) Err() error {
	return s.err
}

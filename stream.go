package byframe

import "io"

// Scanner scan frames based on the length header
type Scanner struct {
	limit int

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

// Limit of frame buffer, panic if exceeds. Zero means no limit.
func (s *Scanner) Limit(size int) *Scanner {
	s.limit = size
	return s
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
		s.checkLimit()
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
		s.checkLimit()
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

func (s *Scanner) checkLimit() {
	if s.limit != 0 && len(s.buf) > s.limit {
		panic("[byframe] exceeded the limit")
	}
}

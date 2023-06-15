package byframe_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/ysmood/byframe/v3"
)

func ExampleScanner() {
	var buf bytes.Buffer

	for i := 0; i < 3; i++ {
		frame := byframe.Encode([]byte(fmt.Sprintf("%d", i)))
		buf.Write(frame)
	}

	s := byframe.NewScanner(&buf)

	for s.Scan() {
		fmt.Println(string(s.Frame()))
	}

	// Output:
	// 0
	// 1
	// 2
}

func (t T) Scanner() {
	frame := byframe.Encode([]byte("test data"))
	r := bytes.NewReader(frame)
	s := byframe.NewScanner(r)

	for s.Scan() {
		t.Eq([]byte("test data"), s.Frame())
	}
}

func (t T) ScannerNext() {
	frame := byframe.Encode([]byte("test data"))
	r := bytes.NewReader(frame)
	s := byframe.NewScanner(r)

	b, err := s.Next()
	t.E(err)
	t.Eq([]byte("test data"), b)
}

func (t T) ScannerMultiFrames() {
	frameA := byframe.Encode([]byte("test a"))
	frameB := byframe.Encode([]byte("test b"))
	r := bytes.NewReader(append(frameA, frameB...))
	s := byframe.NewScanner(r)

	list := [][]byte{}
	for s.Scan() {
		list = append(list, s.Frame())
	}
	t.Eq([][]byte{[]byte("test a"), []byte("test b")}, list)
}

func (t T) ScannerOptions() {
	frame := byframe.Encode([]byte("test data test data"))
	r := bytes.NewReader(frame)
	s := byframe.NewScanner(r).Limit(10).BufferSize(1)

	for s.Scan() {
		nop()
	}
	t.Eq(s.Err(), byframe.ErrLimitExceeded)
}

func (t T) ScannerLargeHeaderErr() {
	r := bytes.NewReader(bytes.Repeat([]byte{0b1000_0000}, 20))
	s := byframe.NewScanner(r)

	for s.Scan() {
		nop()
	}
	t.Eq(s.Err(), byframe.ErrHeaderTooLarge)
}

type MultiRead struct {
	i            int
	returnedZero int
	frame        []byte
}

// read only one byte each time
func (mr *MultiRead) Read(buf []byte) (int, error) {
	// simulate (0, nil) edge case
	if mr.i == mr.returnedZero {
		mr.i++
		return 0, nil
	}

	copy(buf, mr.frame[mr.i:mr.i+1])
	mr.i++
	if mr.i == len(mr.frame) {
		return 0, io.EOF
	}
	return 1, nil
}

func (t T) ScannerMultiRead() {
	data := []byte(strings.Repeat("test data", 100))
	frame := byframe.Encode(data)

	s := byframe.NewScanner(&MultiRead{frame: frame, returnedZero: 5})

	for s.Scan() {
		t.Eq(data, s.Frame())
	}
}

type ErrRead struct {
}

func (mr *ErrRead) Read([]byte) (int, error) {
	return 0, errors.New("err")
}

func (t T) ScannerReadErr() {
	s := byframe.NewScanner(&ErrRead{})

	hit := false
	for s.Scan() {
		hit = true
	}
	t.False(hit)
	t.Eq(errors.New("err"), s.Err())
}

func (t T) StreamMonkey() {
	list := [][]byte{}
	buf := bytes.NewBuffer(nil)

	for i := 0; i < 1000; i++ {
		data := bytes.Repeat([]byte{1}, rand.Intn(128*1024))
		buf.Write(byframe.Encode(data))
		list = append(list, data)
	}

	s := byframe.NewScanner(buf)

	for s.Scan() {
		item := list[0]
		list = list[1:]
		t.True(bytes.Equal(s.Frame(), item))
	}
}

func nop() {}

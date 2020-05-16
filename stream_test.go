package byframe_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/byframe/v2"
)

func ExampleScanner() {
	type User struct {
		Name string
		Age  int
	}

	var buf bytes.Buffer

	frame, _ := byframe.Encode(User{"Ann", 10})
	buf.Write(frame)

	frame, _ = byframe.Encode(User{"Lee", 12})
	buf.Write(frame)

	s := byframe.NewScanner(&buf)

	for s.Scan() {
		var user User
		_ = s.Decode(&user)
		fmt.Println(user)
	}

	// Output:
	// {Ann 10}
	// {Lee 12}
}

func TestScanner(t *testing.T) {
	frame := byframe.EncodeBytes([]byte("test data"))
	r := bytes.NewReader(frame)
	s := byframe.NewScanner(r)

	for s.Scan() {
		assert.Equal(t, []byte("test data"), s.Frame())
	}
	assert.Nil(t, s.Err())
}

func TestScannerMultiFrames(t *testing.T) {
	frameA := byframe.EncodeBytes([]byte("test a"))
	frameB := byframe.EncodeBytes([]byte("test b"))
	r := bytes.NewReader(append(frameA, frameB...))
	s := byframe.NewScanner(r)

	list := [][]byte{}
	for s.Scan() {
		list = append(list, s.Frame())
	}
	assert.Equal(t, [][]byte{[]byte("test a"), []byte("test b")}, list)
	assert.Nil(t, s.Err())
}

func TestScannerLimit(t *testing.T) {
	frame := byframe.EncodeBytes([]byte("test data test data"))
	r := bytes.NewReader(frame)
	s := byframe.NewScanner(r).Limit(10)

	assert.Panics(t, func() {
		s.Scan()
	})
}

type MultiRead struct {
	i            int
	returnedZero bool
	frame        []byte
}

// read only one byte each time
func (mr *MultiRead) Read(buf []byte) (int, error) {
	// simulate (0, nil) edge case
	if !mr.returnedZero && mr.i == 5 {
		mr.returnedZero = true
		return 0, nil
	}

	copy(buf, mr.frame[mr.i:mr.i+1])
	mr.i++
	if mr.i == len(mr.frame) {
		return 0, io.EOF
	}
	return 1, nil
}

func TestScannerMultiRead(t *testing.T) {
	frame := byframe.EncodeBytes([]byte("test data"))

	s := byframe.NewScanner(&MultiRead{i: 0, frame: frame})

	for s.Scan() {
		assert.Equal(t, []byte("test data"), s.Frame())
	}
	assert.Nil(t, s.Err())
}

type ErrRead struct {
}

func (mr *ErrRead) Read(buf []byte) (int, error) {
	return 0, errors.New("err")
}

func TestScannerReadErr(t *testing.T) {
	s := byframe.NewScanner(&ErrRead{})

	hit := false
	for s.Scan() {
		hit = true
	}
	assert.False(t, hit)
	assert.Equal(t, errors.New("err"), s.Err())
}

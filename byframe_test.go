package byframe_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/byframe"
)

func ExampleEncodeHeader() {
	header := byframe.EncodeHeader(1000)

	dataLen, headerLen, sufficient := byframe.DecodeHeader(header)

	fmt.Println(headerLen, dataLen, sufficient)
	// Output: 2 1000 true
}

func ExampleEncode() {
	frame := byframe.Encode([]byte("test"))

	data, _ := byframe.Decode(frame)

	fmt.Println(string(data))
	// Output: test
}

func ExampleScanner() {
	frame := byframe.Encode([]byte("test"))
	s := byframe.NewScanner(bytes.NewReader(frame))

	for s.Scan() {
		fmt.Println(string(s.Frame()))
		// Output: test
	}
}

func TestEncodeHeader20(t *testing.T) {
	h := byframe.EncodeHeader(20)
	assert.Equal(t, []byte{20}, h)
}

func TestEncodeHeader200(t *testing.T) {
	h := byframe.EncodeHeader(200)
	assert.Equal(t, []byte{0xc8, 0x1}, h)
}

func TestDecode(t *testing.T) {
	frame := byframe.Encode([]byte("test data"))
	data, err := byframe.Decode(frame)
	assert.Nil(t, err)
	assert.Equal(t, []byte("test data"), data)
}

func TestDecodeErrHeaderInsufficient(t *testing.T) {
	frame := byframe.Encode([]byte{})
	_, err := byframe.Decode(frame)
	assert.Equal(t, byframe.ErrHeaderInsufficient, err)
}

func TestDecodeErrInsufficient(t *testing.T) {
	_, err := byframe.Decode([]byte{10})
	assert.Equal(t, byframe.ErrInsufficient, err)
}

func TestHeaderInsufficient(t *testing.T) {
	h, l, s := byframe.DecodeHeader([]byte{135})
	assert.Equal(t, 1, h)
	assert.Equal(t, 7, l)
	assert.Equal(t, false, s)
}

func TestScanner(t *testing.T) {
	frame := byframe.Encode([]byte("test data"))
	r := bytes.NewReader(frame)
	s := byframe.NewScanner(r)

	for s.Scan() {
		assert.Equal(t, []byte("test data"), s.Frame())
	}
	assert.Nil(t, s.Err())
}

func TestScannerMultiFrames(t *testing.T) {
	frame := byframe.Encode([]byte("test data"))
	r := bytes.NewReader(append(frame, frame...))
	s := byframe.NewScanner(r)

	list := [][]byte{}
	for s.Scan() {
		list = append(list, s.Frame())
	}
	assert.Equal(t, [][]byte{[]byte("test data"), []byte("test data")}, list)
	assert.Nil(t, s.Err())
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
	frame := byframe.Encode([]byte("test data"))

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

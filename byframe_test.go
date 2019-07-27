package byframe_test

import (
	"bytes"
	"fmt"
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
	}
	// Output: test
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

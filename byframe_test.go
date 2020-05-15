package byframe_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/byframe/v2"
)

func ExampleEncodeHeader() {
	header := byframe.EncodeHeader(1000)

	dataLen, headerLen, sufficient := byframe.DecodeHeader(header)

	fmt.Println(headerLen, dataLen, sufficient)
	// Output: 2 1000 true
}

func ExampleEncode() {
	type User struct {
		Name string
		Age  int
	}

	data, _ := byframe.Encode(User{"Ann", 10})

	var decoded User
	_ = byframe.Decode(data, &decoded)

	fmt.Println(decoded)
	// Output: {Ann 10}
}

func ExampleEncodeBytes() {
	frame := byframe.EncodeBytes([]byte("test"))

	data, _, _ := byframe.DecodeBytes(frame)

	fmt.Println(string(data))
	// Output: test
}

func TestEncodeHeader0(t *testing.T) {
	h := byframe.EncodeHeader(0)
	assert.Equal(t, []byte{0}, h)
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
	frame := byframe.EncodeBytes([]byte("test data"))
	data, _, err := byframe.DecodeBytes(frame)
	assert.Nil(t, err)
	assert.Equal(t, []byte("test data"), data)
}

func TestDecodeErrHeaderInsufficient(t *testing.T) {
	_, _, err := byframe.DecodeBytes([]byte{128})
	assert.Equal(t, byframe.ErrHeaderInsufficient, err)
}

func TestDecodeErrInsufficient(t *testing.T) {
	_, _, err := byframe.DecodeBytes([]byte{10})
	assert.Equal(t, byframe.ErrInsufficient, err)
}

func TestHeaderInsufficient(t *testing.T) {
	h, l, s := byframe.DecodeHeader([]byte{135})
	assert.Equal(t, 1, h)
	assert.Equal(t, 7, l)
	assert.Equal(t, false, s)
}

func TestErr(t *testing.T) {
	_, err := byframe.Encode(nil)
	assert.Error(t, err)

	err = byframe.Decode(nil, nil)
	assert.Error(t, err)

	err = byframe.Decode(byframe.EncodeBytes([]byte{}), nil)
	assert.Error(t, err)
}

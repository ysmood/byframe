package byframe_test

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/ysmood/byframe/v3"
	"github.com/ysmood/got"
)

type T struct {
	got.G
}

func Test(t *testing.T) {
	got.Each(t, T{})
}

func ExampleEncodeHeader() {
	header := byframe.EncodeHeader(1000)

	headerLen, dataLen := byframe.DecodeHeader(header)

	fmt.Println(headerLen, dataLen)

	// Output: 2 1000
}

func ExampleEncode() {
	data := byframe.Encode([]byte("test"))

	decoded, _ := byframe.Decode(data)

	fmt.Println(string(decoded))

	// Output: test
}

func (t T) EncodeHeader0() {
	h := byframe.EncodeHeader(0)
	t.Eq([]byte{0}, h)
}

func (t T) EncodeHeader20() {
	h := byframe.EncodeHeader(20)
	t.Eq([]byte{20}, h)
}

func (t T) EncodeHeader200() {
	h := byframe.EncodeHeader(200)
	t.Eq([]byte{0xc8, 0x1}, h)
}

func (t T) EncodeHeaderHuge() {
	h := byframe.EncodeHeader(math.MaxInt)
	t.Eq(len(h), 9)
}

func (t T) Decode() {
	frame := byframe.Encode([]byte("test data"))
	data, err := byframe.Decode(frame)
	t.E(err)
	t.Eq([]byte("test data"), data)
}

func (t T) DecodeWithExtraData() {
	frame := byframe.Encode([]byte("test data"))
	data, err := byframe.Decode(append(frame, []byte("test")...))
	t.E(err)
	t.Eq([]byte("test data"), data)
}

func (t T) DecodeErrHeaderInsufficient() {
	_, err := byframe.Decode([]byte{0b1000_0000})
	t.Eq(byframe.ErrHeaderInsufficient, err)
}

func (t T) DecodeHeaderMax() {
	data := bytes.Repeat([]byte{0b1111_1111}, 8)
	data = append(data, 0b0111_1111)
	h, d := byframe.DecodeHeader(data)
	t.Eq(h, 9)
	t.Eq(d, math.MaxInt)
}

func (t T) DecodeErrHeaderTooLarge() {
	data := bytes.Repeat([]byte{0b1000_0000}, 20)
	_, err := byframe.Decode(data)
	t.Eq(err, byframe.ErrHeaderTooLarge)
}

func (t T) DecodeErrInsufficient() {
	_, err := byframe.Decode([]byte{10})
	t.Eq(byframe.ErrInsufficient, err)
}

func (t T) EncodeDecodeNil() {
	byframe.Encode(nil)

	_, err := byframe.Decode(nil)
	t.Err(err)
}

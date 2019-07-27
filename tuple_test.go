package byframe_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/byframe"
)

func ExampleEncodeTuple() {
	type Name struct {
		First []byte
		Last  []byte
	}

	name := Name{[]byte("Jack"), []byte("Black")}

	data := byframe.EncodeTuple(&name.First, &name.Last)

	var newName Name
	_ = byframe.DecodeTuple(data, &newName.First, &newName.Last)

	fmt.Println(string(newName.First), string(newName.Last))
	// Output: Jack Black
}

func TestEncodeTuple(t *testing.T) {
	type Name struct {
		First []byte
		Last  []byte
	}

	name := Name{[]byte("a"), []byte("b")}

	data := byframe.EncodeTuple(&name.First, &name.Last)
	assert.Len(t, data, 3)

	var newName Name
	assert.Nil(t, byframe.DecodeTuple(data, &newName.First, &newName.Last))
	assert.Equal(t, name, newName)
}

func TestEncodeZeroLengthTuple(t *testing.T) {
	data := byframe.EncodeTuple()
	assert.Nil(t, byframe.DecodeTuple(data))
}

func TestDecodeTupleErr(t *testing.T) {
	assert.Nil(t, byframe.DecodeTuple([]byte{}))

	d := []byte("")
	err := byframe.DecodeTuple([]byte{1}, &d, &d)
	assert.Equal(t, byframe.ErrInsufficient, err)
}

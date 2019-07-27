package byframe

// EncodeTuple encode an array of []byte into an []byte.
// The format:
//
// | itemA length | itemA bytes | ... | itemN length | itemN bytes | itemLast bytes |
//
// The last item doesn't need length header. The protection for data curruption can be done
// on other levels of data handling, this algorithm is designed for small size.
func EncodeTuple(items ...*[]byte) []byte {
	itemsLen := len(items)

	last := itemsLen - 1
	if itemsLen == 0 {
		return []byte{0}
	}

	data := []byte{}

	for _, item := range items[:last] {
		data = append(data, Encode(*item)...)
	}
	data = append(data, *items[last]...)

	return data
}

// DecodeTuple zero copy decode
func DecodeTuple(data []byte, items ...*[]byte) error {
	itemsLen := len(items)

	if itemsLen == 0 {
		return nil
	}

	last := itemsLen - 1

	for i := 0; i < last; i++ {
		item, n, err := Decode(data)
		if err != nil {
			return err
		}
		*items[i] = item
		data = data[n:]
	}
	*items[last] = data
	return nil
}

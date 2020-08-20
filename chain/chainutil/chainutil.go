package chainutil

// ReverseBytes receives a byte slice and returns a byte slice after toggling
// the endianness
func ReverseBytes(inputBytes []byte) []byte {
	size := len(inputBytes)
	for i := 0; i < size/2; i++ {
		inputBytes[i], inputBytes[size-1-i] = inputBytes[size-1-i], inputBytes[i]
	}
	return inputBytes
}

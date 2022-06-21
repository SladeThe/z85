package z85

import (
	"encoding/binary"
	"math"
)

const (
	encoding = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.-:+=^!/*?&<>()[]{}@%$#"
)

var (
	decoding = make([]byte, math.MaxUint8)
)

func init() {
	for i := range decoding {
		decoding[i] = math.MaxUint8
	}

	for i, r := range encoding {
		decoding[r] = byte(i)
	}
}

// EncodeToString uses EncodeTo under the hood and returns the result as a string.
func EncodeToString(plain []byte) string {
	return string(Encode(plain))
}

// DecodeString a string encoded with EncodeToString, returning plain bytes.
func DecodeString(encoded string) ([]byte, error) {
	return Decode([]byte(encoded))
}

// Encode uses EncodeTo under the hood, taking care of allocating the output buffer.
func Encode(plain []byte) []byte {
	encoded := make([]byte, EncodedLen(plain))
	_, _ = EncodeTo(plain, encoded)
	return encoded
}

// Decode uses DecodeTo under the hood, taking care of allocating the output buffer.
func Decode(encoded []byte) ([]byte, error) {
	decodedCap, errCap := DecodedCap(encoded)
	if errCap != nil || decodedCap == 0 {
		return nil, errCap
	}

	decoded := make([]byte, decodedCap)

	n, errDecode := DecodeTo(encoded, decoded)
	if errDecode != nil {
		return nil, errDecode
	}

	return decoded[:n], nil
}

// EncodedLen returns the exact length of the buffer needed to hold all the bytes of the Z85 encoded data.
func EncodedLen(plain []byte) int {
	n := len(plain)
	if n == 0 {
		return 0
	}

	return (n + 4) / 4 * 5
}

// DecodedCap returns the buffer length that is certainly enough to hold all the bytes of the Z85 decoded data.
// The len(encoded) must be divisible by 5, otherwise InvalidEncodedLengthError is returned.
// The actual number of bytes decoded by DecodeTo can be less than the value returned by DecodedCap.
func DecodedCap(encoded []byte) (int, error) {
	n := len(encoded)
	if n == 0 {
		return 0, nil
	}
	if n%5 != 0 {
		return 0, InvalidEncodedLengthError(n)
	}

	return n * 4 / 5, nil
}

// EncodeTo encodes source plain bytes into EncodedLen(plain) bytes of the encoded buffer using Z85 encoding.
// The number of bytes written is always EncodedLen(plain).
// If len(encoded) is less than EncodedLen(plain), InsufficientDestinationLengthError is returned.
func EncodeTo(plain, encoded []byte) (int, error) {
	encodedLen := EncodedLen(plain)

	if encodedLen == 0 {
		return 0, nil
	}

	if len(encoded) < encodedLen {
		return 0, InsufficientDestinationLengthError{want: encodedLen, got: len(encoded)}
	}

	for len(plain) >= 4 {
		v := binary.BigEndian.Uint32(plain)

		encoded[4] = encoding[v%85]
		v /= 85
		encoded[3] = encoding[v%85]
		v /= 85
		encoded[2] = encoding[v%85]
		v /= 85
		encoded[1] = encoding[v%85]
		v /= 85
		encoded[0] = encoding[v%85]

		plain = plain[4:]
		encoded = encoded[5:]
	}

	var v uint32

	switch len(plain) {
	case 0:
		v = uint32(0) | uint32(0)<<8 | uint32(0)<<16 | uint32(1)<<24
	case 1:
		v = uint32(0) | uint32(0)<<8 | uint32(1)<<16 | uint32(plain[0])<<24
	case 2:
		v = uint32(0) | uint32(1)<<8 | uint32(plain[1])<<16 | uint32(plain[0])<<24
	case 3:
		v = uint32(1) | uint32(plain[2])<<8 | uint32(plain[1])<<16 | uint32(plain[0])<<24
	}

	for i := 4; i >= 0; i-- {
		encoded[i] = encoding[v%85]
		v /= 85
	}

	return encodedLen, nil
}

// DecodeTo decodes Z85 encoded source into the target plain buffer, returning the actual number of decoded bytes.
// The length of the buffer must be at least DecodedCap(encoded).
// The len(encoded) must be divisible by 5, otherwise InvalidEncodedLengthError is returned.
// If DecodeTo encounters invalid input bytes, it returns InvalidEncodedByteError.
// Incorrectly padded encoded source leads to InvalidPostfixError.
func DecodeTo(encoded, plain []byte) (int, error) {
	decodedCap, errCap := DecodedCap(encoded)

	if errCap != nil {
		return 0, errCap
	}

	if len(encoded) == 0 {
		return 0, nil
	}

	if len(plain) < decodedCap {
		return 0, InsufficientDestinationLengthError{want: decodedCap, got: len(plain)}
	}

	decoded := plain

	for len(encoded) > 0 {
		r0 := encoded[0]
		r1 := encoded[1]
		r2 := encoded[2]
		r3 := encoded[3]
		r4 := encoded[4]

		m0 := uint32(decoding[r0])
		m1 := uint32(decoding[r1])
		m2 := uint32(decoding[r2])
		m3 := uint32(decoding[r3])
		m4 := uint32(decoding[r4])

		if m0|m1|m2|m3|m4 == math.MaxUint8 {
			if m0 == math.MaxUint8 {
				return 0, InvalidEncodedByteError(r0)
			} else if m1 == math.MaxUint8 {
				return 0, InvalidEncodedByteError(r1)
			} else if m2 == math.MaxUint8 {
				return 0, InvalidEncodedByteError(r2)
			} else if m3 == math.MaxUint8 {
				return 0, InvalidEncodedByteError(r3)
			} else {
				return 0, InvalidEncodedByteError(r4)
			}
		}

		binary.BigEndian.PutUint32(decoded, m0*52200625+m1*614125+m2*7225+m3*85+m4)

		encoded = encoded[5:]
		decoded = decoded[4:]
	}

	for i := decodedCap - 1; i >= 0; i-- {
		switch plain[i] {
		case 1:
			return i, nil
		case 0:
			continue
		default:
			return 0, InvalidPostfixError(plain[i])
		}
	}

	return 0, InvalidPostfixError(0)
}

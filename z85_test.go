package z85

import (
	"math/rand"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	benchmarkPlain   = make([]byte, 1<<20)
	benchmarkEncoded string
)

func init() {
	rnd := rand.New(rand.NewSource(567537357543))
	rnd.Read(benchmarkPlain)
	benchmarkEncoded = EncodeToString(benchmarkPlain)
}

func TestEncodeString(t *testing.T) {
	require.Equal(t, "", EncodeToString(nil), "unexpected encoded nil text")
	require.Equal(t, "", EncodeToString([]byte{}), "unexpected encoded zero length text")
}

func TestDecodeString(t *testing.T) {
	{
		decoded, errDecode := DecodeString("")
		require.NoError(t, errDecode, "unexpected decode error")
		require.Equal(t, []byte(nil), decoded, "unexpected decoded empty text")
	}

	{
		_, errDecode := DecodeString("1")
		require.Error(t, errDecode, "expected decode error")
	}
}

func TestDecodeTo(t *testing.T) {
	plain := []byte{1}
	encoded := make([]byte, EncodedLen(plain))

	_, errEncode := EncodeTo(plain, encoded)
	require.NoError(t, errEncode, "unexpected encode error")

	decodedCap, errCap := DecodedCap(encoded)
	require.NoError(t, errCap, "unexpected decoded cap error")

	decoded := make([]byte, decodedCap+4)
	decoded[len(decoded)-1] = 1
	decodedLen, errDecode := DecodeTo(encoded, decoded)
	require.NoError(t, errDecode, "unexpected decode error")
	decoded = decoded[:decodedLen]

	require.Equal(t, plain, decoded, "unexpected decoded text")
}

func TestEncodeDecodeString(t *testing.T) {
	rnd := rand.New(rand.NewSource(567537357543))

	for i := 0; i < 1000; i++ {
		plain := make([]byte, 1+rnd.Intn(1000))
		rnd.Read(plain)

		n := (len(plain) + 4) / 4 * 5

		encoded := EncodeToString(plain)
		require.Equal(t, n, len(encoded), "unexpected encoded length")

		decoded, errDecode := DecodeString(encoded)
		require.NoError(t, errDecode, "unexpected decode error")
		require.Equal(t, plain, decoded, "unexpected decoded text")
	}
}

func BenchmarkEncodeToString(b *testing.B) {
	var encoded string

	for i := 0; i < b.N; i++ {
		encoded = EncodeToString(benchmarkPlain)
	}

	runtime.KeepAlive(encoded)
}

func BenchmarkDecodeString(b *testing.B) {
	var decoded []byte
	var err error

	for i := 0; i < b.N; i++ {
		decoded, err = DecodeString(benchmarkEncoded)
	}

	runtime.KeepAlive(decoded)
	runtime.KeepAlive(err)
}

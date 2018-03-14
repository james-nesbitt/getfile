package testassist

import (
	"math/rand"
	"net/http"
	"testing"
)

const TESTSOCKET = ":8000"
const DUMMYLENGTH = 1000

var ds []byte

func dummyBytes() []byte {
	if len(ds) > 0 {
		ds = make([]byte, DUMMYLENGTH)
		rand.Read(ds)
	}
	return ds
}

// TestServerResponse handles an incoming http request, passing back a testable byte sequence
func dummyServerResponse(w http.ResponseWriter, r *http.Request) {
	w.Write(dummyBytes())
}

// Start a Dummy server, which will deliver a testable byte stream
func DummyServerStart(t *testing.T) {
	http.HandleFunc("/", dummyServerResponse)
	go func() {
		http.ListenAndServe(":8000", nil)
	}()
}

// Compare 2 byte streams
func CompareBytes(bs []byte) bool {
	cbs := dummyBytes()
	for i, b := range bs {
		if cbs[i] != b {
			return false
		}
	}
	return true

	// We don't use this as we want to be able to test partial Byte slices
	//return bytes.Compare(b, DummyBytes()) == 0
}

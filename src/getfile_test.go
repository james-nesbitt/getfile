package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"testing"
	"context"
)

/**
 * Here we control what source is used for testing.
 *
 * Originally we tried to use an http.Server to serve a controlled bytestream, but it didn't properly handle the
 * Content Length header, nor the range header.
 *
 * In this version, we compare what contents we retrieve to a manually calculated MD5 hash of the first 4MiB. It is not
 * the ideal test, as we had no perfect test of the first 4MiB.
 *
 */

const URL = "http://4ca5b8f6.bwtest-aws.pravala.com/384MB.jar" // File source
const MD54MB = "778a5540263ef2df98016cef6ce3a94c"              // MD5 hash of 4MiB of the source (manually created)

const COUNT = 4       // 4 pieces
const SIZE = 1.049e+6 // 1MiB

// run some rudimentary tests on our GetFile methods
func TestGetFile(t *testing.T) {
	gf := NewGetFile(URL)

	if gfm, err := gf.Meta(); err != nil {
		t.Error("Error retreiving meta for url: ", err.Error())
	} else {
		l := gfm.Length

		t.Log("Length: ", l)

		if pr, err := gf.Pieces(context.Background(), COUNT, SIZE); err != nil {
			t.Error("failed to get file pieces: ", err.Error())
		} else if bs, err := ioutil.ReadAll(pr); err != nil {
			t.Error("Failed to read bytes: ", err.Error())
		} else if len(bs) == 0 {
			t.Error("no pieces retreived: ", pr)
		} else if fmt.Sprintf("%x", md5.Sum(bs)) != MD54MB {
			t.Error("Incorrect MD5 of results: " + fmt.Sprintf("%x", md5.Sum(bs)))
		}
	}
}

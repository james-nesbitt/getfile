package main

import (
	"io"
	"strconv"
	"testing"
)

/**
 * As FilePieces is considered one our our risky structs, we emphasize it for testing.
 *
 * As a struct it needs to handle concurrency, and should be able fault and memory tolerant.
 */

// Basic testing of FilePieces functionality - add local bytes, and retrieve them
func TestFilePieces(t *testing.T) {
	p := NewFilePieces(10)

	if pbr, err := p.All(); err == nil {
		t.Error("Empty FilePieces did not give an error when retrieving")
	} else {
		pb := make([]byte, 10)
		if n, err := pbr.Read(pb); err != io.EOF {
			t.Error("FilePieces reader did not give an EOF error when reading: ", err.Error())
		} else if n != 0 {
			t.Error("Empty FilePieces gave the wrong number of bytes")
		}
	}

	for i := 0; i < 10; i++ {
		if err := p.Set(i, []byte{byte(i)}); err != nil {
			t.Error("Couldn't set a byte array for index:", strconv.Itoa(i)+" :", err.Error())
		}
	}

	if pbr, err := p.All(); err != nil {
		t.Error("Error retrieving byte reader :", err.Error())
	} else {

		pb := make([]byte, 10)
		if n, err := pbr.Read(pb); err != nil {
			t.Error("FilePieces returned an error when reading the bytes: ", err.Error())
		} else if n != 10 {
			t.Error("FilePieces returned the wrong number of bytes: ", strconv.Itoa(n))
		} else {
			for i, b := range pb {
				eb := byte(i)
				if b != eb {
					t.Error("FilePieces piece byte didn't match expected value")
				}
			}
		}

	}
}

// @TODO do some FilePieces concurrency testing

// @TODO do some FilePieces testing for large pieces

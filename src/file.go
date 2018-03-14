package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

/**
 * file.go
 *
 * Here are some structs related to file data.  They are utility structs for containing metadata and piece data
 * for retrieved files.  The utilities don't in general require any remote connectivity.
 */

// FileMeta delivers Meta Data about a remote file, based on response Headers
type FileMeta struct {
	Length int
}

// NewFileMeta_FromHeaders build a FileMeta from an http response headers map
func NewFileMeta_FromHeaders(headers map[string][]string) (FileMeta, error) {
	var fm FileMeta

	if _, exists := headers["Content-Length"]; !exists {
		return fm, errors.New("no content-length provided in http headers: " + fmt.Sprint(headers))
	} else if l, err := strconv.Atoi(headers["Content-Length"][0]); err != nil {
		return fm, err
	} else {
		fm.Length = l
	}

	return fm, nil
}

/**
 * FilePieces store and retrieve pieces of a file download
 *
 * @Note we needed a custom struct to handle pieces of a file.
 *
 * There are many risks related to piecing a file:
 *   - piece failure may result in a missing piece,
 *   - pieces may be really big, and we may need advanced memory/storage management
 *   - pieces may be retrieved concurrently, so we need to be safe with assignments
 *
 * Because of this, we start with a premade tool for handling Pieces, which has a defined structure, and can be
 * debugged/improved on it's own.
 *
 * @TODO make this handle piece failure
 * @TODO make this consider memory usage with large pieces
 * @TODO this is not safe for concurrent execution
 * @TODO this should have a lock to prevent concurrent retrival/sets
 */
type FilePieces struct {
	num    int
	pieces map[int][]byte
}

// Construct a new FilePieces struct of expected length
func NewFilePieces(n int) *FilePieces {
	return &FilePieces{
		num:    n,
		pieces: map[int][]byte{},
	}
}

// Set a piece index value
// @TODO this needs to be safe for cuncurrent execution
func (fp *FilePieces) Set(i int, b []byte) error {
	if fp.pieces == nil {
		fp.pieces = map[int][]byte{}
	}

	fp.pieces[i] = b
	return nil
}

// Get all of the pieces concatenated byte reader
// @TODO this needs to be safe for cuncurrent execution
func (fp *FilePieces) All() (io.Reader, error) {
	if fp.pieces == nil {
		fp.pieces = map[int][]byte{}
	}

	// use a simple bytes.Buffer to concatenate all of the bytes.
	body := bytes.NewBuffer([]byte{})

	if len(fp.pieces) == 0 {
		return body, errors.New("no pieces have been added")
	}

	for i := 0; i < fp.num; i++ {
		if p, exists := fp.pieces[i]; exists {
			body.Write(p)
		} else {
			return body, errors.New("missing a piece : " + strconv.Itoa(i))
		}
	}
	return io.Reader(body), nil
}

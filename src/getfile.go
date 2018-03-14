package main

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"context"
)

// A Struct that can be used to retrieve files in parts
// @TODO this doesn't need to be public, as we have a public constructor
type GetFile struct {
	source string // target file path, expected to be a URL
}

// NewGetFile constructor for a GetFile
func NewGetFile(source string) *GetFile {
	return &GetFile{
		source: source,
	}
}

// GetMeta gets meta information about the file by retreiving Http headers
func (g *GetFile) Meta() (FileMeta, error) {
	res, _ := http.Head(g.source)
	return NewFileMeta_FromHeaders(res.Header)
}

// GetPieces
// @NOTE does not do any validation on filesize
func (g *GetFile) Pieces(ctx context.Context,num, size int) (io.Reader, error) {
	var wg sync.WaitGroup

	pieces := NewFilePieces(num) // a place to keep all of the downloaded pieces
	gErr := make(chan error)     // a global error signal handler

	// for each piece, run a parallel retrieval, trigger a global error if needed.
	for i := 0; i < num; i++ {
		wg.Add(1)

		min := size * i       // Min range
		max := size * (i + 1) // Max range

		go func(min int, max int, i int) {
			if p, err := g.Piece(min, max); err != nil {
				go func() {
					gErr <- err
				}()
			} else {
				pieces.Set(i, p)
			}
			wg.Done()
		}(min, max, i)
	}

	// use a finished chan to track when all of the parrallel requests are completed.
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	// now wait for done or an error/context end
	select {
	case <- ctx.Done():
		if err := ctx.Err(); err != nil {
			return nil, err
		} else {
			return nil, errors.New("limiting context ended")
		}
	case err := <-gErr:
		return nil, err
	case <-done:
		return pieces.All()
	}
}

/**
 * GetPiece retrieves a piece of a file download from the source
 *
 * I considered read into an io.Reader, but it provides no advantage, as the response Body needs closing
 * This is intended to be safe for parrallel usage, as each http.client is independent
 */
func (g *GetFile) Piece(min, max int) ([]byte, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", g.source, nil)
	// Add the data for the Range header of the form "bytes=0-100"
	rangeHeader := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
	req.Header.Add("Range", rangeHeader)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return []byte{}, err
	}
	if bs, err := ioutil.ReadAll(resp.Body); err != nil {
		return []byte{}, err
	} else {
		return bs, nil
	}
}

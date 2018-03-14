package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

/**
 * A POC for downloading a file in parts using the Range header.
 *
 * This POC is heavily based on:
 *    https://coderwall.com/p/uz2noa/fast-parallel-downloads-in-golang-with-accept-ranges-and-goroutines
 *
 * This POC demonstrates:
 *  - Concurrent part downloads
 *  - http client & range header
 * This POC does not:
 *  - allow partial downloads
 *  - have any reusable components
 *  - have any tests (not in focus for the POC)
 *  - would have pretty bad memory hits for large files (all downloaded pieces are kept in memory during collection
 */

var url string
var name string

var limit int
var coroutines int

var wg sync.WaitGroup

var wd string

func init() {
	url = "http://4ca5b8f6.bwtest-aws.pravala.com/384MB.jar"
	name = "384.jar"
	limit = -1 // stop after limit pieces, <1 means unlimited
	coroutines = 25
	wd, _ = os.Getwd()
}

func main() {
	log.Info("Starting: ", url)

	res, _ := http.Head(url)
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0]) // Get the content length from the header request

	len_sub := length / limit // Bytes for each Go-routine
	diff := length % limit    // Get the remaining for the last request

	log.Info("Sizes", length, limit, len_sub, diff)

	numRuns := coroutines
	if limit > 0 && limit < coroutines {
		numRuns = limit
	}

	body := make([]string, limit+1) // Make up a temporary array to hold the data to be written to the file
	for i := 0; i < numRuns; i++ {
		wg.Add(1)

		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range

		if i == limit-1 {
			max += diff // Add the remaining bytes in the last request
		}

		log.Info("Firing routing:", i)
		go func(min int, max int, i int) {
			log.Info("Starting routing:", i, min, max)
			client := &http.Client{}
			req, _ := http.NewRequest("GET", url, nil)
			range_header := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1) // Add the data for the Range header of the form "bytes=0-100"
			req.Header.Add("Range", range_header)
			resp, _ := client.Do(req)
			defer resp.Body.Close()
			reader, _ := ioutil.ReadAll(resp.Body)
			body[i] = string(reader)

			log.Info("Finished routing:", i)
			wg.Done()

		}(min, max, i)
	}

	log.Info("Waiting for processes to finish")
	wg.Wait()

	log.Info("Collecting file")
	temp, _ := ioutil.TempFile(wd, "getfile")
	for i := 0; i < limit; i++ {
		temp.Write([]byte(body[i]))
	}
	fullpath := filepath.Join(wd, name, ".new")
	log.Info("Writing file:", fullpath)
	os.Rename(temp.Name(), fullpath)

	log.Info("Stopping")
}

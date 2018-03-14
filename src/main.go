package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var url string  // url (http) to the source file
var name string // file name for the destination (later we choose a dir)

var limit int // how many pieces to download
var size int  // how big (in bytes) each piece should be

var wd string // what path to download the file to

const (
	TEST_URL  = "http://4ca5b8f6.bwtest-aws.pravala.com/384MB.jar"
	TEST_NAME = "384.jar" // file name for downloaded files
)

// Init the program
func init() {
	wd, _ = os.Getwd() // put any downloaded files into the pwd
}

// main exec function
func main() {
	/**
	 * use urfave/cli to set up the cli
	 */

	app := cli.NewApp()
	app.Name = "getfile"
	app.Usage = "Use this command to partially download a file over http, in pieces"
	app.Version = VERSION

	// @NOTE it would be nice if urfave/cli did something with this
	app.ArgsUsage = "{url} : http path to file to be downloaded"

	// use flags for the configuration, and attach them directly to the variables
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file, f",
			Value:       "gotfile.bin",
			Usage:       "local file path where the file will be saved",
			Destination: &name,
		},
		cli.IntFlag{
			Name:        "limit, l",
			Value:       4,
			Usage:       "Number of pieces to download",
			Destination: &limit,
		},
		cli.IntFlag{
			Name:        "size, s",
			Value:       1.049e+6,
			Usage:       "Size of downloaded pieces",
			Destination: &size,
		},
	}

	// define a single command for getfile execution, use a lamda to pipe flags into the real handler
	app.Action = func(c *cli.Context) error {
		url = c.Args().Get(0)
		if url == "" {
			return errors.New("Missing url.  You have to pass a valid url to the command.")
		}

		return getfile(url, path.Join(wd, name), limit, size)
	}

	// Run the cli app, catch any errors
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Local handoff that uses the GetFile struct to perform the download, with a fair amount of logging output
func getfile(url string, file string, limit int, size int) error {
	log.Info("Starting: ", url)

	gf := NewGetFile(url)

	if pr, err := gf.Pieces(limit, size); err != nil {
		log.WithError(err).Error("Could not retrieve file pieces")
		return err
	} else if f, err := ioutil.TempFile(wd, "getfile"); err != nil {
		log.WithError(err).Error("Could not create a TempFile target")
		return err
	} else if l, err := io.Copy(f, pr); err != nil || l == 0 {
		log.WithError(err).Error("Could not write bytes to the TempFile")
		return err
	} else {
		f.Close()
		if err := os.Rename(f.Name(), file); err != nil {
			log.WithError(err).Error("Could not write bytes to the TempFile")
			return err
		}
	}

	// test if we like the contents: make sure we made a file, and output an md5 for the user
	if f, err := os.Open(file); err != nil {
		log.WithError(err).Error("Couldn't open new file")
		return err
	} else if bs, err := ioutil.ReadAll(f); err != nil {
		log.WithError(err).Error("Couldn't Read bytes from file")
		return err
	} else {
		log.Info("File Created: " + file + " [MD5:" + fmt.Sprintf("%x", md5.Sum(bs)) + "]")
	}

	log.Info("Stopping")
	return nil
}

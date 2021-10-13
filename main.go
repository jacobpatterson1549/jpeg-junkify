// Package main provides a command-line-interface to reduce the quality of jpeg images.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jacobpatterson1549/jpeg-junkify/jpeg"
)

func main() {
	var filename, inDirectory, outDirectory, bytesLimit string
	flag.CommandLine = flagSet(&filename, &inDirectory, &outDirectory, &bytesLimit)
	flag.Parse()
	junkifyFiles(bytesLimit, filename, inDirectory, outDirectory)
}

func junkifyFiles(b, filename, inDirectory, outDirectory string) {
	bytesLimit, err := parseBytesLimit(b)
	switch {
	case err != nil:
		log.Fatal(fmt.Errorf("parsing max file size: %v", err))
	case bytesLimit <= 0:
		log.Fatal("positive file size limit required")
	case len(inDirectory) > 0:
		junkifyBulk(inDirectory, outDirectory, bytesLimit)
	case len(filename) == 0:
		log.Fatal("file name required")
	default:
		outFilename := filepath.Join(".", filename)
		junkifyFile(filename, outFilename, bytesLimit)
	}
}

func junkifyBulk(inDirectory, outDirectory string, bytesLimit int) {
	dirEntries, err := os.ReadDir(inDirectory)
	if err != nil {
		log.Fatal(fmt.Errorf("reading input file directory: %v", err))
	}
	errs := make(chan error)
	for _, e := range dirEntries {
		if e.IsDir() {
			continue
		}
		inFilename := filepath.Join(inDirectory, e.Name())
		outFilename := filepath.Join(outDirectory, e.Name())
		go func() {
			errs <- junkifyFile(inFilename, outFilename, bytesLimit)
		}()
	}
	for range dirEntries {
		if err := <-errs; err != nil {
			log.Print(err)
		}
	}
}

func junkifyFile(inFilename, outFilename string, bytesLimit int) error {
	if !isJPEG(inFilename) {
		return fmt.Errorf("%v is not a jpeg, skipping", inFilename)
	}
	f, err := os.Open(inFilename)
	if err != nil {
		return fmt.Errorf("opening input file: %v", err)
	}
	b, q, err := jpeg.Junkfiy(f, bytesLimit)
	if err != nil {
		return fmt.Errorf("error shrinking %v: %v", inFilename, err)
	}
	outFilename2 := fmt.Sprintf("%v_quality_%v.jpg", outFilename, q)
	w, err := os.OpenFile(outFilename2, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("opening output to save: %v", err)
	}
	defer w.Close()
	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("writing output file: %v", err)
	}
	fmt.Printf("processed %v to %v\n", inFilename, outFilename2)
	return nil
}

func flagSet(file, inDirectory, outDirectory, bytesLimit *string) *flag.FlagSet {
	programName := os.Args[0]
	fs := flag.NewFlagSet(programName, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Reduces the quality of JPEG image(s).")
		fs.PrintDefaults()
	}
	fs.StringVar(bytesLimit, "b", "1MB", "The maximum image size, in bytes.  Example suffixes: 1000, 1KB, 1K, 2M (the first three are equivalent)")
	fs.StringVar(file, "f", "", "The name of the jpeg file to junkify.")
	fs.StringVar(inDirectory, "in-dir", "", "The directory of jpeg images to junkify.  Takes precedence over the 'f' parameter.")
	fs.StringVar(outDirectory, "out-dir", ".", "The directory to write junkified images to.")
	return fs
}

func parseBytesLimit(s string) (int, error) {
	suffixMultipliers := map[string]int{
		"K":  1000,
		"KB": 1000,
		"M":  1000 * 1000,
		"MB": 1000 * 1000,
	}
	m := 1
	for suffix, multiplier := range suffixMultipliers {
		if strings.HasSuffix(s, suffix) {
			s = strings.TrimSuffix(s, suffix)
			m = multiplier
			break
		}
	}
	bytesLimit, err := strconv.Atoi(s)
	if err != nil {
		return bytesLimit, fmt.Errorf("parsing bytes limit value: %v", s)
	}
	bytesLimit *= m
	return bytesLimit, nil
}

func isJPEG(filename string) bool {
	extension := filepath.Ext(filename)
	extension = strings.ToLower(extension)
	switch extension {
	case ".jpg", "jpeg": // check so the "mime" package will not be imported
		return true
	}
	return false
}

# jpeg-junkify

[![Go Report Card](https://goreportcard.com/badge/github.com/jacobpatterson1549/jpeg-junkify)](https://goreportcard.com/report/github.com/jacobpatterson1549/jpeg-junkify)


## an image quality reducer

Jpeg-junkify reduces the quality of JPEG images without decreasing their sizes.  Runs on single files or whole folders.

## Dependencies

[Go 1.17](https://golang.org/dl/) is used to build the application.
[Make](https://www.gnu.org/software/make/) is used to by [Makefile](Makefile) to build and runs the application.

## Build

Run `make` to build and run the application.  It creates the application `jpeg-junkify` in the `build` folder.
To compile for Windows, run `make GO_ARGS="GOOS=windows" OBJ="jpeg-junkify.exe"`.
The `GOARCH` build flag can be added after `GOOS` to specify the CPU architecture: `make "GOOS=linux GOARCH=386"`.  Common values are `amd64`, and `386`.

## Testing

Run `make test` to run the tests for the application.

## Running

The executable application runs on the command line.  Run it with the `-h` parameter for more information: `./build/jpeg-junkify -h`

Examples:
* `./build/jpeg-junkify -b 3M -in-dir ~/Desktop/maps/ -out-dir ~/Desktop/out/`
* `./build/jpeg-junkify -f cows.jpg -b 750KB`
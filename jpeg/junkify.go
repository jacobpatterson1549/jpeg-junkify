// Package jpeg provides utilities on jpeg images.
package jpeg

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
)

// Junkify reduces the quality of the image to be no more than the limit.
// The raw bytes of the reduced-quality image is returned.
// The quality percentage (1-100) is also returned.
// If the image is smaller than the limit, its unaltered bytes are returned.
// If the reader is not a jpeg image, an error is always returned.
func Junkfiy(r io.Reader, bytesLimit int) ([]byte, int, error) {
	if r == nil {
		return nil, 0, fmt.Errorf("missing image")
	}
	var buf1, buf2 bytes.Buffer
	size, err := io.Copy(&buf1, r) // read image
	if err != nil {
		return nil, 0, fmt.Errorf("reading image: %v", err)
	}
	tr := io.TeeReader(&buf1, &buf2)
	img, err := jpeg.Decode(tr) // read image
	if err != nil {
		return nil, 0, fmt.Errorf("decoding image: %v", err)
	}
	if int(size) <= bytesLimit {
		return buf2.Bytes(), 100, nil
	}
	return decreaseQuality(img, bytesLimit)
}

// decreaseQuality shrinks the image to be barely below the limit.
// The quality is decreased using a binary search from quality 1-100.
// This is based on the GuessingGame example at in the Golang source.
// See https://pkg.go.dev/sort#Search
// This will write the image no more than O(8) times to determine the best
// quality search (2^7 = 128, so max 7 searches) + final write = 8
func decreaseQuality(img image.Image, bytesLimit int) (b []byte, quality int, err error) {
	var buf bytes.Buffer
	min, max := 1, 100
	shrinkImage := func(q int) error {
		buf.Reset()
		return jpeg.Encode(&buf, img, &jpeg.Options{Quality: q})
	}
	for max-min > 1 {
		q := min + ((max - min) / 2)
		err := shrinkImage(q)
		switch {
		case err != nil:
			return nil, 0, fmt.Errorf("shrinking image with quality %v%%: %v", q, err)
		case buf.Len() == bytesLimit:
			return buf.Bytes(), q, nil
		case buf.Len() <= bytesLimit:
			min = q
		default:
			max = q
		}
	}
	if err := shrinkImage(min); err != nil {
		return nil, 0, fmt.Errorf("final image shrink with quality %v%%: %v", min, err)
	}
	if buf.Len() > bytesLimit {
		return nil, 0, fmt.Errorf("could not shrink image to have limit of %v bytes with minimal compression", bytesLimit)
	}
	return buf.Bytes(), min, nil
}

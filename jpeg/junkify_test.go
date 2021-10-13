package jpeg

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"testing"
	"testing/iotest"
)

//go:embed test.jpg
var test_jpg []byte // 16536 bytes (~16K)
//go:embed junkify_test.go
var junkify_test_go []byte

func TestJunkify(t *testing.T) {
	tests := []struct {
		img         io.Reader
		bytesLimit  int
		wantOk      bool
		wantSize    int
		wantQuality int
	}{
		{},
		{ // not an image
			img: bytes.NewReader(junkify_test_go),
		},
		{ // read error
			img: iotest.ErrReader(errors.New("read error")),
		},
		{ // smaller than limit
			img:         bytes.NewReader(test_jpg),
			bytesLimit:  1000000,
			wantOk:      true,
			wantSize:    16536,
			wantQuality: 100,
		},
		{ // happy path: quality 73
			img:         bytes.NewReader(test_jpg),
			bytesLimit:  10000,
			wantOk:      true,
			wantSize:    9722,
			wantQuality: 98,
		},
		{ // happy path: quality 50 - should short circut
			img:         bytes.NewReader(test_jpg),
			bytesLimit:  2612,
			wantOk:      true,
			wantSize:    2612,
			wantQuality: 50,
		},
		{ // limit too small
			img:        bytes.NewReader(test_jpg),
			bytesLimit: 1,
		},
	}
	for i, test := range tests {
		gotBytes, gotQuality, err := Junkfiy(test.img, test.bytesLimit)
		switch {
		case !test.wantOk:
			if err == nil {
				t.Errorf("test %v: wanted error", i)
			}
		case err != nil:
			t.Errorf("test %v: unwanted error: %v", i, err)
		case len(gotBytes) != test.wantSize:
			t.Errorf("test %v: sizes not equal:\nwanted: %v\ngot:    %v", i, test.wantSize, len(gotBytes))
		case gotQuality != test.wantQuality:
			t.Errorf("test %v: qualities not equal: wanted: %v got: %v", i, test.wantQuality, gotQuality)
		}
	}
}

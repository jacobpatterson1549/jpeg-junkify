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
	tests := map[string]struct {
		img         io.Reader
		bytesLimit  int
		wantOk      bool
		wantSize    int
		wantQuality int
	}{
		"no image": {},
		"not an image": {
			img: bytes.NewReader(junkify_test_go),
		},
		"image read error": {
			img: iotest.ErrReader(errors.New("read error")),
		},
		"smaller than limit": {
			img:         bytes.NewReader(test_jpg),
			bytesLimit:  1000000,
			wantOk:      true,
			wantSize:    16536,
			wantQuality: 100,
		},
		"smaller size": {
			img:         bytes.NewReader(test_jpg),
			bytesLimit:  10000,
			wantOk:      true,
			wantSize:    9722,
			wantQuality: 98,
		},
		"short circut 50% quality": {
			img:         bytes.NewReader(test_jpg),
			bytesLimit:  2612,
			wantOk:      true,
			wantSize:    2612,
			wantQuality: 50,
		},
		"limit too small": {
			img:        bytes.NewReader(test_jpg),
			bytesLimit: 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotBytes, gotQuality, err := Junkfiy(test.img, test.bytesLimit)
			switch {
			case !test.wantOk:
				if err == nil {
					t.Errorf("wanted error")
				}
			case err != nil:
				t.Errorf("unwanted error: %v", err)
			case len(gotBytes) != test.wantSize:
				t.Errorf("sizes not equal:\nwanted: %v\ngot:    %v", test.wantSize, len(gotBytes))
			case gotQuality != test.wantQuality:
				t.Errorf("qualities not equal: wanted: %v got: %v", test.wantQuality, gotQuality)
			}
		})
	}
}

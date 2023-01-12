package qrcode

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"testing"

	"github.com/chai2010/webp"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	text := "동해물과 백두산이"

	type args struct {
		encode func(w io.Writer, img image.Image) error
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{`png`, args{png.Encode}, false},
		{`jpg`, args{func(w io.Writer, img image.Image) error { return jpeg.Encode(w, img, nil) }}, false},
		{`gif`, args{func(w io.Writer, img image.Image) error { return gif.Encode(w, img, nil) }}, false},
		{`webp`, args{func(w io.Writer, img image.Image) error { return webp.Encode(w, img, nil) }}, false},
		{`avif`, args{func(w io.Writer, img image.Image) error { return webp.Encode(w, img, nil) }}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := Text(text)
			require.NoError(t, err)

			img, err := q.Render(200, 200)
			require.NoError(t, err)

			buf := new(bytes.Buffer)
			err = tt.args.encode(buf, img)
			require.NoError(t, err)

			got, err := Decode(img)
			require.Truef(t, (err != nil) == tt.wantErr, `doSomething() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}

			require.Equal(t, text, got)
		})
	}
}

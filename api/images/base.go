package images

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
)

var typePrefix = []byte("data:image/")
var dataPrefix = []byte("base64,")
var nullBytes = []byte("null")

type Image struct {
	ContentType string
	Data        []byte
}

func New(data []byte) (*Image, error) {
	p := bytes.SplitN(data, []byte{';'}, 2)
	if len(p) < 2 || !bytes.HasPrefix(p[0], typePrefix) || !bytes.HasPrefix(p[1], dataPrefix) {
		return nil, errors.New("invalid image data")
	}

	var rawB64 = p[1][7:]

	var img = Image{
		ContentType: string(p[0][5:]),
		Data:        make([]byte, base64.StdEncoding.DecodedLen(len(rawB64))),
	}
	_, err := base64.StdEncoding.Decode(img.Data, rawB64)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (i Image) Encode() ([]byte, error) {
	if i.ContentType == "" {
		i.ContentType = http.DetectContentType(i.Data)
		if !(i.ContentType == "image/png" || i.ContentType == "image/jpeg" || i.ContentType == "image/gif") {
			return nil, errors.New("invalid content type. accepts png, jpeg and gif")
		}
	}

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(i.Data)))
	base64.StdEncoding.Encode(encoded, i.Data)

	return bytes.Join([][]byte{
		[]byte("data:"),
		[]byte(i.ContentType),
		[]byte(";base64,"),
		encoded,
	}, nil), nil
}

func (i *Image) UnmarshalJSON(data []byte) error {
	data = data[1 : len(data)-1]
	if data[0] == 'n' && data[1] == 'u' && data[2] == 'l' && data[3] == 'l' {
		return nil
	}
	img, err := New(data)
	if err != nil {
		return err
	}
	*i = *img
	return nil
}

func (i *Image) MarshalJSON() ([]byte, error) {
	if len(i.Data) == 0 {
		return nullBytes, nil
	}

	data, err := i.Encode()
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{
		{'"'}, data, {'"'},
	}, nil), nil
}

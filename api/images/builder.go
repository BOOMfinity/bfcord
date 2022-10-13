package images

import (
	"bytes"
	"encoding/base64"
	"github.com/BOOMfinity/go-utils/ubytes"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
)

type MediaBuilder struct {
	err  error
	data []byte
}

func (m *MediaBuilder) Base64(str string) *MediaBuilder {
	data := ubytes.ToBytes(str)
	m.data = data
	return m
}

func (m *MediaBuilder) Reader(reader io.Reader) *MediaBuilder {
	data, err := io.ReadAll(reader)
	if err != nil {
		m.err = err
		return m
	}
	t := http.DetectContentType(data)
	prefix := []byte("data:" + t + ";base64,")
	m.data = make([]byte, base64.StdEncoding.EncodedLen(len(data))+len(prefix))
	copy(m.data, prefix)
	base64.StdEncoding.Encode(m.data[len(prefix):], data)
	return m
}

func (m *MediaBuilder) HTTP(url string) *MediaBuilder {
	resp := fasthttp.AcquireResponse()
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	err := fasthttp.Do(req, resp)
	if err != nil {
		m.err = err
		return m
	}
	return m.Reader(bytes.NewReader(resp.Body()))
}

func (m *MediaBuilder) ToBase64() (string, error) {
	return ubytes.ToString(m.data), m.err
}

func (m *MediaBuilder) Nil() bool {
	return len(m.data) == 0
}

func NewMediaBuilder() *MediaBuilder {
	return &MediaBuilder{}
}

package media

import (
	"encoding/base64"
	"errors"
	"net/http"
	"os"

	"github.com/BOOMfinity/go-utils/ubytes"
	"github.com/valyala/fasthttp"
)

const AttachmentSizeLimit = uint(1024 * 1024 * 25)

type AttachmentResolverFn func() ([]byte, error)

type AttachmentResolver interface {
	Network(url string, limit ...uint) AttachmentResolverFn
	Local(path string) AttachmentResolverFn
	Raw(data []byte) AttachmentResolverFn
	Base64(b64 string) AttachmentResolverFn
}

func ResolveAttachment() AttachmentResolver {
	return attachmentResolver{}
}

type attachmentResolver struct{}

func (a attachmentResolver) Network(url string, limits ...uint) AttachmentResolverFn {
	limit := AttachmentSizeLimit
	if len(limits) > 0 {
		limit = limits[0]
		if limit > AttachmentSizeLimit {
			limit = AttachmentSizeLimit
		}
	}
	return func() ([]byte, error) {
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)
		req.SetRequestURI(url)
		err := fasthttp.Do(req, resp)
		if err != nil {
			return nil, err
		}
		if uint(resp.Header.ContentLength()) > limit {
			return nil, errors.New("maximum body size limit reached")
		}
		return resp.Body(), nil
	}
}

func (a attachmentResolver) Local(path string) AttachmentResolverFn {
	return func() ([]byte, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (a attachmentResolver) Raw(data []byte) AttachmentResolverFn {
	return func() ([]byte, error) {
		return data, nil
	}
}

func (a attachmentResolver) Base64(b64 string) AttachmentResolverFn {
	return func() ([]byte, error) {
		return ubytes.ToBytes(b64), nil
	}
}

func AsImageData(format string, img []byte) string {
	encoded := base64.StdEncoding.EncodeToString(img)
	return "data:" + format + ";base64;" + encoded
}

func ImageDataFromNetwork(url string) (string, error) {
	_, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(body)
	return AsImageData(contentType, body), nil
}

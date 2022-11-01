package cdn

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/errs"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

var Resolver = cdnResolver{}

type cdnResolver struct{}

func (v cdnResolver) EmojiUrl(id snowflake.ID, format ImageFormat) string {
	return fmt.Sprintf("%v/emojis/%v.%v", Url, id, format)
}

func (v cdnResolver) Emoji(id snowflake.ID, format ImageFormat) ([]byte, error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	req.SetRequestURI(v.EmojiUrl(id, format))
	err := fasthttp.Do(req, res)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() == fasthttp.StatusNotFound {
		return nil, errs.HTTPNotFound
	}
	return res.Body(), nil
}

func (v cdnResolver) UserDefaultAvatar(tag uint32) string {
	return fmt.Sprintf("%v/embed/avatars/%v.png", Url, tag%5)
}

func (v cdnResolver) UserAvatar(id snowflake.ID, hash string, size ImageSize, format ImageFormat) string {
	if format == "" {
		format = ImageFormatPNG
	}
	if size == 0 {
		size = ImageSize512
	}
	return fmt.Sprintf("%v/avatars/%v/%v.%v?size=%v", Url, id, hash, format, size)
}

func (v cdnResolver) GuildIconUrl(id snowflake.ID, hash string, format ImageFormat) string {
	return fmt.Sprintf("%v/icons/%v/%v.%v", Url, id, hash, format)
}

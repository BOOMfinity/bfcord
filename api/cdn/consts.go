package cdn

const Url = "https://cdn.discordapp.com"

type ImageSize uint32

const (
	ImageSize16   ImageSize = 8 << 1
	ImageSize32   ImageSize = 8 << 2
	ImageSize64   ImageSize = 8 << 3
	ImageSize128  ImageSize = 8 << 4
	ImageSize256  ImageSize = 8 << 5
	ImageSize512  ImageSize = 8 << 6
	ImageSize1024 ImageSize = 8 << 7
	ImageSize2048 ImageSize = 8 << 8
	ImageSize4096 ImageSize = 8 << 9
)

type ImageFormat string

const (
	ImageFormatPNG  ImageFormat = "png"
	ImageFormatWEBP ImageFormat = "webp"
	ImageFormatJPG  ImageFormat = "jpg"
	ImageFormatJPEG ImageFormat = "jpeg"
	ImageFormatGIF  ImageFormat = "gif"
)

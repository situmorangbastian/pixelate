package pixelate

import "mime/multipart"

type ImageService interface {
	ConvertPngToJpg(file *multipart.FileHeader) (fileName string, err error)
	Resize(file *multipart.FileHeader, scale string) (fileName string, err error)
	Compress(file *multipart.FileHeader) (fileName string, err error)
}

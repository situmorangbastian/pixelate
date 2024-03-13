package pixelate

type ImageService interface {
	ConvertPngToJpg(file string) (fileName string, err error)
	Resize(file string, scale string) (fileName string, err error)
	Compress(file string) (fileName string, err error)
}

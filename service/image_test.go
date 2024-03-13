package service_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/situmorangbastian/pixelate/service"
)

func TestConvertPngToJpg(t *testing.T) {
	tests := []struct {
		testName        string
		invalidFileName string
		expectedResult  string
		expectedError   bool
	}{
		{
			testName:       "success",
			expectedResult: "converted.jpg",
		},
		{
			testName:        "error on open file",
			invalidFileName: "invalid.png",
			expectedError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			pngContent := createPNGFile()
			pngFile, err := os.CreateTemp("", "test-*.png")
			require.NoError(t, err)
			defer os.Remove(pngFile.Name())
			defer pngFile.Close()
			_, err = pngFile.Write(pngContent)
			require.NoError(t, err)

			service := service.NewImageService()

			fileHeader := &multipart.FileHeader{
				Filename: pngFile.Name(),
				Size:     int64(len(pngContent)),
			}

			if test.invalidFileName != "" {
				fileHeader.Filename = test.invalidFileName
			}

			fileName, err := service.ConvertPngToJpg(fileHeader.Filename)
			if test.expectedError {
				require.Error(t, err)
				return
			}
			defer os.Remove(test.expectedResult)
			require.NoError(t, err)
			require.Equal(t, test.expectedResult, fileName)
		})
	}
}

func TestResize(t *testing.T) {
	tests := []struct {
		testName        string
		invalidFileName string
		expectedResult  string
		expectedError   bool
	}{
		{
			testName:       "success",
			expectedResult: "resized.png",
		},
		{
			testName:        "error on open file",
			invalidFileName: "invalid.png",
			expectedError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			pngContent := createPNGFile()
			pngFile, err := os.CreateTemp("", "test-*.png")
			require.NoError(t, err)
			defer os.Remove(pngFile.Name())
			defer pngFile.Close()
			_, err = pngFile.Write(pngContent)
			require.NoError(t, err)

			service := service.NewImageService()

			fileHeader := &multipart.FileHeader{
				Filename: pngFile.Name(),
				Size:     int64(len(pngContent)),
			}

			if test.invalidFileName != "" {
				fileHeader.Filename = test.invalidFileName
			}

			fileName, err := service.Resize(fileHeader.Filename, "10:10")
			if test.expectedError {
				require.Error(t, err)
				return
			}
			defer os.Remove(test.expectedResult)
			require.NoError(t, err)
			require.Equal(t, test.expectedResult, fileName)
		})
	}
}

func TestCompress(t *testing.T) {
	tests := []struct {
		testName        string
		invalidFileName string
		expectedResult  string
		expectedError   bool
	}{
		{
			testName:       "success",
			expectedResult: "compressed.png",
		},
		{
			testName:        "error on open file",
			invalidFileName: "invalid.png",
			expectedError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			pngContent := createPNGFile()
			pngFile, err := os.CreateTemp("", "test-*.png")
			require.NoError(t, err)
			defer os.Remove(pngFile.Name())
			defer pngFile.Close()
			_, err = pngFile.Write(pngContent)
			require.NoError(t, err)

			service := service.NewImageService()

			fileHeader := &multipart.FileHeader{
				Filename: pngFile.Name(),
				Size:     int64(len(pngContent)),
			}

			if test.invalidFileName != "" {
				fileHeader.Filename = test.invalidFileName
			}

			fileName, err := service.Compress(fileHeader.Filename)
			if test.expectedError {
				require.Error(t, err)
				return
			}
			defer os.Remove(test.expectedResult)
			require.NoError(t, err)
			require.Equal(t, test.expectedResult, fileName)
		})
	}
}

func createPNGFile() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	blue := color.RGBA{0, 0, 255, 255}

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, blue)
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

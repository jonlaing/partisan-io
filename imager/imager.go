package imager

import (
	"crypto/rand"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
)

var imgExt = map[string]string{
	"image/jpeg": "jpg",
	"image/jpg":  "jpg",
	"image/png":  "png",
}

// ImageProcessor holds info about image for processing
type ImageProcessor struct {
	File           multipart.File
	origImage      image.Image
	Image          image.Image
	fileType       string
	filenameString string
}

// Resize saves a resized image
func (i *ImageProcessor) Resize(maxWidth uint) (err error) {
	err = i.decode()
	if err != nil {
		return err
	}

	i.Image = resize.Resize(maxWidth, 0, i.origImage, resize.Bicubic)
	if err != nil {
		return err
	}

	return nil
}

// Thumbnail makes a thumbnail of the image
func (i *ImageProcessor) Thumbnail(size int) (err error) {
	err = i.decode()
	if err != nil {
		return err
	}

	err = i.resizeAndCrop(size, size)

	return
}

// Save saves the new image in the specified path
func (i *ImageProcessor) Save(path string) (newPath string, err error) {
	fileName, err := i.FileName()
	if err != nil {
		return
	}

	newPath = path + "/" + fileName

	fullFile, err := os.OpenFile("."+newPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer fullFile.Close()

	switch i.fileType {
	case "image/jpeg", "image/jpg":
		jpeg.Encode(fullFile, i.Image, &jpeg.Options{Quality: 100})
		break
	case "image/png":
		png.Encode(fullFile, i.Image)
		break
	default:
		return "", fmt.Errorf("Unsupported image type: %s", i.fileType)
		break
	}

	return
}

// non-destructive resizing
func (i ImageProcessor) resizeImg(minBound int) (finalImg image.Image) {
	bounds := i.Image.Bounds().Max

	w, h := calcResizeDimensions(bounds, minBound)
	finalImg = resize.Resize(uint(w), uint(h), i.Image, resize.Bicubic)

	return
}

func (i *ImageProcessor) resizeAndCrop(w, h int) (err error) {
	var min int
	if w < h {
		min = w
	} else {
		min = h
	}

	i.Image, err = cutter.Crop(i.resizeImg(min), cutter.Config{
		Width:  w,
		Height: h,
		Mode:   cutter.Centered,
	})

	return
}

func calcResizeDimensions(bounds image.Point, minBound int) (width, height int) {
	if bounds.X > bounds.Y {
		ratio := float64(bounds.X) / float64(bounds.Y)
		width = int(ratio * float64(minBound))
		height = minBound
	} else if bounds.Y > bounds.X {
		ratio := float64(bounds.Y) / float64(bounds.X)
		width = minBound
		height = int(ratio * float64(minBound))
	}

	return
}

// FileName generates a random filename plus extension
func (i *ImageProcessor) FileName() (string, error) {
	if i.filenameString != "" {
		return i.filenameString, nil
	}

	if err := i.detectFileType(); err != nil {
		return "", err
	}

	extension, ok := imgExt[i.fileType]
	if !ok {
		return "", fmt.Errorf("Filetype not supported: %s", i.fileType)
	}

	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 16)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	i.filenameString = string(bytes) + "." + extension
	return i.filenameString, nil
}

func (i *ImageProcessor) detectFileType() error {
	i.File.Seek(0, 0)

	if i.fileType != "" {
		return nil
	}

	buf := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err := i.File.Read(buf)
	if err != nil {
		return err
	}

	i.fileType = http.DetectContentType(buf)

	// Resetting read of tmpFile (otherwise we'd copy an incomplete file)
	_, err = i.File.Seek(0, 0)
	return err
}

func (i *ImageProcessor) decode() (err error) {
	i.origImage, _, err = image.Decode(i.File)
	return
}

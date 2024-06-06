package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	// "io"

	"github.com/disintegration/imaging"
)

type ImageUser struct {
	ImageKTP []byte `json:"image_ktp" `
}

func main() {
	src, err := imaging.Open("original/6.jpg", imaging.AutoOrientation(true))

	if err != nil {
		// Handle errors while opening the image
		fmt.Printf("failed to open image: %v\n", err)
		return
	}

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	src = imaging.Resize(src, 800, 0, imaging.Lanczos)

	// defer src.Close()
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, src, nil)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	files := buf.Bytes()

	catchmodels := ImageUser{
		ImageKTP: files,
	}

	cam1 := base64.StdEncoding.EncodeToString(catchmodels.ImageKTP)

	fmt.Println(cam1)

}

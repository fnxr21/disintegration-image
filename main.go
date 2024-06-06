package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	// "os"
	"time"

	"github.com/disintegration/imaging"
)

type ImageUser struct {
	ImageKTP []byte `json:"image_ktp" `
}

func main() {
	timestamp := time.Now().Add(time.Second).Unix()
	src, err := imaging.Open("original/6.jpg", imaging.AutoOrientation(true))

	if err != nil {
		// Handle errors while opening the image
		fmt.Printf("failed to open image: %v\n", err)
		return
	}

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	src = imaging.Resize(src, 800, 0, imaging.Lanczos)

	// Create a filename with timestamp and .jpg extension
	filename := fmt.Sprintf("3L/rnd-"+"%d.jpg", timestamp)
	err = imaging.Save(src, filename)

	if err != nil {
		// Handle errors while saving the image
		fmt.Printf("failed to save image: %v\n", err)
		return
	}

	files, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	
	defer files.Close()
	
	data, err := io.ReadAll(files)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	catchmodels := ImageUser{
		ImageKTP: data,
		}

	cam1 := base64.StdEncoding.EncodeToString(catchmodels.ImageKTP)

	fmt.Println(cam1)
	defer os.Remove(filename)

	//save to database
	// Create and save image

}

package main

import (
	"log"

	"github.com/disintegration/imaging"
)

func main() {
	src, err := imaging.Open("original/1.jpg")

	if err != nil {
		log.Printf("failed to open :%v", err)
	}

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	src = imaging.Resize(src, 800, 0, imaging.Lanczos)

	err = imaging.Save(src, "3L/1.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

}

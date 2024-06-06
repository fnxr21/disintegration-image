package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"

	// "io"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
)

// type ImageUser struct {
// 	ImageKTP []byte `json:"image_ktp" `
// }

func main() {

	e := echo.New()
	e.GET("/", (
		func(c echo.Context) error {

			imageUser := c.Get("dataFile")
			Image := imageUser.([]byte)

			Imageuser := base64.StdEncoding.EncodeToString(Image)
			return c.JSON(http.StatusOK, Imageuser)
		}))
	e.Logger.Fatal(e.Start(":1323"))
}

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		files, err := c.FormFile("image")

		f, _ := files.Open()

		defer f.Close()
		imgpng, _ := png.Decode(f)
		fmt.Println(imgpng)
		// defer imgpng.Close()
		img, _ := jpeg.Decode(f)

		srcs := imaging.Resize(img, 800, 0, imaging.Lanczos)

		// defer src.Close()
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, srcs, nil)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return err
		}

		filesa := buf.Bytes()

		c.Set("dataFile", filesa)
		return next(c)
	}
}

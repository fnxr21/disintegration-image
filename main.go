package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	e.POST("/", UploadFile(func(c echo.Context) error {

		imageUser := c.Get("dataFile")
		Image := imageUser.([]byte)

		Imageuser := base64.StdEncoding.EncodeToString(Image)
		return c.JSON(http.StatusOK, Imageuser)
	}))
	e.Logger.Fatal(e.Start(":1323"))
}

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("image")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing image file")
		}

		if file != nil {

			imageBytes, err := file.Open()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error opening uploaded file: %v", err))
			}

			defer imageBytes.Close()

			buffer := make([]byte, file.Size)
			_, err = imageBytes.Read(buffer)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error reading uploaded file data: %v", err))
			}

			img, err := imaging.Decode(bytes.NewReader(buffer), imaging.AutoOrientation(true))

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid image format: %v", err))
			}

			srcs := imaging.Resize(img, 500, 0, imaging.Lanczos)

			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, srcs, nil)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error encoding resized image: %v", err))
			}

			filesa := buf.Bytes()
			c.Set("dataFile", filesa)
			return next(c)
		}
		c.Set("dataFile", "")
		return next(c)
	}
}

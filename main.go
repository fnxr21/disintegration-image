package main

import (
	"bytes"
	"encoding/base64"
	"strconv"

	// "reflect"
	// "strconv"

	"fmt"
	"image/jpeg"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	e.POST("/", UploadFiles(func(c echo.Context) error {

		imageUser := c.Get("dataFile").(map[string][]byte)

		// var imageMap map[string]string
		imageMap := make(map[string]string)

		TotalImage := len(imageUser)
		Counter := 0

		for filename, imageData := range imageUser {
			// Separate logic here:
			fmt.Println("Processing image:", filename)
			Imageuser := base64.StdEncoding.EncodeToString(imageData)

			Counter = Counter + 1
			TotalName := strconv.Itoa(Counter)

			imageMap["image"+TotalName] = Imageuser

			if TotalImage == Counter {
				break
			}

		}
		fmt.Println(imageMap)
		// fmt.Println(reflect.TypeOf(imageUser))

		// Image := imageUser.([]byte)

		// Imageuser := base64.StdEncoding.EncodeToString(Image)
		return c.JSON(http.StatusOK, imageMap)
	}))
	e.Logger.Fatal(e.Start(":1323"))
}
func UploadFiles(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		form, err := c.MultipartForm()
		if err != nil {
			return echo.ErrBadRequest
		}

		files, ok := form.File["images"]
		if !ok {
			return echo.ErrNotFound
		}
		// var myMap map[string][]byte
		myMap := make(map[string][]byte)
		for _, file := range files {
			// Extract file extension
			fileExt := filepath.Ext(file.Filename)

			// Generate unique filename
			originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
			now := time.Now()
			filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt

			imageSize := resizeImage(file.Size)

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid image format: %v", err))
			}

			imageBytes, _ := file.Open()

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

			srcs := imaging.Resize(img, imageSize, 0, imaging.Lanczos)
			// ok := strconv.Itoa(int(time.Now().Add(time.Second).Unix()))

			// imaging.Save(srcs, "3L/"+ok+".jpg")
			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, srcs, nil)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error encoding resized image: %v", err))
			}

			datafile := buf.Bytes()

			myMap[filename] = datafile

		}

		c.Set("dataFile", myMap)
		return next(c)
	}
}

func resizeImage(fileSize int64) int {

	if fileSize == 0 {
		return 0
	}

	if fileSize > 10*1000000 {
		return 500
	}
	return 800
}

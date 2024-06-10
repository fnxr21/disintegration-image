package main

import (
	"bytes"
	// "encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "strconv"
	// "time"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	e.POST("/", UploadFiles(func(c echo.Context) error {

		// imageUser := c.Get("dataFile")
		// Image := imageUser.([]byte)

		// Imageuser := base64.StdEncoding.EncodeToString(Image)
		return c.JSON(http.StatusOK, "Imageuser")
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
		// filePaths := []string{}
		for _, file := range files {
			// Extract file extension
			fileExt := filepath.Ext(file.Filename)

			// Generate unique filename
			originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
			now := time.Now()
			filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt

			// Create destination file path
			filePath := "./3L/" + filename

			// filePaths = append(filePaths, filePath)

			// Open destination file
			out, err := os.Create(filePath)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			defer out.Close()

			// Open uploaded file
			readerFile, err := file.Open()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			defer readerFile.Close()

			// Copy uploaded file content
			_, err = io.Copy(out, readerFile)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}
		c.Set("dataFile", "")
		return next(c)
	}
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
			// Check for valid file type (JPG or PNG)
			if contentType := file.Header.Get("Content-Type"); !(contentType == "image/jpeg" || contentType == "image/png") {
				return c.JSON(http.StatusBadRequest, "Invalid file type. Only JPG and PNG are allowed")
			}

			imageSize := resizeImage(file.Size)

			fmt.Println(imageSize, "ukuran")
			fmt.Println(file.Size, "ukuran size")

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

			filesa := buf.Bytes()

			fmt.Println(len(filesa), "size terakhir")
			c.Set("dataFile", filesa)
			return next(c)
		}
		c.Set("dataFile", "")
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

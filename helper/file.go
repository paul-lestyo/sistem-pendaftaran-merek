package helper

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
)

func CheckInputFile(c *fiber.Ctx, inputFileName string) (fileInput FileInput, ok bool) {
	file, err := c.FormFile(inputFileName)
	fileInput = FileInput{}
	ok = false

	if err == nil && file.Size > 0 {
		fileInput = FileInput{
			Path:        filepath.Dir(file.Filename),
			Filename:    filepath.Base(file.Filename),
			Ext:         filepath.Ext(file.Filename),
			ContentType: file.Header.Get("Content-Type"),
			Size:        file.Size,
		}
		ok = true
	}

	return fileInput, ok
}

func UploadFile(c *fiber.Ctx, inputFileName string, filepath string) (string, bool) {
	filename := ""
	file, err := c.FormFile(inputFileName)
	if err != nil {
		return filename, false
	}

	if file.Size != 0 {
		filename = "/uploads/" + filepath + "/" + file.Filename
		fmt.Println("hoho: ", filename)
		err := c.SaveFile(file, "assets"+filename)
		PanicIfError(err)
	}
	return filename, true
}

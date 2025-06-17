package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UplaodProductImageHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "image data not found",
		})
	}

	//validasi gambar

	//validasi extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowExts[ext] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "image extension is not allowed (jpg, jpeg, png, webp)",
		})
	}

	//validasi content type
	contentType := file.Header.Get("Content-Type")
	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedContentTypes[contentType] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "content type is not allowed (image/jpeg, image/png, image/webp)",
		})
	}

	//product_16121.png
	timestamp := time.Now().UnixNano()
	fileName := fmt.Sprintf("product_%d%s", timestamp, filepath.Ext(file.Filename))
	uploadPath := "./storage/product/" + fileName
	err = c.SaveFile(file, uploadPath)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error saving image",
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "image uploaded successfully",
		"file_name": fileName,
	})
}

package main

import (
	"context"
	"log"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/handler"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/service"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func handleGetFileName(c *fiber.Ctx) error {
	fileNameParam := c.Params("filename")
	filePath := path.Join("storage", "product", fileNameParam)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).SendString("File not found")
		}

		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Internal server error")
	}

	// buka file
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Internal server error")
	}

	// kirimkan file sbg response
	ext := path.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	c.Set("Content-Type", mimeType)
	return c.SendStream(file)
}

func main() {
	godotenv.Load()
	ctx := context.Background()
	app := fiber.New()

	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))
	orderRepository := repository.NewOrderRepository(db)
	webhookSerivce := service.NewWebhookService(orderRepository)
	webhookHandler := handler.NewWebhookHandler(webhookSerivce)

	app.Use(cors.New())

	// app.Static("/storage/product", "./storage/product")
	// localhost:3000/storage/product/name_file.png
	app.Get("/storage/product/:filename", handleGetFileName)
	app.Post("/product/upload", handler.UplaodProductImageHandler)
	app.Post("/webhook/xendit/invoice", webhookHandler.ReceiveInvoice)

	app.Listen(":3000")
}

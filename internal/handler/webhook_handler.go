package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type webhookHandler struct {
}

func (w *webhookHandler) ReceiveInvoice(c *fiber.Ctx) error {
	fmt.Println(string(c.Body()))
	return c.SendStatus(http.StatusOK)
}

func NewWebhookHandler() *webhookHandler {
	return &webhookHandler{}
}

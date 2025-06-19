package handler

import (
	"log"
	"net/http"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/dto"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/service"
	"github.com/gofiber/fiber/v2"
)

type webhookHandler struct {
	webhookService service.IWebhookService
}

func (w *webhookHandler) ReceiveInvoice(c *fiber.Ctx) error {
	var request dto.XenditInvoiceRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	w.webhookService.ReceiveInvoice(c.UserContext(), &request)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func NewWebhookHandler(webhootService service.IWebhookService) *webhookHandler {
	return &webhookHandler{
		webhookService: webhootService,
	}
}

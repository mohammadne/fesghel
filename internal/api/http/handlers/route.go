package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/mohammadne/fesghel/internal/urls"
)

func NewRoute(r fiber.Router, logger *zap.Logger, urls urls.Service) {
	handler := &route{
		logger: logger,
		urls:   urls,
	}

	r.Get("/:id", handler.moveURL)
}

type route struct {
	logger *zap.Logger
	urls   urls.Service
}

func (r *route) moveURL(c fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	url, err := r.urls.Retrieve(c.Context(), id)
	if err != nil {
		if errors.Is(err, urls.ErrIDNotExists) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusMovedPermanently).JSON(map[string]string{
		"url": string(url),
	})
}

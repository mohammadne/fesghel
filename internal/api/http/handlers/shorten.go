package handlers

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/mohammadne/fesghel/internal/api/http/i18n"
	"github.com/mohammadne/fesghel/internal/api/http/models"
	"github.com/mohammadne/fesghel/internal/entities"
	"github.com/mohammadne/fesghel/internal/urls"
)

func NewShorten(r fiber.Router, logger *zap.Logger, i18n i18n.I18N, urls urls.Service) {
	handler := &shorten{
		logger: logger,
		i18n:   i18n,
		urls:   urls,
	}

	g := r.Group("shorten")
	g.Post("/", handler.listUsers)
	g.Get("/:id", handler.listUsers)
}

type shorten struct {
	logger *zap.Logger
	i18n   i18n.I18N
	urls   urls.Service
}

func (s *shorten) listUsers(c fiber.Ctx) error {
	response := &models.Response{}
	language, _ := c.Locals("language").(entities.Language)

	// response.Request = s.usecase.ListUsers(c.Context())
	response.Message = s.i18n.Translate("users.list_users.success", language)
	return response.Write(c, fiber.StatusOK)
}

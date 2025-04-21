package handlers

import (
	"errors"

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
	g.Post("/", handler.shortenURL)
	g.Get("/:id", handler.retrieveURL)
}

type shorten struct {
	logger *zap.Logger
	i18n   i18n.I18N
	urls   urls.Service
}

func (s *shorten) shortenURL(c fiber.Ctx) error {
	response := &models.Response{}
	language, _ := c.Locals("language").(entities.Language)

	request := models.ShortenRequest{}
	if err := c.Bind().Body(&request); err != nil {
		response.Message = s.i18n.Translate("shorten.shorten_url.error_request", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	id, err := s.urls.Shorten(c.Context(), entities.URL(request.URL))
	if err != nil {
		s.logger.Error("error retreiving the data", zap.Error(err))
		response.Message = s.i18n.Translate("shorten.shorten_url.error_shorten", language)
		return response.Write(c, fiber.StatusInternalServerError)
	}

	response.Request = models.ShortenURLResponse{ID: id}
	response.Message = s.i18n.Translate("shorten.shorten_url.success", language)
	return response.Write(c, fiber.StatusCreated)
}

func (s *shorten) retrieveURL(c fiber.Ctx) error {
	response := &models.Response{}
	language, _ := c.Locals("language").(entities.Language)

	id := c.Params("id")
	if len(id) == 0 {
		response.Message = s.i18n.Translate("shorten.retrieve_url.id_not_given", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	url, err := s.urls.Retrieve(c.Context(), id)
	if err != nil {
		if errors.Is(err, urls.ErrIDNotExists) {
			s.logger.Error("error id not exists", zap.String("id", id))
			response.Message = s.i18n.Translate("shorten.retrieve_url.not_exists", language)
			return response.Write(c, fiber.StatusNotFound)
		}

		s.logger.Error("error retreiving the url", zap.Error(err))
		response.Message = s.i18n.Translate("shorten.retrieve_url.error", language)
		return response.Write(c, fiber.StatusInternalServerError)
	}

	response.Request = models.RetrieveURLResponse{URL: url}
	response.Message = s.i18n.Translate("shorten.retrieve_url.success", language)
	return response.Write(c, fiber.StatusOK)
}

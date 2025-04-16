package models

import (
	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Message string `json:"message"`
	Code    string `json:"data,omitempty"`
	Request any    `json:"request,omitempty"`
}

func (response *Response) Write(ctx fiber.Ctx, statusCode int) error {
	ctx.Set("Content-Type", "application/json")
	return ctx.Status(statusCode).JSON(&response)
}
